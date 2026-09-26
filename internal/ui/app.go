package ui

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/contexts"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const refreshInterval = 30 * time.Second

type Model struct {
	version string

	registry     *resource.Registry
	resource     resource.Resource
	resourceRows []resource.Row

	table        table.Model
	tableXOffset int

	err error

	itemCount    int
	status       string
	loading      bool
	loadingLabel string
	showLoading  bool
	loaded       bool
	loadID       uint64

	autoRefreshTimestamp *time.Time
	autoRefreshPaused    bool

	width  int
	height int

	activatingContext bool
	context           *openstack.Context

	navigation []navigationEntry
	navigateID string
	filter     *resourceFilter
	scope      resourceScope

	detailMode    bool
	detail        viewport.Model
	detailID      string
	detailContent string

	commandMode bool
	command     textinput.Model
}

func New(registry *resource.Registry, openstackContext *openstack.Context, version string) Model {
	r, ok := registry.Get("contexts")
	if !ok {
		return Model{
			version:  version,
			registry: registry,
			context:  openstackContext,
			err:      fmt.Errorf("contexts resource not registered"),
		}
	}

	columns := make([]table.Column, 0, len(r.Columns()))

	for _, column := range r.Columns() {
		columns = append(columns, table.Column{
			Title: column.Title,
			Width: column.MinWidth,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	t.SetStyles(tableStyles())

	command := textinput.New()
	command.Prompt = ":"
	command.CharLimit = 64

	detail := viewport.New(1, 1)
	detail.SetHorizontalStep(4)

	return Model{
		version: version,

		registry: registry,
		resource: r,
		table:    t,
		command:  command,
		context:  openstackContext,
		detail:   detail,

		loading:      true,
		showLoading:  true,
		loadingLabel: "Loading...",
		loaded:       false,
		loadID:       1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadResource(),
		autoRefresh(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "esc", "enter":
				m.err = nil

				if len(m.navigation) > 0 {
					return m, m.navigateBack()
				}

				return m, nil
			case "ctrl+x":
				return m, tea.Quit
			}
		}

		return m, nil
	}

	switch msg := msg.(type) {
	case errMsg:
		if msg.loadID != 0 && msg.loadID != m.loadID {
			return m, nil
		}

		m.autoRefreshPaused = true
		m.loading = false
		m.showLoading = false
		m.loadingLabel = ""

		m.err = msg.err

		slog.Error(msg.err.Error(), "cloud", m.context.Cloud)
		return m, nil

	case clipboardResultMsg:
		m.status = "copied to clipboard"
		return m, clearStatus(m.status)

	case clearStatusMsg:
		if m.status == msg.status {
			m.status = ""
		}
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		if m.detailMode {
			m.detailMode = false
			m.detailContent = ""
			m.detailID = ""
			return m, nil
		}

		m.status = ""
		m.commandMode = false
		m.command.Blur()
		m.command.SetValue("")

		if len(m.navigation) > 0 {
			return m, m.navigateBack()
		}

		return m, nil
	}

	if m.detailMode {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "ctrl+x":
				return m, tea.Quit
			case "c":
				return m, m.copyDetailsViewport()
			}
		}

		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.Resize()
			return m, nil

		case autoRefreshMsg:
			return m, autoRefresh()
		}

		var cmd tea.Cmd
		m.detail, cmd = m.detail.Update(msg)

		return m, cmd
	}

	if m.commandMode {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				command := m.command.Value()

				m.commandMode = false
				m.command.Blur()
				m.command.SetValue("")

				if m.registry.IsNavigationOnly(command) {
					res, _ := m.registry.Get(command)
					m.status = fmt.Sprintf("resource '%s' is accessible only through its parents", res.Kind())
					return m, nil
				}

				return m, m.switchResource(command, nil)
			}
		}

		var cmd tea.Cmd
		m.command, cmd = m.command.Update(msg)

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.activatingContext {
			return m, nil
		}

		switch msg.String() {
		case "ctrl+x":
			return m, tea.Quit

		case ":":
			m.status = ""
			m.commandMode = true
			m.command.SetValue("")
			m.command.Focus()

			return m, textinput.Blink

		case "r":
			return m, m.refreshResource()

		case "c":
			if m.detailMode {

				return m, m.copyDetailsViewport()
			}

			return m, m.copyTableViewport()

		case "enter":
			return m, m.executeDefaultCommand()

		case "left":
			m.tableXOffset = max(m.tableXOffset-4, 0)
			return m, nil

		case "right":
			m.tableXOffset = min(
				m.tableXOffset+4,
				m.maxTableXOffset(),
			)
			return m, nil
		}

		if cmd := m.executeResourceCommand(commandKey(msg)); cmd != nil {
			return m, cmd
		}

	case resourcesLoadedMsg:
		if msg.loadID != m.loadID {
			return m, nil
		}

		m.loading = false
		m.showLoading = false

		m.loaded = true
		m.autoRefreshPaused = false

		if m.resource == nil {
			m.status = "no active resource"
			m.autoRefreshPaused = true
			return m, nil
		}

		msg.rows = filterResourceRows(msg.rows, m.filter)
		columns := m.resource.Columns()
		rows := make([]table.Row, 0, len(msg.rows))

		for _, row := range msg.rows {
			values := make(table.Row, 0, len(columns))

			for _, column := range columns {
				value := row.Fields[column.Key]

				if m.resource.Kind() == "contexts" &&
					column.Key == "active" &&
					row.ID == m.context.Cloud {
					value = "true"
				}

				values = append(values, value)
			}

			rows = append(rows, values)
		}

		m.resourceRows = msg.rows
		m.table.SetRows(rows)
		m.itemCount = len(rows)

		if m.navigateID != "" {
			msg.selectedID = m.navigateID
			m.navigateID = ""
		}

		cursor := 0

		if len(msg.rows) > 0 {
			found := false

			if msg.selectedID != "" {
				for i, row := range msg.rows {
					if row.ID == msg.selectedID {
						cursor = i
						found = true
						break
					}
				}
			}

			if !found {
				cursor = min(msg.cursor, len(msg.rows)-1)
			}

			m.table.SetCursor(0)
			m.table.MoveDown(cursor)
		}

		t := time.Now()
		m.autoRefreshTimestamp = &t

		return m, nil

	case autoRefreshMsg:
		cmd := m.autoRefreshResource()

		t := time.Now()
		m.autoRefreshTimestamp = &t

		return m, tea.Batch(
			cmd,
			autoRefresh(),
		)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.Resize()

		return m, nil

	case contexts.ActivatedMsg:
		m.activatingContext = false
		m.showLoading = false
		m.loadingLabel = ""

		if msg.Err != nil {
			m.err = fmt.Errorf("context activation failed: %v", msg.Err)
			slog.Error(m.err.Error(), "cloud", msg.Context.Cloud)
			return m, nil
		}

		*m.context = *msg.Context
		m.status = fmt.Sprintf("connected to %s", m.context.Cloud)

		m.loadingLabel = "Loading..."
		cmd := m.refreshResource()

		return m, tea.Batch(
			clearStatus(m.status),
			cmd,
		)

	case resource.NavigateMsg:
		cursor := m.table.Cursor()

		if cursor >= 0 && cursor < len(m.resourceRows) {
			m.navigation = append(m.navigation, navigationEntry{
				resource: m.resource.Kind(),
				id:       m.resourceRows[cursor].ID,
			})
		}

		m.navigateID = msg.ID

		return m, m.navigateResource(msg.Resource)

	case resource.NavigateFilteredMsg:
		cursor := m.table.Cursor()

		if cursor >= 0 && cursor < len(m.resourceRows) {
			m.navigation = append(m.navigation, navigationEntry{
				resource: m.resource.Kind(),
				id:       m.resourceRows[cursor].ID,
				filter:   m.filter,
			})
		}

		m.navigateID = ""
		m.filter = &resourceFilter{
			field: msg.Field,
			value: msg.Value,
		}

		return m, m.navigateFilteredResource(msg.Resource)

	case resource.NavigateScopedMsg:
		cursor := m.table.Cursor()

		if cursor >= 0 && cursor < len(m.resourceRows) {
			m.navigation = append(m.navigation, navigationEntry{
				resource: m.resource.Kind(),
				id:       m.resourceRows[cursor].ID,
				filter:   m.filter,
				scope:    m.scope,
			})
		}

		m.navigateID = ""
		m.filter = nil
		m.scope = msg.Scope

		return m, m.navigateScopedResource(msg.Resource)

	case resource.DetailsMsg:
		m.showLoading = false
		m.loadingLabel = ""

		if msg.Err != nil {
			m.err = fmt.Errorf("loading resource details failed: %v", msg.Err)
			slog.Error(m.err.Error(), "cloud", m.context.Cloud)
			return m, nil
		}

		data, err := json.MarshalIndent(msg.Content, "", "  ")
		if err != nil {
			m.err = fmt.Errorf("encoding details: %v", err)
			return m, nil
		}

		m.detailMode = true
		m.detailID = msg.ID
		m.detailContent = string(data)
		m.detail.SetContent(colorizeJSON(data))
		m.detail.GotoTop()
		m.Resize()

		return m, nil

	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	header := m.renderHeader()

	contentView := m.renderTable()

	if m.detailMode {
		contentView = m.renderDetails()
	}

	resourceLine := ""

	if m.resource != nil {
		var resourceTags strings.Builder

		for _, entry := range m.navigation {
			resourceTags.WriteString(
				navigationTagStyle.Render(" <"+entry.resource+"> ") + " ",
			)
		}

		if m.detailMode {
			resourceTags.WriteString(
				navigationTagStyle.Render(" <"+m.resource.Kind()+"> ") + " ",
			)
			resourceTags.WriteString(
				resourceTagStyle.Render(" <show> "),
			)
		} else {
			resourceTags.WriteString(
				resourceTagStyle.Render(" <" + m.resource.Kind() + "> "),
			)
		}

		rightTag := loadingStyle.Render(" AutoRefresh ")
		autoRefreshMsgState := autoRefreshStateOnStyle.Render(" ON ")
		if m.autoRefreshPaused {
			autoRefreshMsgState = autoRefreshStateOffStyle.Render(" OFF ")
		}
		timeTag := ""
		if m.autoRefreshTimestamp != nil {
			timeTag = loadingStyle.Render(fmt.Sprintf(" %s ", m.autoRefreshTimestamp.Format(time.RFC3339)))
		}
		rightTag = lipgloss.JoinHorizontal(lipgloss.Left, rightTag, autoRefreshMsgState, timeTag)

		if m.showLoading {
			loadingTag := loadingStyle.Render(fmt.Sprintf(" %s ", m.loadingLabel)) + " "
			rightTag = lipgloss.JoinHorizontal(lipgloss.Left, loadingTag, rightTag)
		}

		resourceLine = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().
				Width(max(m.width-lipgloss.Width(rightTag), 1)).
				Render(resourceTags.String()),
			rightTag,
		)
	}

	commandLine := ""

	if m.status != "" {
		commandLine = statusStyle.
			Width(m.width).
			Render(m.status)
	}

	if m.commandMode {
		commandLine = commandStyle.Render(m.command.View())
	}

	view := header +
		"\n\n" +
		contentView +
		"\n" +
		resourceLine +
		"\n" +
		commandLine

	if m.err != nil {
		if m.err != nil {
			return overlayCenter(view, m.renderErrorModal())
		}
	}

	return view
}

func (m *Model) Resize() {
	const (
		headerHeight         = 7
		footerHeight         = 2
		tableContainerBorder = 2
	)

	contentHeight := max(
		m.height-
			headerHeight-
			footerHeight-
			tableContainerBorder,
		1,
	)

	m.table.SetHeight(contentHeight)

	m.detail.Width = max(m.width-2, 1)
	m.detail.Height = contentHeight

	if m.resource == nil {
		return
	}

	resourceColumns := m.resource.Columns()

	minWidth := 0
	totalFlex := 0

	for _, column := range resourceColumns {
		minWidth += column.MinWidth
		totalFlex += column.Flex
	}

	tableWidth := max(m.width-2, 1)

	const tableHorizontalPadding = 2

	contentWidth := max(
		tableWidth-(len(resourceColumns)*tableHorizontalPadding),
		1,
	)

	extra := max(contentWidth-minWidth, 0)

	columns := make([]table.Column, 0, len(resourceColumns))

	for _, column := range resourceColumns {
		width := column.MinWidth

		if totalFlex > 0 {
			width += extra * column.Flex / totalFlex
		}

		columns = append(columns, table.Column{
			Title: column.Title,
			Width: width,
		})
	}

	m.table.SetColumns(columns)

	m.tableXOffset = min(
		m.tableXOffset,
		m.maxTableXOffset(),
	)
}

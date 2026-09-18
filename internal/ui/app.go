package ui

import (
	"encoding/json"
	"fmt"
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
	registry     *resource.Registry
	resource     resource.Resource
	resourceRows []resource.Row
	table        table.Model
	err          error

	itemCount         int
	status            string
	loading           bool
	showLoading       bool
	loaded            bool
	loadID            uint64
	autoRefreshPaused bool

	width  int
	height int

	context *openstack.Context

	navigation []navigationEntry
	navigateID string
	filter     *resourceFilter

	detailMode    bool
	detail        viewport.Model
	detailID      string
	detailContent string

	commandMode bool
	command     textinput.Model
}

func New(registry *resource.Registry, openstackContext *openstack.Context) Model {
	r, ok := registry.Get("contexts")
	if !ok {
		return Model{
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
		registry: registry,
		resource: r,
		table:    t,
		command:  command,
		context:  openstackContext,
		detail:   detail,

		loading:     true,
		showLoading: true,
		loaded:      false,
		loadID:      1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadResource(),
		autoRefresh(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case clipboardResultMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("clipboard failed: %v", msg.err)
			return m, nil
		}

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
				return m, m.copyViewport()
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

				return m, m.switchResource(command)
			}
		}

		var cmd tea.Cmd
		m.command, cmd = m.command.Update(msg)

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			return m, m.copyViewport()

		case "enter":
			return m, m.executeDefaultCommand()
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

		if msg.err != nil {
			m.status = msg.err.Error()
			m.autoRefreshPaused = true
			return m, nil
		}

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

		return m, nil

	case autoRefreshMsg:
		cmd := m.autoRefreshResource()

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
		if msg.Err != nil {
			m.status = fmt.Sprintf("context activation failed: %v", msg.Err)
			return m, nil
		}

		*m.context = *msg.Context
		m.status = fmt.Sprintf("connected to %s", m.context.Cloud)

		cmd := m.autoRefreshResource()

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

	//case servers.ShowMsg:
	//	if msg.Err != nil {
	//		m.status = msg.Err.Error()
	//		return m, nil
	//	}
	//
	//	data, err := json.MarshalIndent(msg.Server, "", "  ")
	//	if err != nil {
	//		m.status = fmt.Sprintf("encoding server details: %v", err)
	//		return m, nil
	//	}
	//
	//	m.detailMode = true
	//	m.detailID = msg.Server.ID
	//	m.detailContent = string(data)
	//	m.detail.SetContent(colorizeJSON(data))
	//	m.detail.GotoTop()
	//	m.Resize()
	//
	//	return m, nil

	case resource.DetailsMsg:
		if msg.Err != nil {
			m.status = msg.Err.Error()
			return m, nil
		}

		data, err := json.MarshalIndent(msg.Content, "", "  ")
		if err != nil {
			m.status = fmt.Sprintf("encoding details: %v", err)
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
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

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

		loadingTag := ""

		if m.showLoading {
			loadingTag = loadingStyle.Render(" Loading... ")
		}

		resourceLine = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().
				Width(max(m.width-lipgloss.Width(loadingTag), 1)).
				Render(resourceTags.String()),
			loadingTag,
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

	return header +
		"\n\n" +
		contentView +
		"\n" +
		resourceLine +
		"\n" +
		commandLine
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
}

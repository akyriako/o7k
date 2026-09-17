package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/contexts"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const refreshInterval = 10 * time.Second

type resourcesLoadedMsg struct {
	loadID     uint64
	rows       []resource.Row
	selectedID string
	cursor     int
	err        error
}

type clearStatusMsg struct {
	status string
}

type navigateMsg struct {
	resource string
	id       string
}

type navigationEntry struct {
	resource string
	id       string
}

type autoRefreshMsg struct{}

type Model struct {
	registry     *resource.Registry
	resource     resource.Resource
	resourceRows []resource.Row
	table        table.Model
	err          error

	itemCount   int
	status      string
	loading     bool
	showLoading bool
	loaded      bool
	loadID      uint64

	width  int
	height int

	context *openstack.Context

	navigation []navigationEntry
	navigateID string

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

	return Model{
		registry: registry,
		resource: r,
		table:    t,
		command:  command,
		context:  openstackContext,

		loading:     true,
		showLoading: true,
		loaded:      false,
		loadID:      1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadResource(),
		autoRefreshCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Esc always means "cancel/dismiss and return to normal mode".
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		m.status = ""
		m.commandMode = false
		m.command.Blur()
		m.command.SetValue("")

		if len(m.navigation) > 0 {
			return m, m.navigateBack()
		}

		return m, nil
	}

	// While command mode is active, keyboard input belongs to the
	// command text input instead of the resource table.
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
		case "q", "ctrl+c":
			return m, tea.Quit

		case ":":
			m.status = ""
			m.commandMode = true
			m.command.SetValue("")
			m.command.Focus()

			return m, textinput.Blink

		case "r":
			return m, m.refreshResource()

		case "enter":
			return m, m.executeDefaultCommand()
		}

		if cmd := m.executeResourceCommand(commandKey(msg)); cmd != nil {
			return m, cmd
		}

	case resourcesLoadedMsg:
		// Ignore responses belonging to an older load.
		if msg.loadID != m.loadID {
			return m, nil
		}

		m.loading = false
		m.showLoading = false

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		m.loaded = true

		if m.resource == nil {
			m.err = fmt.Errorf("no active resource")
			return m, nil
		}

		columns := m.resource.Columns()
		rows := make([]table.Row, 0, len(msg.rows))

		for _, row := range msg.rows {
			values := make(table.Row, 0, len(columns))

			for _, column := range columns {
				value := row.Fields[column.Key]

				if m.resource.Kind() == "contexts" && column.Key == "active" && row.ID == m.context.Cloud {
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

		// tea.Tick fires once, so every tick must schedule the next one.
		return m, tea.Batch(
			cmd,
			autoRefreshCmd(),
		)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.resize()

		return m, nil

	case contexts.ActivatedMsg:
		if msg.Err != nil {
			m.status = fmt.Sprintf("context activation failed: %v", msg.Err)
			return m, nil
		}

		m.context = msg.Context
		m.status = fmt.Sprintf("connected to %s", msg.Context.Cloud)

		cmd := m.autoRefreshResource()

		return m, tea.Batch(
			clearStatusCmd(m.status),
			cmd,
		)

	case clearStatusMsg:
		if m.status == msg.status {
			m.status = ""
		}

		return m, nil

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
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *Model) executeDefaultCommand() tea.Cmd {
	if m.resource == nil {
		return nil
	}

	for _, command := range m.resource.Commands() {
		if command.Default {
			return m.executeResourceCommand(command.Key)
		}
	}

	return nil
}

func (m *Model) executeResourceCommand(key string) tea.Cmd {
	if m.resource == nil {
		return nil
	}

	cursor := m.table.Cursor()
	if cursor < 0 || cursor >= len(m.resourceRows) {
		return nil
	}

	for _, command := range m.resource.Commands() {
		if command.Key == key {
			return m.resource.Execute(command, m.resourceRows[cursor])
		}
	}

	return nil
}

func (m *Model) switchResource(name string) tea.Cmd {
	m.navigation = nil
	m.navigateID = ""

	r, ok := m.registry.Get(name)
	if !ok {
		m.status = fmt.Sprintf("unknown resource: %s", name)
		return nil
	}

	m.resource = r
	m.err = nil
	m.status = ""

	m.loading = true
	m.showLoading = true
	m.loaded = false
	m.itemCount = 0
	m.loadID++

	// A resource switch must not display rows belonging to the
	// previously selected resource.
	m.table.SetRows(nil)
	m.table.SetCursor(0)
	m.resourceRows = nil

	m.resize()

	return m.loadResource()
}

func (m *Model) navigateResource(name string) tea.Cmd {
	navigation := m.navigation
	navigateID := m.navigateID

	cmd := m.switchResource(name)

	m.navigation = navigation
	m.navigateID = navigateID

	return cmd
}

func (m *Model) navigateBack() tea.Cmd {
	if len(m.navigation) == 0 {
		return nil
	}

	last := len(m.navigation) - 1
	entry := m.navigation[last]

	m.navigation = m.navigation[:last]
	m.navigateID = entry.id

	navigation := m.navigation
	navigateID := m.navigateID

	cmd := m.switchResource(entry.resource)

	m.navigation = navigation
	m.navigateID = navigateID

	return cmd
}

func (m *Model) refreshResource() tea.Cmd {
	if m.resource == nil || m.loading {
		return nil
	}

	m.err = nil
	m.status = ""

	// Manual refresh keeps the existing rows visible while the
	// request runs, but gives the user visual feedback.
	m.loading = true
	m.showLoading = true
	m.loadID++

	return m.loadResource()
}

func (m *Model) autoRefreshResource() tea.Cmd {
	if m.resource == nil || m.loading {
		return nil
	}

	// Background refresh deliberately leaves the existing table
	// untouched and does not show the loading badge.
	m.loading = true
	m.showLoading = false
	m.loadID++

	return m.loadResource()
}

func (m Model) loadResource() tea.Cmd {
	loadID := m.loadID
	r := m.resource

	selectedID := ""
	cursor := m.table.Cursor()

	if cursor >= 0 && cursor < len(m.resourceRows) {
		selectedID = m.resourceRows[cursor].ID
	}

	return func() tea.Msg {
		if r == nil {
			return resourcesLoadedMsg{
				loadID:     loadID,
				selectedID: selectedID,
				cursor:     cursor,
				err:        fmt.Errorf("no active resource"),
			}
		}

		rows, err := r.List(context.Background())

		return resourcesLoadedMsg{
			loadID:     loadID,
			rows:       rows,
			selectedID: selectedID,
			cursor:     cursor,
			err:        err,
		}
	}
}

func autoRefreshCmd() tea.Cmd {
	return tea.Tick(
		refreshInterval,
		func(time.Time) tea.Msg {
			return autoRefreshMsg{}
		},
	)
}

func (m *Model) resize() {
	const (
		headerHeight         = 7
		footerHeight         = 2
		tableContainerBorder = 2
	)

	tableHeight := max(
		m.height-
			headerHeight-
			footerHeight-
			tableContainerBorder,
		1,
	)

	m.table.SetHeight(tableHeight)

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

func (m Model) renderTable() string {
	title := ""

	if m.resource != nil {
		title = fmt.Sprintf(
			" %s[%d] | last update: %s ",
			m.resource.Kind(),
			m.itemCount,
			time.Now().Format("02.01.2006 15:04:05"),
		)
	}

	innerWidth := max(m.width-2, 1)

	titleWidth := lipgloss.Width(title)
	remaining := max(innerWidth-titleWidth, 0)

	left := remaining / 2
	right := remaining - left

	borderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00BFFF"))

	top := borderStyle.Render(
		"┌" +
			repeat("─", left) +
			title +
			repeat("─", right) +
			"┐",
	)

	bottom := borderStyle.Render(
		"└" +
			repeat("─", innerWidth) +
			"┘",
	)

	body := lipgloss.NewStyle().
		Width(innerWidth).
		Render(m.table.View())

	// A successfully loaded resource with zero items is not an
	// error. Keep the table header and render an explicit empty
	// state in the row area.
	if m.loaded && !m.loading && m.itemCount == 0 {
		lines := strings.Split(body, "\n")

		if len(lines) > 1 {
			header := lines[0]
			contentHeight := len(lines) - 1

			emptyContent := lipgloss.Place(
				innerWidth,
				contentHeight,
				lipgloss.Center,
				lipgloss.Center,
				"No resources found",
			)

			body = header + "\n" + emptyContent
		}
	}

	bodyLines := strings.Split(body, "\n")

	for i, line := range bodyLines {
		lineWidth := lipgloss.Width(line)

		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}

		bodyLines[i] =
			borderStyle.Render("│") +
				line +
				borderStyle.Render("│")
	}

	return top +
		"\n" +
		strings.Join(bodyLines, "\n") +
		"\n" +
		bottom
}

func repeat(s string, count int) string {
	return strings.Repeat(s, max(count, 0))
}

func clearStatusCmd(status string) tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearStatusMsg{status: status}
	})
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := m.renderHeader()
	tableView := m.renderTable()

	resourceLine := ""

	if m.resource != nil {
		var resourceTags strings.Builder

		for _, entry := range m.navigation {
			resourceTags.WriteString(navigationTagStyle.Render(" <"+entry.resource+"> ") + " ")
		}
		resourceTags.WriteString(resourceTagStyle.Render(" <" + m.resource.Kind() + "> "))

		//resourceTag := resourceTagStyle.Render("< " + m.resource.Kind() + " >")

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
		tableView +
		"\n" +
		resourceLine +
		"\n" +
		commandLine
}

func commandKey(msg tea.KeyMsg) string {
	key := msg.String()

	if len(key) == 1 && key[0] >= 'A' && key[0] <= 'Z' {
		return "shift-" + strings.ToLower(key)
	}

	return key
}

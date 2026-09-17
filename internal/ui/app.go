package ui

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resourcesLoadedMsg struct {
	loadID uint64
	rows   []resource.Row
	err    error
}

type Model struct {
	registry *resource.Registry
	resource resource.Resource
	table    table.Model
	err      error

	itemCount int
	status    string
	loading   bool
	loadID    uint64

	width  int
	height int

	profile string
	region  string
	project string
	domain  string

	commandMode bool
	command     textinput.Model
}

func New(registry *resource.Registry) Model {
	r, ok := registry.Get("servers")
	if !ok {
		return Model{
			registry: registry,
			err:      fmt.Errorf("servers resource not registered"),
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

	styles := table.DefaultStyles()
	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("255")).
		Background(lipgloss.Color("#ED1944")).
		Bold(false)

	t.SetStyles(styles)

	command := textinput.New()
	command.Prompt = ":"
	command.CharLimit = 64

	return Model{
		registry: registry,
		resource: r,
		table:    t,
		command:  command,

		loading: true,
		loadID:  1,

		profile: "default",
		region:  "RegionOne",
		project: "default",
		domain:  "Default",
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadResource()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Esc always means "cancel/dismiss and return to normal mode".
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		m.status = ""
		m.commandMode = false
		m.command.Blur()
		m.command.SetValue("")

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
		}

	case resourcesLoadedMsg:
		if msg.loadID != m.loadID {
			return m, nil
		}

		m.loading = false

		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		if m.resource == nil {
			m.err = fmt.Errorf("no active resource")
			return m, nil
		}

		columns := m.resource.Columns()
		rows := make([]table.Row, 0, len(msg.rows))

		for _, row := range msg.rows {
			values := make(table.Row, 0, len(columns))

			for _, column := range columns {
				values = append(values, row.Fields[column.Key])
			}

			rows = append(rows, values)
		}

		m.table.SetRows(rows)
		m.table.SetCursor(0)
		m.itemCount = len(rows)

		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		m.resize()

		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m *Model) switchResource(name string) tea.Cmd {
	r, ok := m.registry.Get(name)
	if !ok {
		m.status = fmt.Sprintf("unknown resource: %s", name)
		return nil
	}

	m.resource = r
	m.err = nil
	m.status = ""
	m.loading = true
	m.itemCount = 0
	m.loadID++

	// Immediately remove the previous resource's data.
	m.table.SetRows(nil)
	m.table.SetCursor(0)

	m.resize()

	return m.loadResource()
}

func (m Model) loadResource() tea.Cmd {
	loadID := m.loadID
	r := m.resource

	return func() tea.Msg {
		if r == nil {
			return resourcesLoadedMsg{
				loadID: loadID,
				err:    fmt.Errorf("no active resource"),
			}
		}

		rows, err := r.List(context.Background())

		return resourcesLoadedMsg{
			loadID: loadID,
			rows:   rows,
			err:    err,
		}
	}
}

func (m *Model) resize() {
	const (
		headerHeight         = 6
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
		title = fmt.Sprintf(" %s[%d] ", m.resource.Kind(), m.itemCount)
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

	bodyLines := strings.Split(body, "\n")

	for i, line := range bodyLines {
		lineWidth := lipgloss.Width(line)

		if lineWidth < innerWidth {
			line += strings.Repeat(" ", innerWidth-lineWidth)
		}

		bodyLines[i] = borderStyle.Render("│") +
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

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := m.renderHeader()
	tableView := m.renderTable()

	resourceLine := ""

	if m.resource != nil {
		resourceTag := resourceTagStyle.Render(
			"<" + m.resource.Kind() + ">",
		)

		loadingTag := ""

		if m.loading {
			loadingTag = loadingStyle.Render(" Loading... ")
		}

		resourceLine = lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().
				Width(max(m.width-lipgloss.Width(loadingTag), 1)).
				Render(resourceTag),
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

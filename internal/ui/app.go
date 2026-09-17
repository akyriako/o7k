package ui

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resourcesLoadedMsg struct {
	rows []resource.Row
	err  error
}

type Model struct {
	registry *resource.Registry
	resource resource.Resource
	table    table.Model
	err      error
	status   string
	width    int
	height   int

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

	m.resize()

	return m.loadResource()
}

func (m Model) loadResource() tea.Cmd {
	return func() tea.Msg {
		if m.resource == nil {
			return resourcesLoadedMsg{
				err: fmt.Errorf("no active resource"),
			}
		}

		rows, err := m.resource.List(context.Background())

		return resourcesLoadedMsg{
			rows: rows,
			err:  err,
		}
	}
}

func (m *Model) resize() {
	const (
		headerHeight = 2
		footerHeight = 2
	)

	tableHeight := max(
		m.height-headerHeight-footerHeight,
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

	extra := max(m.width-minWidth, 0)

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

var (
	resourceTagStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("#ED1944")).
				Bold(true)

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#87CEFA"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFA500")).
			Align(lipgloss.Center)
)

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := "o7k — OpenStack TUI"

	resourceLine := ""
	if m.resource != nil {
		resourceLine = resourceTagStyle.Render(
			"<" + m.resource.Kind() + ">",
		)
	}

	commandLine := commandStyle.Render("<q> quit")

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
		m.table.View() +
		"\n" +
		resourceLine +
		"\n" +
		commandLine
}

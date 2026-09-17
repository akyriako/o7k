package ui

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resourcesLoadedMsg struct {
	rows []resource.Row
	err  error
}

type Model struct {
	registry *resource.Registry
	table    table.Model
	err      error
	width    int
	height   int
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
			Width: column.Width,
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
		Background(lipgloss.Color("240")).
		Bold(false)

	t.SetStyles(styles)

	return Model{
		registry: registry,
		table:    t,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadServers()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case resourcesLoadedMsg:
		r, ok := m.registry.Get("servers")
		if !ok {
			m.err = fmt.Errorf("servers resource not registered")
			return m, nil
		}

		columns := r.Columns()
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

func (m Model) loadServers() tea.Cmd {
	return func() tea.Msg {
		r, ok := m.registry.Get("servers")
		if !ok {
			return resourcesLoadedMsg{
				err: fmt.Errorf("servers resource not registered"),
			}
		}

		rows, err := r.List(context.Background())

		return resourcesLoadedMsg{
			rows: rows,
			err:  err,
		}
	}
}

func (m *Model) resize() {
	const (
		headerHeight = 2
		footerHeight = 1
	)

	tableHeight := max(m.height-headerHeight-footerHeight, 1)

	m.table.SetHeight(tableHeight)
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	header := "o7k — OpenStack TUI\n"
	footer := "q quit"

	return header +
		"\n" +
		m.table.View() +
		"\n" +
		footer
}

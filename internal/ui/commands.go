package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type resourcesLoadedMsg struct {
	loadID     uint64
	rows       []resource.Row
	selectedID string
	cursor     int
	//err        error
}

type errMsg struct {
	err    error
	op     string
	loadID uint64
}

type clearStatusMsg struct {
	status string
}

type clipboardResultMsg struct {
	//err error
}

type resourceScope map[string]string

type navigationEntry struct {
	resource string
	id       string
	filter   *resourceFilter
	scope    resourceScope
}

type resourceFilter struct {
	field string
	value string
}

type autoRefreshMsg struct{}

func commandKey(msg tea.KeyMsg) string {
	key := msg.String()

	if len(key) == 1 && key[0] >= 'A' && key[0] <= 'Z' {
		return "shift-" + strings.ToLower(key)
	}

	return key
}

func clearStatus(status string) tea.Cmd {
	return tea.Tick(3*time.Second, func(time.Time) tea.Msg {
		return clearStatusMsg{status: status}
	})
}

func autoRefresh() tea.Cmd {
	return tea.Tick(
		refreshInterval,
		func(time.Time) tea.Msg {
			return autoRefreshMsg{}
		},
	)
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
			row := m.resourceRows[cursor]

			if m.resource.Kind() == "contexts" && key == "a" {
				m.activatingContext = true
				m.showLoading = true
				m.loadingLabel = "Connecting to " + row.ID + "..."
			}

			if key == "s" {
				m.showLoading = true
				m.loadingLabel = "Loading..."
			}

			return m.resource.Execute(command, row)
		}
	}

	return nil
}

func (m *Model) switchResource(name string, scope resourceScope) tea.Cmd {
	m.navigation = nil
	m.navigateID = ""
	m.filter = nil
	m.scope = scope

	m.detailMode = false
	m.detailID = ""
	m.detailContent = ""
	m.autoRefreshPaused = false

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
	m.loadingLabel = "Loading..."
	m.loaded = false
	m.itemCount = 0
	m.loadID++

	m.table.SetRows(nil)
	m.table.SetCursor(0)
	m.resourceRows = nil
	m.tableXOffset = 0

	m.Resize()

	return m.loadResource()
}

func (m *Model) navigateResource(name string) tea.Cmd {
	navigation := m.navigation
	navigateID := m.navigateID

	cmd := m.switchResource(name, nil)

	m.navigation = navigation
	m.navigateID = navigateID

	return cmd
}

func (m *Model) navigateFilteredResource(name string) tea.Cmd {
	navigation := m.navigation
	filter := m.filter

	cmd := m.switchResource(name, nil)

	m.navigation = navigation
	m.filter = filter

	return cmd
}

func (m *Model) navigateScopedResource(name string) tea.Cmd {
	navigation := m.navigation
	scope := m.scope

	cmd := m.switchResource(name, scope)

	m.navigation = navigation
	m.scope = scope

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
	m.filter = entry.filter
	m.scope = entry.scope

	navigation := m.navigation
	navigateID := m.navigateID
	filter := m.filter
	scope := m.scope

	cmd := m.switchResource(entry.resource, entry.scope)

	m.navigation = navigation
	m.navigateID = navigateID
	m.filter = filter
	m.scope = scope

	return cmd
}

func filterResourceRows(rows []resource.Row, filter *resourceFilter) []resource.Row {
	if filter == nil {
		return rows
	}

	filtered := make([]resource.Row, 0, len(rows))

	for _, row := range rows {
		if row.Fields[filter.field] == filter.value {
			filtered = append(filtered, row)
		}
	}

	return filtered
}

func (m *Model) refreshResource() tea.Cmd {
	if m.resource == nil || m.loading {
		return nil
	}

	m.err = nil
	m.status = ""

	m.loading = true
	m.showLoading = true
	m.loadingLabel = "Loading..."
	m.loadID++

	return m.loadResource()
}

func (m *Model) autoRefreshResource() tea.Cmd {
	if m.resource == nil ||
		m.loading ||
		m.autoRefreshPaused ||
		m.detailMode {
		return nil
	}

	m.loading = true
	m.showLoading = true
	m.loadingLabel = "Refreshing..."
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
			}
		}

		ctx := context.Background()
		if m.scope != nil {
			ctx = resource.WithScope(ctx, m.scope)
		}
		rows, err := r.List(ctx)
		if err != nil {
			return errMsg{
				err:    err,
				loadID: loadID,
			}
		}

		return resourcesLoadedMsg{
			loadID:     loadID,
			rows:       rows,
			selectedID: selectedID,
			cursor:     cursor,
		}
	}
}

func (m *Model) copyDetailsViewport() tea.Cmd {
	if m.detailContent == "" {
		return nil
	}

	content := m.detailContent

	return func() tea.Msg {
		if err := clipboard.Init(); err != nil {
			return errMsg{err: fmt.Errorf("initializing clipboard failed: %v", err)}
		}

		_, err := clipboard.Write(
			context.Background(),
			clipboard.FmtText,
			[]byte(content),
		)
		if err != nil {
			return errMsg{err: fmt.Errorf("writing to clipboard failed: %v", err)}
		}

		return clipboardResultMsg{}
	}
}

func (m *Model) copyTableViewport() tea.Cmd {
	if m.detailMode {
		return nil
	}

	var content strings.Builder
	columns := m.table.Columns()
	headers := make([]string, len(columns))

	for i, column := range columns {
		headers[i] = column.Title
	}

	content.WriteString(strings.Join(headers, "\t"))
	content.WriteByte('\n')

	for _, row := range m.table.Rows() {
		content.WriteString(strings.Join(row, "\t"))
		content.WriteByte('\n')
	}

	data := content.String()

	return func() tea.Msg {
		if err := clipboard.Init(); err != nil {
			return errMsg{err: fmt.Errorf("initializing clipboard failed: %v", err)}
		}

		_, err := clipboard.Write(
			context.Background(),
			clipboard.FmtText,
			[]byte(data),
		)
		if err != nil {
			return errMsg{err: fmt.Errorf("writing to clipboard failed: %v", err)}
		}

		return clipboardResultMsg{}
	}
}

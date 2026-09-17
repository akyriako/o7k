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
	err        error
}

type clearStatusMsg struct {
	status string
}

type clipboardResultMsg struct {
	err error
}

type navigationEntry struct {
	resource string
	id       string
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
			return m.resource.Execute(command, m.resourceRows[cursor])
		}
	}

	return nil
}

func (m *Model) switchResource(name string) tea.Cmd {
	m.navigation = nil
	m.navigateID = ""
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
	m.loaded = false
	m.itemCount = 0
	m.loadID++

	m.table.SetRows(nil)
	m.table.SetCursor(0)
	m.resourceRows = nil

	m.Resize()

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

	m.loading = true
	m.showLoading = true
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

func (m *Model) copyViewport() tea.Cmd {
	if !m.detailMode || m.detailContent == "" {
		return nil
	}

	content := m.detailContent

	return func() tea.Msg {
		if err := clipboard.Init(); err != nil {
			return clipboardResultMsg{err: err}
		}

		clipboard.Write(
			context.Background(),
			clipboard.FmtText,
			[]byte(content),
		)

		return clipboardResultMsg{}
	}
}

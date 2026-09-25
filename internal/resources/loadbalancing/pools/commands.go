package pools

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) members(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "members",
			Scope: map[string]string{
				"pool_id": row.ID,
			},
		}
	}
}

func (r *Resource) healthMonitor(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "healthmonitors",
			Field:    "id",
			Value:    row.Fields["healthmonitor_id"],
		}
	}
}

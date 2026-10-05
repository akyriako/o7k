package projects

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) quotas(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "quotas",
			Scope: map[string]string{
				"project_id": row.ID,
			},
		}
	}
}

package endpoints

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToService(row resource.Row) tea.Cmd {
	serviceID := row.Fields["service_id"]

	return func() tea.Msg {
		return resource.NavigateMsg{
			Resource: "services",
			ID:       serviceID,
		}
	}
}

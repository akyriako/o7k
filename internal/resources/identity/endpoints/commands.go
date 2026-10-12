package endpoints

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) service(row resource.Row) tea.Cmd {
	serviceID := row.Fields["service_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "services",
			Field:    "id",
			Value:    serviceID,
		}
	}
}

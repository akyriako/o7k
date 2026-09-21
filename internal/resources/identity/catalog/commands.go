package catalog

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToService(row resource.Row) tea.Cmd {
	serviceName := row.Fields["service"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "services",
			Field:    "name",
			Value:    serviceName,
		}
	}
}

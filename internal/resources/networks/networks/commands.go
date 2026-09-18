package networks

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToSubnets(row resource.Row) tea.Cmd {
	networkID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "subnets",
			Field:    "network_id",
			Value:    networkID,
		}
	}
}

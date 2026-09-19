package ports

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToNetwork(row resource.Row) tea.Cmd {
	networkID := row.Fields["network_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "networks",
			Field:    "id",
			Value:    networkID,
		}
	}
}

func (r *Resource) navigateToFloatingIPs(row resource.Row) tea.Cmd {
	portID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "floatingips",
			Field:    "port_id",
			Value:    portID,
		}
	}
}

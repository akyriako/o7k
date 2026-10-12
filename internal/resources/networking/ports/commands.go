package ports

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) network(row resource.Row) tea.Cmd {
	networkID := row.Fields["network_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "networks",
			Field:    "id",
			Value:    networkID,
		}
	}
}

func (r *Resource) floatingips(row resource.Row) tea.Cmd {
	portID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "floatingips",
			Field:    "port_id",
			Value:    portID,
		}
	}
}

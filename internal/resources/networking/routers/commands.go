package routers

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToPorts(row resource.Row) tea.Cmd {
	routerID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "ports",
			Field:    "device_id",
			Value:    routerID,
		}
	}
}

func (r *Resource) navigateToFloatingIPs(row resource.Row) tea.Cmd {
	routerID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "floatingips",
			Field:    "router_id",
			Value:    routerID,
		}
	}
}

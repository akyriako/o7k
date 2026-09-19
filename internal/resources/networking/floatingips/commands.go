package floatingips

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToRouter(row resource.Row) tea.Cmd {
	routerID := row.Fields["router_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "routers",
			Field:    "id",
			Value:    routerID,
		}
	}
}

func (r *Resource) navigateToPort(row resource.Row) tea.Cmd {
	portID := row.Fields["port_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "ports",
			Field:    "id",
			Value:    portID,
		}
	}
}

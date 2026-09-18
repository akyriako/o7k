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

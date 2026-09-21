package services

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToEndpoints(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "endpoints",
			Field:    "service_id",
			Value:    row.ID,
		}
	}
}

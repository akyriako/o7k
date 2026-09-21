package recordsets

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToZone(row resource.Row) tea.Cmd {
	zoneID := row.Fields["zone_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "zones",
			Field:    "id",
			Value:    zoneID,
		}
	}
}

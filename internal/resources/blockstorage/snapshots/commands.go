package snapshots

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToVolume(row resource.Row) tea.Cmd {
	volumeID := row.Fields["volume_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "volumes",
			Field:    "id",
			Value:    volumeID,
		}
	}
}

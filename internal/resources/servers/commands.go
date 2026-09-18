package servers

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToImage(row resource.Row) tea.Cmd {
	imageID := row.Fields["image_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    imageID,
		}
	}
}

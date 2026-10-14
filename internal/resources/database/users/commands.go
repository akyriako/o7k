package users

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.DetailsMsg{
			ID:      row.ID,
			Content: row.Fields,
		}
	}
}

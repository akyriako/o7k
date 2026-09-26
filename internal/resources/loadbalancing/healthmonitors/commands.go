package healthmonitors

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) pool(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "pools",
			Field:    "healthmonitor_id",
			Value:    row.ID,
		}
	}
}

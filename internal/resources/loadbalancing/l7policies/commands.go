package l7policies

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) rules(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "l7rules",
			Scope: map[string]string{
				"l7policy_id": row.ID,
			},
		}
	}
}

func (r *Resource) listener(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "listeners",
			Field:    "id",
			Value:    row.Fields["listener_id"],
		}
	}
}

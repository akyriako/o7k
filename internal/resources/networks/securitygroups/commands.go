package securitygroups

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToRules(row resource.Row) tea.Cmd {
	securityGroupID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "securitygrouprules",
			Field:    "security_group_id",
			Value:    securityGroupID,
		}
	}
}

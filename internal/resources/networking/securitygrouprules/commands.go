package securitygrouprules

import (
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) navigateToSecurityGroup(row resource.Row) tea.Cmd {
	securityGroupID := row.Fields["security_group_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "securitygroups",
			Field:    "id",
			Value:    securityGroupID,
		}
	}
}

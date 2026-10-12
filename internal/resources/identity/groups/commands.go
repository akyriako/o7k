package groups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identitygroups "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.IdentityV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		group, err := identitygroups.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting group %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      group.ID,
			Content: group,
		}
	}
}

func (r *Resource) domain(row resource.Row) tea.Cmd {
	domainID := row.Fields["domain_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "domains",
			Field:    "id",
			Value:    domainID,
		}
	}
}

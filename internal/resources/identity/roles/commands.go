package roles

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityroles "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.IdentityV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		role, err := identityroles.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting role %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      role.ID,
			Content: role,
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

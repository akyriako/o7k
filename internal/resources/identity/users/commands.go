package users

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityusers "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.IdentityV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		user, err := identityusers.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting user %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      user.ID,
			Content: user,
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

func (r *Resource) project(row resource.Row) tea.Cmd {
	projectID := row.Fields["default_project_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "projects",
			Field:    "id",
			Value:    projectID,
		}
	}
}

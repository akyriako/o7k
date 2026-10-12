package domains

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identitydomains "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.IdentityV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		domain, err := identitydomains.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting domain %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      domain.ID,
			Content: domain,
		}
	}
}

func (r *Resource) projects(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "projects",
			Field:    "domain_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) users(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "users",
			Field:    "domain_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) groups(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "groups",
			Field:    "domain_id",
			Value:    row.ID,
		}
	}
}

package projects

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityprojects "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.IdentityV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		project, err := identityprojects.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting project %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      project.ID,
			Content: project,
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

func (r *Resource) parentProject(row resource.Row) tea.Cmd {
	parentID := row.Fields["parent_id"]

	return func() tea.Msg {
		if parentID == "" {
			return resource.ErrorMsg{
				Err: fmt.Errorf("project %q has no parent project", row.ID),
			}
		}

		return resource.NavigateFilteredMsg{
			Resource: "projects",
			Field:    "id",
			Value:    parentID,
		}
	}
}

func (r *Resource) users(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "users",
			Field:    "default_project_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) quotas(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "quotas",
			Scope: map[string]string{
				"project_id": row.ID,
			},
		}
	}
}

package stacks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackresources"
	orchestrationstacks "github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		stack, err := orchestrationstacks.Get(
			context.Background(),
			client,
			row.Fields["name"],
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting stack %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      stack.ID,
			Content: stack,
		}
	}
}

func (r *Resource) resources(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "stack-resources",
			Scope: map[string]string{
				"stack_name": row.Fields["name"],
				"stack_id":   row.ID,
			},
		}
	}
}

func (r *Resource) events(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "stack-events",
			Scope: map[string]string{
				"stack_name": row.Fields["name"],
				"stack_id":   row.ID,
			},
		}
	}
}

func (r *Resource) deployments(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		pages, err := stackresources.List(client, row.Fields["name"], row.ID, &stackresources.ListOpts{Depth: 99}).AllPages(context.Background())
		if err != nil {
			return resource.ErrorMsg{Err: fmt.Errorf("listing stack resources: %w", err)}
		}

		items, err := stackresources.ExtractResources(pages)
		if err != nil {
			return resource.ErrorMsg{Err: fmt.Errorf("extracting stack resources: %w", err)}
		}

		ids := make([]string, 0)
		for _, item := range items {
			if item.Type == "OS::Heat::SoftwareDeployment" && item.PhysicalID != "" {
				ids = append(ids, item.PhysicalID)
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "software-deployments",
			Field:    "id",
			Values:   ids,
		}
	}
}

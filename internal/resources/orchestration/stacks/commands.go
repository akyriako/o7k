package stacks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
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

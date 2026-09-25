package stackevents

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackevents"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		stackName := row.Fields["stack_name"]
		stackID := row.Fields["stack_id"]
		resourceName := row.Fields["resource"]

		event, err := stackevents.Get(
			context.Background(),
			client,
			stackName,
			stackID,
			resourceName,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting stack event %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      event.ID,
			Content: event,
		}
	}
}

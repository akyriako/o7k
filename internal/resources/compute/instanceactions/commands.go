package instanceactions

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeactions "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/instanceactions"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		action, err := computeactions.Get(
			context.Background(),
			client,
			row.Fields["server_id"],
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting instance action %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      row.ID,
			Content: action,
		}
	}
}

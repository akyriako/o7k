package flavors

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeflavors "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		flavor, err := computeflavors.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting flavor %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      flavor.ID,
			Content: flavor,
		}
	}
}

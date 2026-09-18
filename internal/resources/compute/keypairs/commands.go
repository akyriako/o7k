package keypairs

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
)

func (r *Resource) show(name string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		keypair, err := keypairs.Get(context.Background(), client, name, keypairs.GetOpts{}).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting keypair %q: %w", name, err),
			}
		}

		return resource.DetailsMsg{
			ID:      keypair.Name,
			Content: keypair,
		}
	}
}

package sharenetworks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasharenetworks "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		network, err := manilasharenetworks.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share network %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      network.ID,
			Content: network,
		}
	}
}

func (r *Resource) network(row resource.Row) tea.Cmd {
	networkID := row.Fields["neutron_net_id"]

	return func() tea.Msg {
		//if networkID == "" {
		//	return resource.ErrorMsg{
		//		Err: fmt.Errorf("share network %q has no neutron network", row.ID),
		//	}
		//}

		return resource.NavigateFilteredMsg{
			Resource: "networks",
			Field:    "id",
			Value:    networkID,
		}
	}
}

func (r *Resource) subnet(row resource.Row) tea.Cmd {
	subnetID := row.Fields["neutron_subnet_id"]

	return func() tea.Msg {
		//if subnetID == "" {
		//	return resource.ErrorMsg{
		//		Err: fmt.Errorf("share network %q has no neutron subnet", row.ID),
		//	}
		//}

		return resource.NavigateFilteredMsg{
			Resource: "subnets",
			Field:    "id",
			Value:    subnetID,
		}
	}
}

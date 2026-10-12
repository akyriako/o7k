package networks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
)

func (r *Resource) subnets(row resource.Row) tea.Cmd {
	networkID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "subnets",
			Field:    "network_id",
			Value:    networkID,
		}
	}
}

func (r *Resource) ports(row resource.Row) tea.Cmd {
	networkID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "ports",
			Field:    "network_id",
			Value:    networkID,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.NetworkV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		network, err := networks.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting network %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      network.ID,
			Content: network,
		}
	}
}

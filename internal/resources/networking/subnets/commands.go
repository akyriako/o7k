package subnets

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

func (r *Resource) network(row resource.Row) tea.Cmd {
	networkID := row.Fields["network_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "networks",
			Field:    "id",
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

		subnet, err := subnets.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting subnet %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      subnet.ID,
			Content: subnet,
		}
	}
}

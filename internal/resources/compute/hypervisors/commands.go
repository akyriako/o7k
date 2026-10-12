package hypervisors

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		hypervisor, err := hypervisors.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting hypervisor %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      hypervisor.ID,
			Content: hypervisor,
		}
	}
}

func (r *Resource) servers(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		host := row.Fields["host"]
		if host == "" {
			return resource.ErrorMsg{
				Err: fmt.Errorf("hypervisor %q has no compute host", row.ID),
			}
		}

		return resource.NavigateFilteredMsg{
			Resource: "servers",
			Field:    "host",
			Value:    host,
		}
	}
}

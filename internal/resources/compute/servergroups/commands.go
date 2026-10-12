package servergroups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		serverGroup, err := servergroups.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting server group %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      serverGroup.ID,
			Content: serverGroup,
		}
	}
}

func (r *Resource) members(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		serverGroup, err := servergroups.Get(context.Background(), client, row.ID).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting server group %q: %w", row.ID, err),
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "servers",
			Field:    "id",
			Values:   serverGroup.Members,
		}
	}
}

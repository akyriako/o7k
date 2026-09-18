package servers

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeservers "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)

func (r *Resource) navigateToImage(row resource.Row) tea.Cmd {
	imageID := row.Fields["image_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    imageID,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		server, err := computeservers.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting server %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      server.ID,
			Content: server,
		}
	}
}

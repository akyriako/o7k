package snapshots

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
)

func (r *Resource) volume(row resource.Row) tea.Cmd {
	volumeID := row.Fields["volume_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "volumes",
			Field:    "id",
			Value:    volumeID,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.BlockStorageV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		snapshot, err := snapshots.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting volume snapshot %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      snapshot.ID,
			Content: snapshot,
		}
	}
}

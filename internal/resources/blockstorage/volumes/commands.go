package volumes

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.BlockStorageV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		volume, err := volumes.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting volume %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      volume.ID,
			Content: volume,
		}
	}
}

func (r *Resource) navigateToBackups(row resource.Row) tea.Cmd {
	volumeID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "backups",
			Field:    "volume_id",
			Value:    volumeID,
		}
	}
}

func (r *Resource) navigateToSnapshots(row resource.Row) tea.Cmd {
	volumeID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "snapshots",
			Field:    "volume_id",
			Value:    volumeID,
		}
	}
}

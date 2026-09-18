package backups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/backups"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.BlockStorageV3()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		backup, err := backups.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting volume backup %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      backup.ID,
			Content: backup,
		}
	}
}

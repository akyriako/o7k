package shares

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilashares "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		share, err := manilashares.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      share.ID,
			Content: share,
		}
	}
}

func (r *Resource) snapshots(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "share-snapshots",
			Field:    "share_id",
			Value:    row.ID,
		}
	}
}

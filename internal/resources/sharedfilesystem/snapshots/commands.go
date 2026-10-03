package snapshots

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasnapshots "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/snapshots"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		snapshot, err := manilasnapshots.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share snapshot %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      snapshot.ID,
			Content: snapshot,
		}
	}
}

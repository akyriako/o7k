package sharetypes

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasharetypes "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		shareType, err := manilasharetypes.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share type %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      shareType.ID,
			Content: shareType,
		}
	}
}

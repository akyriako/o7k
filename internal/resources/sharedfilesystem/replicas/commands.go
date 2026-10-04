package replicas

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilareplicas "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/replicas"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		replica, err := manilareplicas.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share replica %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      replica.ID,
			Content: replica,
		}
	}
}

package containers

import (
	"context"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	barbicancontainers "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/containers"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.KeyManagerV1()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		item, err := barbicancontainers.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		return resource.DetailsMsg{
			ID:      row.ID,
			Content: item,
		}
	}
}

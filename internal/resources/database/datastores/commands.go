package datastores

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbdatastores "github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.DBV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		datastore, err := dbdatastores.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting datastore %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      id,
			Content: datastore,
		}
	}
}

package orders

import (
	"context"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	barbicanorders "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/orders"
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

		item, err := barbicanorders.Get(
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

func (r *Resource) secret(row resource.Row) tea.Cmd {
	secretID := row.Fields["secret_id"]
	if secretID == "" {
		return nil
	}

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "secrets",
			Field:    "id",
			Value:    secretID,
		}
	}
}

func (r *Resource) container(row resource.Row) tea.Cmd {
	containerID := row.Fields["container_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "secret-containers",
			Field:    "id",
			Value:    containerID,
		}
	}
}

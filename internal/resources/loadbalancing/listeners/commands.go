package listeners

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
)

func (r *Resource) l7Policies(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "l7policies",
			Field:    "listener_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.LoadBalancerV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		listener, err := listeners.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting listener %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      listener.ID,
			Content: listener,
		}
	}
}

func (r *Resource) pools(row resource.Row) tea.Cmd {
	//return func() tea.Msg {
	//	return resource.NavigateScopedMsg{
	//		Resource: "pools",
	//		Scope: map[string]string{
	//			"loadbalancer_id": row.Fields["default_pool_id"],
	//		},
	//	}
	//}
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "pools",
			Field:    "id",
			Value:    row.Fields["default_pool_id"],
		}
	}
}

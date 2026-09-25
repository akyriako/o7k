package loadbalancers

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.LoadBalancerV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		loadbalancer, err := loadbalancers.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting load balancer %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      loadbalancer.ID,
			Content: loadbalancer,
		}
	}
}

func (r *Resource) listeners(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "listeners",
			Scope: map[string]string{
				"loadbalancer_id": row.ID,
			},
		}
	}
}

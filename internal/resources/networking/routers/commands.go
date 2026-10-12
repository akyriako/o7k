package routers

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
)

func (r *Resource) ports(row resource.Row) tea.Cmd {
	routerID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "ports",
			Field:    "device_id",
			Value:    routerID,
		}
	}
}

func (r *Resource) floatingips(row resource.Row) tea.Cmd {
	routerID := row.ID

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "floatingips",
			Field:    "router_id",
			Value:    routerID,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.NetworkV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		router, err := routers.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting router %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      router.ID,
			Content: router,
		}
	}
}

package zones

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dnszones "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
)

func (r *Resource) navigateToRecordsets(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "recordsets",
			Field:    "zone_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.DNSV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		zone, err := dnszones.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting dns zone %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      zone.ID,
			Content: zone,
		}
	}
}

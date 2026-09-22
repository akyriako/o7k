package recordsets

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dnsrecordsets "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
)

func (r *Resource) navigateToZone(row resource.Row) tea.Cmd {
	zoneID := row.Fields["zone_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "zones",
			Field:    "id",
			Value:    zoneID,
		}
	}
}

func (r *Resource) show(zoneID, rrsetID string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.DNSV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		recordset, err := dnsrecordsets.Get(context.Background(), client, zoneID, rrsetID).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting dns recordset %q: %w", rrsetID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      recordset.ID,
			Content: recordset,
		}
	}
}

package servers

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeservers "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

func (r *Resource) image(row resource.Row) tea.Cmd {
	imageID := row.Fields["image_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    imageID,
		}
	}
}

func (r *Resource) flavor(row resource.Row) tea.Cmd {
	flavor := row.Fields["flavor"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "flavors",
			Field:    "name",
			Value:    flavor,
		}
	}
}

// securityGroups violates the "no N+1 calls" rule, but this is a special case because we
// have to first load the ports and then find out the corresponding security group ids.
// this method should not be considered as the norm of jumping to resources
func (r *Resource) securityGroups(row resource.Row) tea.Cmd {
	return func() tea.Msg {

		client, err := r.context.NetworkV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		pages, err := ports.List(client, ports.ListOpts{
			DeviceID: row.ID,
		}).AllPages(context.Background())
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("listing ports for server %q: %w", row.ID, err),
			}
		}

		items, err := ports.ExtractPorts(pages)
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("extracting ports for server %q: %w", row.ID, err),
			}
		}

		uniqueSecurityGroups := make(map[string]struct{})
		securityGroupIDs := make([]string, 0)

		for _, port := range items {
			for _, securityGroupID := range port.SecurityGroups {
				if _, exists := uniqueSecurityGroups[securityGroupID]; exists {
					continue
				}

				uniqueSecurityGroups[securityGroupID] = struct{}{}
				securityGroupIDs = append(securityGroupIDs, securityGroupID)
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "securitygroups",
			Field:    "id",
			Values:   securityGroupIDs,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		server, err := computeservers.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting server %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      server.ID,
			Content: server,
		}
	}
}

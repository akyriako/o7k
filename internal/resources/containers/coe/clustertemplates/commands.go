package clustertemplates

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ContainerInfraV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		template, err := clustertemplates.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting cluster template %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      template.UUID,
			Content: template,
		}
	}
}

func (r *Resource) image(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    row.Fields["image_id"],
		}
	}
}

func (r *Resource) flavor(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "flavors",
			Field:    "name",
			Value:    row.Fields["flavor_id"],
		}
	}
}

func (r *Resource) masterFlavor(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "flavors",
			Field:    "name",
			Value:    row.Fields["master_flavor_id"],
		}
	}
}

func (r *Resource) externalNetwork(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "networks",
			Field:    "name",
			Value:    row.Fields["external_network_id"],
		}
	}
}

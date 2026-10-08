package clustertemplates

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "COE Cluster Templates"
}

func (r *Resource) Kind() string {
	return "clustertemplates"
}

func (r *Resource) Aliases() []string {
	return []string{"clustertemplate"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "coe", Title: "COE", MinWidth: 16, Flex: 1},
		{Key: "server_type", Title: "SERVER TYPE", MinWidth: 16, Flex: 1},
		{Key: "cluster_distro", Title: "DISTRO", MinWidth: 16, Flex: 1},
		{Key: "network_driver", Title: "NETWORK DRIVER", MinWidth: 18, Flex: 1},
		{Key: "public", Title: "PUBLIC", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-i", Description: "Image"},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-m", Description: "Master Flavor"},
		{Key: "shift-n", Description: "External Network"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ContainerInfraV1()
	if err != nil {
		return nil, fmt.Errorf("getting container infrastructure client: %w", err)
	}

	pages, err := clustertemplates.List(
		client,
		clustertemplates.ListOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing cluster templates: %w", err)
	}

	items, err := clustertemplates.ExtractClusterTemplates(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting cluster templates: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, template := range items {
		rows = append(rows, resource.Row{
			ID: template.UUID,
			Fields: map[string]string{
				"id":                  template.UUID,
				"name":                template.Name,
				"coe":                 template.COE,
				"server_type":         template.ServerType,
				"cluster_distro":      template.ClusterDistro,
				"network_driver":      template.NetworkDriver,
				"public":              fmt.Sprintf("%t", template.Public),
				"image_id":            template.ImageID,
				"flavor_id":           template.FlavorID,
				"master_flavor_id":    template.MasterFlavorID,
				"external_network_id": template.ExternalNetworkID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-i":
		return r.image(row)
	case "shift-f":
		return r.flavor(row)
	case "shift-m":
		return r.masterFlavor(row)
	case "shift-n":
		return r.externalNetwork(row)
	}

	return nil
}

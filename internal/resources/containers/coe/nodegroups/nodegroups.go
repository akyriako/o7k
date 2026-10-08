package nodegroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Magnum Node Groups"
}

func (r *Resource) Kind() string {
	return "nodegroups"
}

func (r *Resource) Aliases() []string {
	return []string{"nodegroup"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "role", Title: "ROLE", MinWidth: 14, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 20, Flex: 1},
		{Key: "node_count", Title: "NODE COUNT", MinWidth: 12, Flex: 0},
		{Key: "flavor_id", Title: "FLAVOR", MinWidth: 16, Flex: 1},
		{Key: "is_default", Title: "DEFAULT", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Servers"},
		{Key: "ctrl+u", Description: "Scale Up   (+1)", StatusLabel: "Scaling up"},
		{Key: "ctrl+d", Description: "Scale Down (-1)", StatusLabel: "Scaling down"},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-i", Description: "Image"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	clusterID := scope["cluster_id"]

	if clusterID == "" {
		return nil, fmt.Errorf("nodegroups requires cluster_id")
	}

	client, err := r.context.ContainerInfraV1()
	if err != nil {
		return nil, fmt.Errorf("getting container infrastructure client: %w", err)
	}

	pages, err := nodegroups.List(
		client,
		clusterID,
		nodegroups.ListOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing node groups for cluster %q: %w", clusterID, err)
	}

	items, err := nodegroups.ExtractNodeGroups(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting node groups: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, group := range items {
		rows = append(rows, resource.Row{
			ID: group.UUID,
			Fields: map[string]string{
				"id":         group.UUID,
				"name":       group.Name,
				"role":       group.Role,
				"status":     group.Status,
				"node_count": strconv.Itoa(group.NodeCount),
				"flavor_id":  group.FlavorID,
				"is_default": strconv.FormatBool(group.IsDefault),
				"image_id":   group.ImageID,
				"cluster_id": clusterID,
				"stack_id":   group.StackID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "shift-f":
		return r.flavor(row)
	case "shift-i":
		return r.image(row)
	case "shift-s":
		return r.servers(row)
	case "ctrl+u":
		return r.scale(row, 1)
	case "ctrl+d":
		return r.scale(row, -1)
	}

	return nil
}

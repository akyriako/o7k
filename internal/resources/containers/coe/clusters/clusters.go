package clusters

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "COE Clusters"
}

func (r *Resource) Kind() string {
	return "clusters"
}

func (r *Resource) Aliases() []string {
	return []string{
		"cluster",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 2},
		{Key: "status", Title: "STATUS", MinWidth: 16, Flex: 2},
		{Key: "flavor_id", Title: "FLAVOR", MinWidth: 16, Flex: 0},
		{Key: "nodes", Title: "NODES (W/M)", MinWidth: 12, Flex: 1},
		{Key: "health_status", Title: "HEALTH", MinWidth: 16, Flex: 1},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ContainerInfraV1()
	if err != nil {
		return nil, err
	}

	pages, err := clusters.List(client, clusters.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing clusters: %w", err)
	}

	items, err := clusters.ExtractClusters(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting clusters: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, cluster := range items {
		rows = append(rows, resource.Row{
			ID: cluster.UUID,
			Fields: map[string]string{
				"id":                  cluster.UUID,
				"name":                cluster.Name,
				"status":              cluster.Status,
				"health_status":       cluster.HealthStatus,
				"master_flavor_id":    cluster.MasterFlavorID,
				"flavor_id":           cluster.FlavorID,
				"node_count":          fmt.Sprintf("%d", cluster.NodeCount),
				"master_count":        fmt.Sprintf("%d", cluster.MasterCount),
				"nodes":               fmt.Sprintf("%d+%d", cluster.NodeCount, cluster.MasterCount),
				"keypair":             cluster.KeyPair,
				"cluster_template_id": cluster.ClusterTemplateID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-t", Description: "Cluster Template"},
		{Key: "shift-n", Description: "Node Groups"},
		{Key: "k", Description: "Show Kubeconfig"},
		{Key: "ctrl+k", Description: "Get Kubeconfig", StatusLabel: "Downloading"},
	}
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-t":
		return r.clusterTemplate(row)
	case "shift-n":
		return r.nodeGroups(row)
	case "k":
		return r.kubeconfig(row.ID)
	case "ctrl+k":
		return r.getKubeconfig(row)
	}

	return nil
}

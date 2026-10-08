package nodegroups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/nodegroups"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ContainerInfraV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		clusterID := row.Fields["cluster_id"]
		if clusterID == "" {
			return resource.DetailsMsg{
				Err: fmt.Errorf("node group %q has no cluster_id", row.ID),
			}
		}

		group, err := nodegroups.Get(
			context.Background(),
			client,
			clusterID,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting node group %q: %w", row.ID, err),
			}
		}

		return resource.DetailsMsg{
			ID:      group.UUID,
			Content: group,
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

func (r *Resource) image(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    row.Fields["image_id"],
		}
	}
}

func (r *Resource) servers(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ContainerInfraV1()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		clusterID := row.Fields["cluster_id"]

		cluster, err := clusters.Get(context.Background(), client, clusterID).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting cluster %q: %w", clusterID, err),
			}
		}

		result := clustertemplates.Get(context.Background(), client, cluster.ClusterTemplateID)

		if result.Err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"getting cluster template %q: %w",
					cluster.ClusterTemplateID,
					result.Err,
				),
			}
		}

		var template struct {
			Driver string `json:"driver"`
		}

		if err := result.ExtractInto(&template); err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"extracting cluster template driver: %w",
					err,
				),
			}
		}

		if template.Driver != "k8s_capi_helm_v1" {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"nova server navigation is unsupported for Magnum driver %q",
					template.Driver,
				),
			}
		}

		stackID := row.Fields["stack_id"]
		name := row.Fields["name"]
		role := row.Fields["role"]

		if stackID == "" {
			return resource.ErrorMsg{
				Err: fmt.Errorf("node group %q has no stack ID", name),
			}
		}

		var prefix string

		switch {
		case role == "master" && name == "default-master":
			prefix = stackID + "-control-plane-"

		case role == "worker":
			prefix = stackID + "-" + name + "-"

		default:
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"nova server navigation is unsupported for node group %q (role %q)",
					name,
					role,
				),
			}
		}

		return resource.NavigateScopedMsg{
			Resource: "servers",
			Scope: map[string]string{
				"name_prefix": prefix,
			},
		}
	}
}

func (r *Resource) scale(row resource.Row, delta int) tea.Cmd {
	return func() tea.Msg {
		clusterID := row.Fields["cluster_id"]
		if clusterID == "" {
			return resource.ErrorMsg{
				Err: fmt.Errorf("node group %q has no cluster ID", row.ID),
			}
		}

		client, err := r.context.ContainerInfraV1()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		cluster, err := clusters.Get(
			context.Background(),
			client,
			clusterID,
		).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting cluster %q: %w", clusterID, err),
			}
		}

		if cluster.Status != "CREATE_COMPLETE" && cluster.Status != "UPDATE_COMPLETE" {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"cannot scale node group: cluster %q is not ready (status: %s)",
					cluster.Name,
					cluster.Status,
				),
			}
		}

		if cluster.HealthStatus != "HEALTHY" {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"cannot scale node group: cluster %q is not healthy (health: %s)",
					cluster.Name,
					cluster.HealthStatus,
				),
			}
		}

		nodeGroup, err := nodegroups.Get(context.Background(), client, clusterID, row.ID).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting node group %q: %w", row.ID, err),
			}
		}

		nodeCount := nodeGroup.NodeCount + delta

		if nodeCount < 0 {
			return resource.ErrorMsg{
				Err: fmt.Errorf("node group cannot be scaled below 0 nodes"),
			}
		}

		if nodeCount < nodeGroup.MinNodeCount {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"node group %q cannot be scaled below %d nodes",
					nodeGroup.Name,
					nodeGroup.MinNodeCount,
				),
			}
		}

		if nodeGroup.MaxNodeCount != nil && nodeCount > *nodeGroup.MaxNodeCount {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"node group %q cannot be scaled above %d nodes",
					nodeGroup.Name,
					nodeGroup.MaxNodeCount,
				),
			}
		}

		err = clusters.Resize(context.Background(), client, clusterID, clusters.ResizeOpts{
			NodeCount: &nodeCount,
			NodeGroup: row.ID,
		}).Err
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"scaling node group %q from %d to %d nodes: %w",
					nodeGroup.Name,
					nodeGroup.NodeCount,
					nodeCount,
					err,
				),
			}
		}

		return resource.CommandCompletedMsg{}
	}
}

package replicas

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilareplicas "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/replicas"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Replicas"
}

func (r *Resource) Kind() string {
	return "share-replicas"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-replica",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "share_id", Title: "SHARE ID", MinWidth: 40, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "state", Title: "REPLICA STATE", MinWidth: 16, Flex: 0},
		{Key: "availability_zone", Title: "AZ", MinWidth: 12, Flex: 0},
		{Key: "host", Title: "HOST", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	pages, err := manilareplicas.ListDetail(client, manilareplicas.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing share replicas: %w", err)
	}

	items, err := manilareplicas.ExtractReplicas(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting share replicas: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, replica := range items {
		rows = append(rows, resource.Row{
			ID: replica.ID,
			Fields: map[string]string{
				"id":                replica.ID,
				"share_id":          replica.ShareID,
				"status":            replica.Status,
				"state":             replica.State,
				"availability_zone": replica.AvailabilityZone,
				"host":              replica.Host,
				"share_network_id":  replica.ShareNetworkID,
				"share_server_id":   replica.ShareServerID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	}

	return nil
}

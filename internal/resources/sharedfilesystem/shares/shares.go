package shares

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilashares "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Shares"
}

func (r *Resource) Kind() string {
	return "shares"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "size", Title: "SIZE", MinWidth: 8, Flex: 0},
		{Key: "protocol", Title: "PROTOCOL", MinWidth: 10, Flex: 0},
		{Key: "availability_zone", Title: "AVAILABILITY ZONE", MinWidth: 12, Flex: 0},
		{Key: "share_type", Title: "SHARE TYPE", MinWidth: 16, Flex: 1},
		{Key: "has_replicas", Title: "REPLICAS", MinWidth: 10, Flex: 0},
		{Key: "snapshot_support", Title: "SNAPSHOTS", MinWidth: 10, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},

		//{Key: "host", Title: "HOST", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Snapshots"},
		{Key: "shift-r", Description: "Replicas"},
		{Key: "shift-n", Description: "Share Network"},
		{Key: "shift-a", Description: "Access Rules"},
		{Key: "shift-t", Description: "Share Type"},
		{Key: "shift-e", Description: "Export Locations"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	pages, err := manilashares.ListDetail(client, manilashares.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing shares: %w", err)
	}

	items, err := manilashares.ExtractShares(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting shares: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, share := range items {
		rows = append(rows, resource.Row{
			ID: share.ID,
			Fields: map[string]string{
				"id":                share.ID,
				"name":              share.Name,
				"size":              strconv.Itoa(share.Size),
				"protocol":          share.ShareProto,
				"status":            share.Status,
				"share_type":        share.ShareTypeName,
				"share_type_id":     share.ShareType,
				"share_network_id":  share.ShareNetworkID,
				"share_server_id":   share.ShareServerID,
				"host":              share.Host,
				"availability_zone": share.AvailabilityZone,
				"has_replicas":      strconv.FormatBool(share.HasReplicas),
				"snapshot_support":  strconv.FormatBool(share.SnapshotSupport),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-s":
		return r.snapshots(row)
	case "shift-r":
		return r.replicas(row)
	case "shift-n":
		return r.shareNetwork(row)
	case "shift-a":
		return r.accessRules(row)
	case "shift-t":
		return r.shareType(row)
	case "shift-e":
		return r.exportLocations(row)
	}

	return nil
}

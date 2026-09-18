package volumes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "volumes"
}

func (r *Resource) Title() string {
	return "Volumes"
}

func (r *Resource) Aliases() []string {
	return []string{"volume", "vol"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "size", Title: "SIZE (GB)", MinWidth: 12, Flex: 0},
		{Key: "type", Title: "TYPE", MinWidth: 20, Flex: 1},
		{Key: "availability_zone", Title: "AVAILABILITY ZONE", MinWidth: 20, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Snapshots"},
		{Key: "shift-b", Description: "Backups"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.BlockStorageV3()
	if err != nil {
		return nil, err
	}

	pages, err := volumes.List(client, volumes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing volumes: %w", err)
	}

	items, err := volumes.ExtractVolumes(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting volumes: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, volume := range items {
		rows = append(rows, resource.Row{
			ID: volume.ID,
			Fields: map[string]string{
				"id":                volume.ID,
				"name":              volume.Name,
				"status":            volume.Status,
				"size":              strconv.Itoa(volume.Size),
				"type":              volume.VolumeType,
				"availability_zone": volume.AvailabilityZone,
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
		return r.navigateToSnapshots(row)
	case "shift-b":
		return r.navigateToBackups(row)
	}

	return nil
}

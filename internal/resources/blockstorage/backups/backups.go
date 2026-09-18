package backups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/backups"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "backups"
}

func (r *Resource) Title() string {
	return "Volume Backups"
}

func (r *Resource) Aliases() []string {
	return []string{"backup", "volume-backups", "volume-backup"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "size", Title: "SIZE (GB)", MinWidth: 12, Flex: 0},
		{Key: "volume_id", Title: "VOLUME ID", MinWidth: 40, Flex: 0},
		//{Key: "availability_zone", Title: "AVAILABILITY ZONE", MinWidth: 20, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.BlockStorageV3()
	if err != nil {
		return nil, err
	}

	pages, err := backups.List(client, backups.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing volume backups: %w", err)
	}

	items, err := backups.ExtractBackups(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting volume backups: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, backup := range items {
		rows = append(rows, resource.Row{
			ID: backup.ID,
			Fields: map[string]string{
				"id":        backup.ID,
				"name":      backup.Name,
				"status":    backup.Status,
				"size":      strconv.Itoa(backup.Size),
				"volume_id": backup.VolumeID,
				//"availability_zone": availabilityZone,
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

package snapshots

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "volume-snapshots"
}

func (r *Resource) Title() string {
	return "Cinder Snapshots"
}

func (r *Resource) Aliases() []string {
	return []string{"volume-snapshot"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "size", Title: "SIZE (GB)", MinWidth: 12, Flex: 0},
		//{Key: "volume_id", Title: "VOLUME ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-v", Description: "Volume"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.BlockStorageV3()
	if err != nil {
		return nil, err
	}

	pages, err := snapshots.List(client, snapshots.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing snapshots: %w", err)
	}

	items, err := snapshots.ExtractSnapshots(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting snapshots: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, snapshot := range items {
		rows = append(rows, resource.Row{
			ID: snapshot.ID,
			Fields: map[string]string{
				"id":        snapshot.ID,
				"name":      snapshot.Name,
				"status":    snapshot.Status,
				"size":      strconv.Itoa(snapshot.Size),
				"volume_id": snapshot.VolumeID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-v":
		return r.volume(row)
	}

	return nil
}

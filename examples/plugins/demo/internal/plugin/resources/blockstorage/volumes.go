package blockstorage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/v3/volumes"
)

type EvsVolumes struct {
	plugin *plugin.Plugin
}

func NewVolumes(p *plugin.Plugin) *EvsVolumes {
	return &EvsVolumes{
		plugin: p,
	}
}

func (r *EvsVolumes) Service() string {
	return "opentelekomcloud-block-storage"
}

func (r *EvsVolumes) Kind() string {
	return "evs-volumes"
}

func (r *EvsVolumes) Title() string {
	return "EVS Volumes"
}

func (r *EvsVolumes) Aliases() []string {
	return nil
}

func (r *EvsVolumes) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "Name", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "Status", MinWidth: 12},
		{Key: "size", Title: "Size", MinWidth: 10},
	}
}

func (r *EvsVolumes) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

//func (r *EvsVolumes) List(ctx context.Context) ([]pluginsdk.Row, error) {
//	client, err := r.plugin.BlockStorageV3(ctx)
//	if err != nil {
//		return nil, fmt.Errorf("getting block storage client: %w", err)
//	}
//
//	// Try direct call without pagination
//	result, err := volumes.List(client, volumes.ListOpts{})
//	if err != nil {
//		return nil, fmt.Errorf("listing volumes: %w", err)
//	}
//
//	// Check what result actually is and what methods it has
//	// You might need: result.Extract(), result.Body, or something else
//
//	rows := make([]pluginsdk.Row, 0)
//	return rows, nil
//}

func (r *EvsVolumes) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.BlockStorageV3(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting block storage client: %w", err)
	}

	result, err := volumes.List(client, volumes.ListOpts{})
	if err != nil {
		return nil, fmt.Errorf("listing volumes: %w", err)
	}

	rows := make([]pluginsdk.Row, len(result.Volumes))

	for i, volume := range result.Volumes {
		serverID := ""

		if len(volume.Attachments) > 0 {
			serverID = volume.Attachments[0].ServerID
		}

		rows[i] = pluginsdk.Row{
			ID: volume.ID,
			Fields: map[string]string{
				"id":        volume.ID,
				"name":      volume.Name,
				"status":    volume.Status,
				"size":      strconv.Itoa(volume.Size),
				"server_id": serverID,
			},
		}
	}

	return rows, nil
}

func (r *EvsVolumes) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row.ID)
	}

	return pluginsdk.Result{}, nil
}

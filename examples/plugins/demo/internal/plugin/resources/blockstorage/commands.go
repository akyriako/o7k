package blockstorage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/blockstorage/v3/volumes"
)

func (r *EvsVolumes) show(ctx context.Context, id string) (pluginsdk.Result, error) {
	client, err := r.plugin.BlockStorageV3(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting block storage client: %w", err)
	}

	volume, err := volumes.Get(client, id)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting volume %q: %w", id, err)
	}

	content, err := json.Marshal(volume)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding volume %q: %w", id, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      volume.ID,
			Content: content,
		},
	}, nil
}

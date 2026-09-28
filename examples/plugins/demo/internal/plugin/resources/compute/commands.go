package compute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
)

func (r *EcsServers) show(ctx context.Context, id string) (pluginsdk.Result, error) {
	client, err := r.plugin.ComputeV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting compute client: %w", err)
	}

	server, err := servers.Get(client, id).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting server %q: %w", id, err)
	}

	content, err := json.Marshal(server)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding server %q: %w", id, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      server.ID,
			Content: content,
		},
	}, nil
}

func (r *EcsServers) volumes(row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "evs-volumes",
			Field:    "server_id",
			Value:    row.ID,
		},
	}, nil
}

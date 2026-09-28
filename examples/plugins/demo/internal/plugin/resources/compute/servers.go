package compute

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
)

type EcsServers struct {
	plugin *plugin.Plugin
}

func NewServers(p *plugin.Plugin) *EcsServers {
	return &EcsServers{
		plugin: p,
	}
}

func (r *EcsServers) Service() string {
	return "opentelekomcloud-compute"
}

func (r *EcsServers) Kind() string {
	return "ecs-servers"
}

func (r *EcsServers) Title() string {
	return "ECS Servers"
}

func (r *EcsServers) Aliases() []string {
	return nil
}

func (r *EcsServers) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "Name", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "Status", MinWidth: 12},
	}
}

func (r *EcsServers) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *EcsServers) List(ctx context.Context) ([]pluginsdk.Row, error) {
	current, err := r.plugin.Host().Context(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting context: %w", err)
	}

	provider, err := r.plugin.Provider().Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting provider: %w", err)
	}

	client, err := openstack.NewComputeV2(provider, golangsdk.EndpointOpts{
		Region: current.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("creating compute client: %w", err)
	}

	pages, err := servers.List(client, servers.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing servers: %w", err)
	}

	allServers, err := servers.ExtractServers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting servers: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(allServers))

	for _, server := range allServers {
		rows = append(rows, pluginsdk.Row{
			ID: server.ID,
			Fields: map[string]string{
				"name":   server.Name,
				"status": server.Status,
				"id":     server.ID,
			},
		})
	}

	return rows, nil
}

func (r *EcsServers) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row.ID)
	}

	return pluginsdk.Result{}, nil
}

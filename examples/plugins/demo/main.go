package main

import (
	"context"
	"strconv"

	plugin "github.com/akyriako/o7k/plugin"
	"github.com/gophercloud/gophercloud/v2"
)

type demoPlugin struct {
	host      plugin.Host
	openstack *plugin.ClientProvider[*gophercloud.ProviderClient]
}

func (p *demoPlugin) Metadata() plugin.Metadata {
	return plugin.Metadata{
		Name:    "demo",
		Version: "0.0.1",
	}
}

func (p *demoPlugin) Resources() []plugin.Resource {
	return []plugin.Resource{
		&contextResource{plugin: p},
	}
}

func (p *demoPlugin) SetHost(host plugin.Host) {
	p.host = host
	p.openstack = newClient(host)
}

type contextResource struct {
	plugin *demoPlugin
}

func (r *contextResource) Service() string {
	return "demo"
}

func (r *contextResource) Kind() string {
	return "demo"
}

func (r *contextResource) Title() string {
	return "Demo Plugin"
}

func (r *contextResource) Aliases() []string {
	return nil
}

func (r *contextResource) Columns() []plugin.Column {
	return []plugin.Column{
		{Key: "cloud", Title: "Cloud", MinWidth: 20, Flex: 1},
		{Key: "region", Title: "Region", MinWidth: 20, Flex: 1},
		{Key: "generation", Title: "Generation", MinWidth: 10},
		{Key: "identity", Title: "Identity Endpoint", MinWidth: 30, Flex: 2},
	}
}

func (r *contextResource) Commands() []plugin.Command {
	return nil
}

func (r *contextResource) List(ctx context.Context) ([]plugin.Row, error) {
	current, err := r.plugin.host.Context(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := r.plugin.openstack.Client(ctx)
	if err != nil {
		return nil, err
	}

	return []plugin.Row{
		{
			ID: current.Cloud,
			Fields: map[string]string{
				"cloud":      current.Cloud,
				"region":     current.Region,
				"generation": strconv.FormatUint(current.Generation, 10),
				"identity":   provider.IdentityEndpoint,
			},
		},
	}, nil
}

func (r *contextResource) Execute(context.Context, plugin.Command, plugin.Row) (plugin.Result, error) {
	return plugin.Result{}, nil
}

func main() {
	plugin.Serve(&demoPlugin{})
}

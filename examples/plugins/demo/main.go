package main

import (
	"context"
	"strconv"

	plugin "github.com/akyriako/o7k/plugin"
)

type demoPlugin struct {
	host plugin.Host
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
}

type contextResource struct {
	plugin *demoPlugin
}

func (r *contextResource) Service() string {
	return "demo"
}

func (r *contextResource) Kind() string {
	return "plugin-context"
}

func (r *contextResource) Title() string {
	return "Plugin Context"
}

func (r *contextResource) Aliases() []string {
	return nil
}

func (r *contextResource) Columns() []plugin.Column {
	return []plugin.Column{
		{Key: "cloud", Title: "Cloud", MinWidth: 20, Flex: 1},
		{Key: "region", Title: "Region", MinWidth: 20, Flex: 1},
		{Key: "generation", Title: "Generation", MinWidth: 10},
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

	return []plugin.Row{
		{
			ID: current.Cloud,
			Fields: map[string]string{
				"cloud":      current.Cloud,
				"region":     current.Region,
				"generation": strconv.FormatUint(current.Generation, 10),
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

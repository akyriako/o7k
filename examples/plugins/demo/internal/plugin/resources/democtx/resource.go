package democtx

import (
	"context"
	"strconv"

	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin"
	"github.com/akyriako/o7k/pluginsdk"
)

type DemoCtx struct {
	plugin *plugin.Plugin
}

func NewDemoCtx(p *plugin.Plugin) *DemoCtx {
	return &DemoCtx{
		plugin: p,
	}
}

func (r *DemoCtx) Service() string {
	return "demo"
}

func (r *DemoCtx) Kind() string {
	return "demo"
}

func (r *DemoCtx) Title() string {
	return "Demo Plugin"
}

func (r *DemoCtx) Aliases() []string {
	return nil
}

func (r *DemoCtx) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "cloud", Title: "CLOUD", MinWidth: 20, Flex: 1},
		{Key: "region", Title: "REGION", MinWidth: 20, Flex: 1},
		{Key: "generation", Title: "GENERATION", MinWidth: 10},
	}
}

func (r *DemoCtx) Commands() []pluginsdk.Command {
	return nil
}

func (r *DemoCtx) List(ctx context.Context) ([]pluginsdk.Row, error) {
	current, err := r.plugin.Host().Context(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.plugin.Provider().Client(ctx)
	if err != nil {
		return nil, err
	}

	return []pluginsdk.Row{
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

func (r *DemoCtx) Execute(context.Context, pluginsdk.Command, pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{}, nil
}

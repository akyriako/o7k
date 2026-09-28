package main

import (
	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin"
	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin/resources/blockstorage"
	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin/resources/compute"
	"github.com/akyriako/o7k/examples/plugins/demo/internal/plugin/resources/democtx"
	"github.com/akyriako/o7k/pluginsdk"
)

func main() {
	p := plugin.New()

	p.Register(
		compute.NewServers(p),
		democtx.NewDemoCtx(p),
		blockstorage.NewVolumes(p),
	)

	pluginsdk.Serve(p)
}

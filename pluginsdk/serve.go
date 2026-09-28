package pluginsdk

import hplugin "github.com/hashicorp/go-plugin"

func Serve(impl Plugin) {
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]hplugin.Plugin{
			PluginName: &GRPCPlugin{Impl: impl},
		},
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}

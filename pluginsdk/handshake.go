package pluginsdk

import hplugin "github.com/hashicorp/go-plugin"

var HandshakeConfig = hplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "O7K_PLUGIN",
	MagicCookieValue: "o7k",
}

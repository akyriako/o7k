package plugins

import (
	"fmt"
	"os/exec"

	"github.com/akyriako/o7k/internal/resource"
	hplugin "github.com/hashicorp/go-plugin"

	plugin "github.com/akyriako/o7k/plugin"
)

type Client struct {
	client *hplugin.Client
	plugin plugin.Client
}

func NewClient(path string, host plugin.Host) (*Client, error) {
	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig,
		Plugins: map[string]hplugin.Plugin{
			plugin.PluginName: &plugin.GRPCPlugin{
				Host: host,
			},
		},
		Cmd: exec.Command(path),
		AllowedProtocols: []hplugin.Protocol{
			hplugin.ProtocolGRPC,
		},
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("starting plugin %q: %w", path, err)
	}

	instance, err := rpcClient.Dispense(plugin.PluginName)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("dispensing plugin %q: %w", path, err)
	}

	remote, ok := instance.(plugin.Client)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("plugin %q returned unexpected client type %T", path, instance)
	}

	return &Client{
		client: client,
		plugin: remote,
	}, nil
}

func (c *Client) Plugin() plugin.Client {
	return c.plugin
}

func (c *Client) Close() {
	c.client.Kill()
}

func (c *Client) Resources() []resource.Resource {
	resources := c.plugin.Resources()
	result := make([]resource.Resource, 0, len(resources))

	for _, r := range resources {
		result = append(result, newResource(r))
	}

	return result
}

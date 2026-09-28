package plugins

import (
	"fmt"
	"os/exec"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/pluginsdk"
	hplugin "github.com/hashicorp/go-plugin"
)

type Client struct {
	client *hplugin.Client
	plugin pluginsdk.Client
}

func NewClient(path string, host pluginsdk.Host) (*Client, error) {
	client := hplugin.NewClient(&hplugin.ClientConfig{
		HandshakeConfig: pluginsdk.HandshakeConfig,
		Plugins: map[string]hplugin.Plugin{
			pluginsdk.PluginName: &pluginsdk.GRPCPlugin{
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

	instance, err := rpcClient.Dispense(pluginsdk.PluginName)
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("dispensing plugin %q: %w", path, err)
	}

	remote, ok := instance.(pluginsdk.Client)
	if !ok {
		client.Kill()
		return nil, fmt.Errorf("plugin %q returned unexpected client type %T", path, instance)
	}

	return &Client{
		client: client,
		plugin: remote,
	}, nil
}

func (c *Client) Plugin() pluginsdk.Client {
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

func (c *Client) Register(registry *resource.Registry) error {
	if err := registry.RegisterAll(c.Resources()); err != nil {
		return fmt.Errorf("registering plugin resources: %w", err)
	}

	return nil
}

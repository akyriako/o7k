package main

import (
	"context"
	"fmt"

	plugin "github.com/akyriako/o7k/plugin"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
)

func newClient(host plugin.Host) *plugin.ClientProvider[*gophercloud.ProviderClient] {
	return plugin.NewClientProvider(host, connectClient)
}

func connectClient(ctx context.Context, current plugin.Context) (*gophercloud.ProviderClient, error) {
	authOpts, _, tlsConfig, err := clouds.Parse(
		clouds.WithLocations(current.CloudsPath),
		clouds.WithCloudName(current.Cloud),
	)
	if err != nil {
		return nil, fmt.Errorf("parsing cloud %q: %w", current.Cloud, err)
	}

	authOpts.AllowReauth = true

	provider, err := config.NewProviderClient(ctx, authOpts, config.WithTLSConfig(tlsConfig))
	if err != nil {
		return nil, fmt.Errorf("authenticating cloud %q: %w", current.Cloud, err)
	}

	return provider, nil
}

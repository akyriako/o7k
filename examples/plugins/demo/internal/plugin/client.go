package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"gopkg.in/yaml.v3"
)

func newClient(host pluginsdk.Host) *pluginsdk.ClientProvider[*golangsdk.ProviderClient] {
	return pluginsdk.NewClientProvider(host, connectClient)
}

func connectClient(_ context.Context, current pluginsdk.Context) (*golangsdk.ProviderClient, error) {
	data, err := os.ReadFile(current.CloudsPath)
	if err != nil {
		return nil, fmt.Errorf("reading clouds.yaml %q: %w", current.CloudsPath, err)
	}

	var config openstack.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing clouds.yaml %q: %w", current.CloudsPath, err)
	}

	cloud, ok := config.Clouds[current.Cloud]
	if !ok {
		return nil, fmt.Errorf("cloud %q not found in %q", current.Cloud, current.CloudsPath)
	}

	if current.Region != "" {
		cloud.RegionName = current.Region
	}

	provider, err := openstack.AuthenticatedClientFromCloud(&cloud)
	if err != nil {
		return nil, fmt.Errorf("authenticating cloud %q: %w", current.Cloud, err)
	}

	return provider, nil
}

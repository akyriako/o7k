package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	identitytokens "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

type Context struct {
	Cloud    string
	Region   string
	Project  string
	Domain   string
	Identity string

	Provider *gophercloud.ProviderClient
}

func (c *Context) Connect(ctx context.Context, cloudsPath string) error {
	authOpts, _, tlsConfig, err := clouds.Parse(
		clouds.WithLocations(cloudsPath),
		clouds.WithCloudName(c.Cloud),
	)
	if err != nil {
		return fmt.Errorf("parsing cloud %q: %w", c.Cloud, err)
	}

	provider, err := config.NewProviderClient(ctx, authOpts, config.WithTLSConfig(tlsConfig))
	if err != nil {
		return fmt.Errorf("authenticating cloud %q: %w", c.Cloud, err)
	}

	c.Provider = provider

	return nil
}

func (c *Context) ServiceCatalog() (*identitytokens.ServiceCatalog, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	result := c.Provider.GetAuthResult()
	if result == nil {
		return nil, fmt.Errorf("context %q has no authentication result", c.Cloud)
	}

	switch result := result.(type) {
	case identitytokens.CreateResult:
		catalog, err := result.ExtractServiceCatalog()
		if err != nil {
			return nil, fmt.Errorf("extracting service catalog: %w", err)
		}

		return catalog, nil

	case identitytokens.GetResult:
		catalog, err := result.ExtractServiceCatalog()
		if err != nil {
			return nil, fmt.Errorf("extracting service catalog: %w", err)
		}

		return catalog, nil

	default:
		return nil, fmt.Errorf(
			"unsupported authentication result %T",
			result,
		)
	}
}

func (c *Context) IdentityV3() (*gophercloud.ServiceClient, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	client, err := openstack.NewIdentityV3(c.Provider, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, fmt.Errorf("creating identity client: %w", err)
	}

	return client, nil
}

func (c *Context) ComputeV2() (*gophercloud.ServiceClient, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	client, err := openstack.NewComputeV2(c.Provider, gophercloud.EndpointOpts{
		Region: c.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("creating compute client: %w", err)
	}

	return client, nil
}

func (c *Context) NetworkV2() (*gophercloud.ServiceClient, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	client, err := openstack.NewNetworkV2(c.Provider, gophercloud.EndpointOpts{
		Region: c.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("creating network client: %w", err)
	}

	return client, nil
}

func (c *Context) ImageV2() (*gophercloud.ServiceClient, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	client, err := openstack.NewImageV2(c.Provider, gophercloud.EndpointOpts{
		Region: c.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("creating image client: %w", err)
	}

	return client, nil
}

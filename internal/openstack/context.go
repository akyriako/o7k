package openstack

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
	"github.com/gophercloud/gophercloud/v2/openstack/config/clouds"
	identitytokens "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

type Context struct {
	CloudsPath string

	Cloud    string
	Region   string
	Project  string
	Domain   string
	Identity string

	Provider       *gophercloud.ProviderClient
	serviceClients sync.Map
}

func (c *Context) Connect(ctx context.Context, cloudsPath string) error {
	authOpts, _, tlsConfig, err := clouds.Parse(
		clouds.WithLocations(cloudsPath),
		clouds.WithCloudName(c.Cloud),
	)
	if err != nil {
		return fmt.Errorf("parsing cloud %q: %w", c.Cloud, err)
	}

	authOpts.AllowReauth = true

	provider, err := config.NewProviderClient(ctx, authOpts, config.WithTLSConfig(tlsConfig))
	if err != nil {
		return fmt.Errorf("authenticating cloud %q: %w", c.Cloud, err)
	}

	c.Provider = provider
	c.CloudsPath = cloudsPath

	slog.Info("connected to cloud", "cloud", c.Cloud, "identityEndpoint", provider.IdentityEndpoint)

	return nil
}

func (c *Context) Activate(other *Context) {
	c.Cloud = other.Cloud
	c.CloudsPath = other.CloudsPath
	c.Region = other.Region
	c.Project = other.Project
	c.Domain = other.Domain
	c.Identity = other.Identity
	c.Provider = other.Provider
	c.serviceClients = sync.Map{}
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

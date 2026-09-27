package openstack

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

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

	slog.Info("connected to cloud", "cloud", c.Cloud, "identityEndpoint", provider.IdentityEndpoint)

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
	return c.getClientService("identity", gophercloud.EndpointOpts{}, openstack.NewIdentityV3)
}

func (c *Context) ComputeV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("compute", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewComputeV2)
}

func (c *Context) NetworkV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("network", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewNetworkV2)
}

func (c *Context) ImageV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("image", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewImageV2)
}

func (c *Context) BlockStorageV3() (*gophercloud.ServiceClient, error) {
	return c.getClientService("block-storage", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewBlockStorageV3)
}

func (c *Context) OrchestrationV1() (*gophercloud.ServiceClient, error) {
	return c.getClientService("orchestration", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewOrchestrationV1)
}

func (c *Context) DNSV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("dns", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewDNSV2)
}

func (c *Context) LoadBalancerV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("load-balancer", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewLoadBalancerV2)
}

func (c *Context) KeyManagerV1() (*gophercloud.ServiceClient, error) {
	return c.getClientService("key-manager", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewKeyManagerV1)
}

func (c *Context) ObjectStorageV1() (*gophercloud.ServiceClient, error) {
	return c.getClientService("object-storage", gophercloud.EndpointOpts{Region: c.Region}, openstack.NewObjectStorageV1)
}

type clientServiceBuilder func(*gophercloud.ProviderClient, gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)

func (c *Context) getClientService(key string, opts gophercloud.EndpointOpts, builder clientServiceBuilder) (*gophercloud.ServiceClient, error) {
	if c.Provider == nil {
		return nil, fmt.Errorf("context %q is not connected", c.Cloud)
	}

	if client, ok := c.serviceClients.Load(key); ok {
		return client.(*gophercloud.ServiceClient), nil
	}

	client, err := builder(c.Provider, opts)
	if err != nil {
		return nil, fmt.Errorf("creating %s client: %w", key, err)
	}

	actual, _ := c.serviceClients.LoadOrStore(key, client)
	return actual.(*gophercloud.ServiceClient), nil
}

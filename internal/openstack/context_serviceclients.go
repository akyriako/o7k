package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/apiversions"
)

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

func (c *Context) SharedFileSystemV2() (*gophercloud.ServiceClient, error) {
	return c.getClientService("shared-file-system", gophercloud.EndpointOpts{Region: c.Region}, newSharedFileSystemV2)
}

func newSharedFileSystemV2(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	client, err := openstack.NewSharedFileSystemV2(provider, opts)
	if err != nil {
		return nil, err
	}

	pages, err := apiversions.List(client).AllPages(context.Background())
	if err != nil {
		return nil, fmt.Errorf("listing shared file system API versions: %w", err)
	}

	versions, err := apiversions.ExtractAPIVersions(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting shared file system API versions: %w", err)
	}

	for _, version := range versions {
		if version.ID != "v2.0" {
			continue
		}

		client.Microversion = min(version.Version, "2.60")
		return client, nil
	}

	return nil, fmt.Errorf("shared file system API v2 not found")
}

func (c *Context) ContainerInfraV1() (*gophercloud.ServiceClient, error) {
	return c.getClientService(
		"container-infra",
		gophercloud.EndpointOpts{Region: c.Region},
		newContainerInfraV1,
	)
}

func newContainerInfraV1(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	client, err := openstack.NewContainerInfraV1(provider, opts)
	if err != nil {
		return nil, err
	}

	client.Microversion = "1.14"

	return client, nil
}

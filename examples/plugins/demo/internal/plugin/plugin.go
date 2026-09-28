package plugin

import (
	"context"

	pluginsdk "github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
)

type Plugin struct {
	host      pluginsdk.Host
	provider  *pluginsdk.ClientProvider[*golangsdk.ProviderClient]
	resources []pluginsdk.Resource
}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) SetHost(host pluginsdk.Host) {
	p.host = host
	p.provider = newClient(host)
}

func (p *Plugin) Register(resources ...pluginsdk.Resource) {
	p.resources = append(p.resources, resources...)
}

func (p *Plugin) Metadata() pluginsdk.Metadata {
	return pluginsdk.Metadata{
		Name:    "demo",
		Version: "0.1.0",
	}
}

func (p *Plugin) Resources() []pluginsdk.Resource {
	return p.resources
}

func (p *Plugin) Host() pluginsdk.Host {
	return p.host
}

func (p *Plugin) Provider() *pluginsdk.ClientProvider[*golangsdk.ProviderClient] {
	return p.provider
}

func (p *Plugin) ComputeV2(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "compute", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewComputeV2(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}

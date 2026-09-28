package plugins

import (
	"context"
	"sync/atomic"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/pluginsdk"
)

type Host struct {
	context    *openstack.Context
	generation atomic.Uint64
}

func NewHost(openstackContext *openstack.Context) *Host {
	host := &Host{
		context: openstackContext,
	}

	host.generation.Store(1)

	return host
}

func (h *Host) Context(context.Context) (pluginsdk.Context, error) {
	return pluginsdk.Context{
		Generation: h.generation.Load(),
		Cloud:      h.context.Cloud,
		CloudsPath: h.context.CloudsPath,
		Region:     h.context.Region,
	}, nil
}

func (h *Host) ContextChanged() {
	h.generation.Add(1)
}

package plugins

import (
	"context"
	"sync"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/pluginsdk"
)

type Host struct {
	mu      sync.RWMutex
	context pluginsdk.Context
}

func NewHost(openstackContext *openstack.Context) *Host {
	return &Host{
		context: pluginsdk.Context{
			Generation: 1,
			Cloud:      openstackContext.Cloud,
			CloudsPath: openstackContext.CloudsPath,
			Region:     openstackContext.Region,
		},
	}
}

func (h *Host) Context(context.Context) (pluginsdk.Context, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.context, nil
}

func (h *Host) ContextChanged(openstackContext *openstack.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.context.Generation++
	h.context.Cloud = openstackContext.Cloud
	h.context.CloudsPath = openstackContext.CloudsPath
	h.context.Region = openstackContext.Region
}

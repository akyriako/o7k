package resource

import (
	"fmt"
	"strings"
)

type Registry struct {
	resources map[string]Resource
}

func NewRegistry() *Registry {
	return &Registry{
		resources: make(map[string]Resource),
	}
}

func (r *Registry) Register(resource Resource) error {
	names := append([]string{resource.Kind()}, resource.Aliases()...)

	for _, name := range names {
		name = normalize(name)

		if name == "" {
			return fmt.Errorf("resource name cannot be empty")
		}

		if _, exists := r.resources[name]; exists {
			return fmt.Errorf("resource name %q already registered", name)
		}

		r.resources[name] = resource
	}

	return nil
}

func (r *Registry) Get(name string) (Resource, bool) {
	resource, ok := r.resources[normalize(name)]
	return resource, ok
}

func normalize(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

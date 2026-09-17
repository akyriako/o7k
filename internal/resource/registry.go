package resource

import (
	"fmt"
	"strings"
)

var reservedCommandKeys = map[string]struct{}{
	"q":     {},
	":":     {},
	"r":     {},
	"esc":   {},
	"enter": {},
}

type Registry struct {
	resources   map[string]Resource
	commandKeys map[string]string
}

func NewRegistry() *Registry {
	return &Registry{
		resources:   make(map[string]Resource),
		commandKeys: make(map[string]string),
	}
}

func (r *Registry) Register(resource Resource) error {
	if err := r.validateCommands(resource); err != nil {
		return err
	}

	names := append([]string{resource.Kind()}, resource.Aliases()...)

	for _, name := range names {
		name = normalize(name)

		if name == "" {
			return fmt.Errorf("resource name cannot be empty")
		}

		if _, exists := r.resources[name]; exists {
			return fmt.Errorf("resource name %q already registered", name)
		}
	}

	for _, name := range names {
		r.resources[normalize(name)] = resource
	}

	for _, command := range resource.Commands() {
		r.commandKeys[normalize(command.Key)] = resource.Kind()
	}

	return nil
}

func (r *Registry) Get(name string) (Resource, bool) {
	resource, ok := r.resources[normalize(name)]
	return resource, ok
}

func (r *Registry) validateCommands(resource Resource) error {
	defaultCount := 0

	for _, command := range resource.Commands() {
		key := normalize(command.Key)

		if key == "" {
			return fmt.Errorf("resource %q has command with empty key", resource.Kind())
		}

		if _, reserved := reservedCommandKeys[key]; reserved {
			return fmt.Errorf("resource %q command key %q is reserved", resource.Kind(), command.Key)
		}

		if owner, exists := r.commandKeys[key]; exists {
			return fmt.Errorf("resource %q command key %q already used by resource %q", resource.Kind(), command.Key, owner)
		}

		if command.Default {
			defaultCount++

			if defaultCount > 1 {
				return fmt.Errorf("resource %q has more than one default command", resource.Kind())
			}
		}
	}

	return nil
}

func normalize(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

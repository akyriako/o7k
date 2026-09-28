package plugins

import (
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
)

type Status string

const (
	StatusLoaded Status = "loaded"
	StatusFailed Status = "failed"
)

type Info struct {
	Path    string
	Name    string
	Version string
	Status  Status
	Err     error
}

type Manager struct {
	host     *Host
	registry *resource.Registry
	clients  []*Client
	plugins  []Info
}

func NewManager(host *Host, registry *resource.Registry) *Manager {
	return &Manager{
		host:     host,
		registry: registry,
	}
}

func (m *Manager) Load(path string) error {
	client, err := NewClient(path, m.host)
	if err != nil {
		m.plugins = append(m.plugins, Info{
			Path:   path,
			Status: StatusFailed,
			Err:    err,
		})
		return err
	}

	metadata := client.Plugin().Metadata()

	if err := client.Register(m.registry); err != nil {
		client.Close()

		m.plugins = append(m.plugins, Info{
			Path:    path,
			Name:    metadata.Name,
			Version: metadata.Version,
			Status:  StatusFailed,
			Err:     err,
		})

		return fmt.Errorf("registering plugin %q: %w", path, err)
	}

	m.clients = append(m.clients, client)
	m.plugins = append(m.plugins, Info{
		Path:    path,
		Name:    metadata.Name,
		Version: metadata.Version,
		Status:  StatusLoaded,
	})

	return nil
}

func (m *Manager) Plugins() []Info {
	return m.plugins
}

func (m *Manager) Close() {
	for _, client := range m.clients {
		client.Close()
	}

	m.clients = nil
}

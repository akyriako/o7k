package plugins

import (
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
)

type Manager struct {
	host     *Host
	registry *resource.Registry
	clients  []*Client
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
		return err
	}

	if err := client.Register(m.registry); err != nil {
		client.Close()
		return fmt.Errorf("registering plugin %q: %w", path, err)
	}

	m.clients = append(m.clients, client)

	return nil
}

func (m *Manager) Close() {
	for _, client := range m.clients {
		client.Close()
	}

	m.clients = nil
}

package plugins

import (
	"context"

	pluginmanager "github.com/akyriako/o7k/internal/plugins"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	manager *pluginmanager.Manager
}

func New(manager *pluginmanager.Manager) *Resource {
	return &Resource{
		manager: manager,
	}
}

func (r *Resource) Kind() string {
	return "plugins"
}

func (r *Resource) Title() string {
	return "Plugins"
}

func (r *Resource) Aliases() []string {
	return nil
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "Name", MinWidth: 20, Flex: 1},
		{Key: "version", Title: "Version", MinWidth: 12},
		{Key: "status", Title: "Status", MinWidth: 10},
		{Key: "path", Title: "Path", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(context.Context) ([]resource.Row, error) {
	plugins := r.manager.Plugins()
	rows := make([]resource.Row, 0, len(plugins))

	for _, plugin := range plugins {
		name := plugin.Name
		if name == "" {
			name = "-"
		}

		version := plugin.Version
		if version == "" {
			version = "-"
		}

		rows = append(rows, resource.Row{
			ID: plugin.Path,
			Fields: map[string]string{
				"name":    name,
				"version": version,
				"status":  string(plugin.Status),
				"path":    plugin.Path,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(resource.Command, resource.Row) tea.Cmd {
	return nil
}

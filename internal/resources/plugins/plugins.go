package plugins

import (
	"context"
	"fmt"

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
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "version", Title: "VERSION", MinWidth: 12},
		{Key: "status", Title: "STATUS", MinWidth: 10},
		{Key: "path", Title: "PATH", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "d", Description: "Details", Default: true},
	}
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

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return func() tea.Msg {
		for _, plugin := range r.manager.Plugins() {
			if plugin.Path != row.ID {
				continue
			}

			content := map[string]any{
				"name":    plugin.Name,
				"version": plugin.Version,
				"status":  plugin.Status,
				"path":    plugin.Path,
			}

			if plugin.Err != nil {
				content["error"] = plugin.Err.Error()
			}

			return resource.DetailsMsg{
				ID:      row.ID,
				Content: content,
			}
		}

		return resource.DetailsMsg{
			ID:  row.ID,
			Err: fmt.Errorf("plugin %q not found", row.ID),
		}
	}
}

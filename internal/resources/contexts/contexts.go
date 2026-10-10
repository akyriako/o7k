package contexts

import (
	"context"
	"log/slog"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	cloudsPaths []string
}

func New(cloudsPaths []string) *Resource {
	return &Resource{cloudsPaths: cloudsPaths}
}

func (r *Resource) Kind() string {
	return "contexts"
}

func (r *Resource) Title() string {
	return "Contexts"
}

func (r *Resource) Aliases() []string {
	return []string{"context", "ctx", "cloud", "clouds"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 2},
		{Key: "active", Title: "ACTIVE", MinWidth: 8, Flex: 0},
		{Key: "region", Title: "REGION", MinWidth: 12, Flex: 1},
		{Key: "domain", Title: "DOMAIN", MinWidth: 20, Flex: 1},
		{Key: "project", Title: "PROJECT", MinWidth: 20, Flex: 2},
		{Key: "cloudsyaml", Title: "SOURCE", MinWidth: 40, Flex: 2},
		//{Key: "identity", Title: "IDENTITY", MinWidth: 24, Flex: 3},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{{Key: "a", Description: "Activate", Default: true}}
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	clouds := &openstack.Clouds{}

	if len(r.cloudsPaths) > 0 {
		loaded, err := openstack.LoadClouds(r.cloudsPaths)
		if err != nil {
			return nil, err
		}

		clouds = loaded
	}

	if cloud, ok := openstack.EnvCloud(); ok {
		if _, exists := clouds.Get(cloud.Name); exists {
			slog.Warn("ignoring OS_* environment variables, clouds.yaml defines a cloud with the same name", "cloud", cloud.Name)
		} else {
			cloud.Path = openstack.EnvCloudSource
			clouds.Items = append(clouds.Items, cloud)
		}
	}

	rows := make([]resource.Row, 0, len(clouds.Items))

	for _, cloud := range clouds.Items {
		rows = append(rows, resource.Row{
			ID: cloud.Name,
			Fields: map[string]string{
				"name":       cloud.Name,
				"region":     cloud.Region,
				"domain":     cloud.Domain,
				"project":    cloud.Project,
				"cloudsyaml": cloud.Path,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "a":
		return r.activate(row)
	}

	return nil
}

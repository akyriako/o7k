package contexts

import (
	"context"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type ActivatedMsg struct {
	Context *openstack.Context
	Err     error
}

type Resource struct {
	cloudsPath string
}

func New(clouds *openstack.Clouds) *Resource {
	return &Resource{cloudsPath: clouds.Path}
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
		{Key: "identity", Title: "IDENTITY", MinWidth: 24, Flex: 3},
		{Key: "region", Title: "REGION", MinWidth: 12, Flex: 1},
		{Key: "domain", Title: "DOMAIN", MinWidth: 16, Flex: 1},
		{Key: "project", Title: "PROJECT", MinWidth: 20, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{
			Key:         "a",
			Description: "Activate",
			Default:     true,
		},
	}
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	clouds, err := openstack.LoadClouds(r.cloudsPath)
	if err != nil {
		return nil, err
	}

	rows := make([]resource.Row, 0, len(clouds.Items))

	for _, cloud := range clouds.Items {
		rows = append(rows, resource.Row{
			ID: cloud.Name,
			Fields: map[string]string{
				"name":     cloud.Name,
				"identity": cloud.Identity,
				"region":   cloud.Region,
				"domain":   cloud.Domain,
				"project":  cloud.Project,
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

func (r *Resource) activate(row resource.Row) tea.Cmd {
	cloudsPath := r.cloudsPath
	cloudName := row.ID

	return func() tea.Msg {
		clouds, err := openstack.LoadClouds(cloudsPath)
		if err != nil {
			return ActivatedMsg{Err: err}
		}

		cloud, ok := clouds.Get(cloudName)
		if !ok {
			return ActivatedMsg{Err: &CloudNotFoundError{Name: cloudName}}
		}

		next := openstack.Context{
			Cloud:    cloud.Name,
			Identity: cloud.Identity,
			Region:   cloud.Region,
			Domain:   cloud.Domain,
			Project:  cloud.Project,
		}

		err = next.Connect(context.Background(), cloudsPath)

		return ActivatedMsg{
			Context: &next,
			Err:     err,
		}
	}
}

type CloudNotFoundError struct {
	Name string
}

func (e *CloudNotFoundError) Error() string {
	return "cloud not found: " + e.Name
}

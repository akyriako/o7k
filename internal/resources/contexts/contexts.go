package contexts

import (
	"context"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
)

type Resource struct {
	clouds *openstack.Clouds
}

func New(clouds *openstack.Clouds) *Resource {
	return &Resource{clouds: clouds}
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
		{Key: "identity", Title: "IDENTITY", MinWidth: 24, Flex: 3},
		{Key: "region", Title: "REGION", MinWidth: 12, Flex: 1},
		{Key: "domain", Title: "DOMAIN", MinWidth: 16, Flex: 1},
		{Key: "project", Title: "PROJECT", MinWidth: 20, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	rows := make([]resource.Row, 0, len(r.clouds.Items))

	for _, cloud := range r.clouds.Items {
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

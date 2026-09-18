package catalog

import (
	"context"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "catalog"
}

func (r *Resource) Title() string {
	return "Catalog"
}

func (r *Resource) Aliases() []string {
	return []string{"cat"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "service", Title: "SERVICE", MinWidth: 20, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 16, Flex: 1},
		{Key: "interface", Title: "INTERFACE", MinWidth: 10, Flex: 0},
		{Key: "region", Title: "REGION", MinWidth: 16, Flex: 1},
		{Key: "url", Title: "URL", MinWidth: 40, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	catalog, err := r.context.ServiceCatalog()
	if err != nil {
		return nil, err
	}

	rows := make([]resource.Row, 0)

	for _, entry := range catalog.Entries {
		for _, endpoint := range entry.Endpoints {
			rows = append(rows, resource.Row{
				ID: endpoint.ID,
				Fields: map[string]string{
					"service":   entry.Name,
					"type":      entry.Type,
					"interface": endpoint.Interface,
					"region":    endpoint.Region,
					"url":       endpoint.URL,
				},
			})
		}
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

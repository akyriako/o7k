package datastores

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbdatastores "github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "db-datastores"
}

func (r *Resource) Title() string {
	return "Trove Datastores"
}

func (r *Resource) Aliases() []string {
	return []string{"dbdatastore", "trovedatastores"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 3},
		{Key: "default_version", Title: "DEFAULT VERSION", MinWidth: 36, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DBV1()
	if err != nil {
		return nil, err
	}

	pages, err := dbdatastores.List(client).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing datastores: %w", err)
	}

	items, err := dbdatastores.ExtractDatastores(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting datastores: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, datastore := range items {
		defaultVersion := datastore.DefaultVersion

		for _, version := range datastore.Versions {
			if version.ID == datastore.DefaultVersion {
				defaultVersion = version.Name
				break
			}
		}

		rows = append(rows, resource.Row{
			ID: datastore.ID,
			Fields: map[string]string{
				"id":              datastore.ID,
				"name":            datastore.Name,
				"default_version": defaultVersion,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	}

	return nil
}

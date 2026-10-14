package databases

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbdatabases "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "db-databases"
}

func (r *Resource) Title() string {
	return "Trove Databases"
}

func (r *Resource) Aliases() []string {
	return []string{"db-database"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 2},
		{Key: "charset", Title: "CHARSET", MinWidth: 16, Flex: 1},
		{Key: "collate", Title: "COLLATE", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	instanceID := resource.Scope(ctx)["instance_id"]
	if instanceID == "" {
		return nil, fmt.Errorf("database instance ID is required")
	}

	client, err := r.context.DBV1()
	if err != nil {
		return nil, err
	}

	pages, err := dbdatabases.List(client, instanceID).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing databases for instance %q: %w", instanceID, err)
	}

	items, err := dbdatabases.ExtractDBs(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting databases: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, database := range items {
		rows = append(rows, resource.Row{
			ID: database.Name,
			Fields: map[string]string{
				"name":    database.Name,
				"charset": database.CharSet,
				"collate": database.Collate,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

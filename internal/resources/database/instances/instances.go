package instances

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbinstances "github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Database Instances"
}

func (r *Resource) Kind() string {
	return "database-instances"
}

func (r *Resource) Aliases() []string {
	return []string{
		"database-instance",
		"db-instances",
		"db-instance",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "datastore", Title: "DATASTORE", MinWidth: 18, Flex: 1},
		{Key: "hostname", Title: "HOSTNAME", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DatabaseV1()
	if err != nil {
		return nil, err
	}

	pages, err := dbinstances.List(client).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing database instances: %w", err)
	}

	items, err := dbinstances.ExtractInstances(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting database instances: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, instance := range items {
		rows = append(rows, resource.Row{
			ID: instance.ID,
			Fields: map[string]string{
				"id":        instance.ID,
				"name":      instance.Name,
				"status":    instance.Status,
				"datastore": instance.Datastore.Type,
				"hostname":  instance.Hostname,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

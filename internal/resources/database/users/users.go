package users

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbusers "github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "db-users"
}

func (r *Resource) Title() string {
	return "Trove Users"
}

func (r *Resource) Aliases() []string {
	return []string{"db-user"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 2},
		{Key: "databases", Title: "DATABASES", MinWidth: 40, Flex: 3},
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

	pages, err := dbusers.List(client, instanceID).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing users for instance %q: %w", instanceID, err)
	}

	items, err := dbusers.ExtractUsers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting database users: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, user := range items {
		databases := make([]string, 0, len(user.Databases))

		for _, database := range user.Databases {
			databases = append(databases, database.Name)
		}

		rows = append(rows, resource.Row{
			ID: user.Name,
			Fields: map[string]string{
				"name":      user.Name,
				"databases": strings.Join(databases, ", "),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

package users

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityusers "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "users"
}

func (r *Resource) Title() string {
	return "Keystone Users"
}

func (r *Resource) Aliases() []string {
	return []string{"user"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 32, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 2},
		{Key: "domain_id", Title: "DOMAIN ID", MinWidth: 32, Flex: 0},
		{Key: "default_project_id", Title: "DEFAULT PROJECT ID", MinWidth: 32, Flex: 0},
		{Key: "enabled", Title: "ENABLED", MinWidth: 7, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		//{Key: "s", Description: "Show", Default: true},
		{Key: "shift-d", Description: "Domain"},
		{Key: "shift-p", Description: "Project"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identityusers.List(client, identityusers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}

	items, err := identityusers.ExtractUsers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting users: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, user := range items {
		rows = append(rows, resource.Row{
			ID: user.ID,
			Fields: map[string]string{
				"id":                 user.ID,
				"name":               user.Name,
				"domain_id":          user.DomainID,
				"default_project_id": user.DefaultProjectID,
				"enabled":            fmt.Sprintf("%t", user.Enabled),
				"description":        user.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	//case "s":
	//	return r.show(row.ID)
	case "shift-d":
		return r.domain(row)
	case "shift-p":
		return r.project(row)
	}

	return nil
}

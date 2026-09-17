package roles

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityroles "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "roles"
}

func (r *Resource) Title() string {
	return "Roles"
}

func (r *Resource) Aliases() []string {
	return []string{"role"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 32, Flex: 1},
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "domain_id", Title: "DOMAIN ID", MinWidth: 32, Flex: 1},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identityroles.List(client, identityroles.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing roles: %w", err)
	}

	items, err := identityroles.ExtractRoles(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting roles: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, role := range items {
		rows = append(rows, resource.Row{
			ID: role.ID,
			Fields: map[string]string{
				"id":          role.ID,
				"name":        role.Name,
				"domain_id":   role.DomainID,
				"description": role.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

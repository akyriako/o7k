package groups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identitygroups "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "groups"
}

func (r *Resource) Title() string {
	return "Keystone Groups"
}

func (r *Resource) Aliases() []string {
	return []string{"group"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 32, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 2},
		{Key: "domain_id", Title: "DOMAIN ID", MinWidth: 32, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		//{Key: "s", Description: "Show", Default: true},
		{Key: "shift-d", Description: "Domain"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identitygroups.List(
		client,
		identitygroups.ListOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing groups: %w", err)
	}

	items, err := identitygroups.ExtractGroups(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting groups: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, group := range items {
		rows = append(rows, resource.Row{
			ID: group.ID,
			Fields: map[string]string{
				"id":          group.ID,
				"name":        group.Name,
				"domain_id":   group.DomainID,
				"description": group.Description,
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
	}

	return nil
}

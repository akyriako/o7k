package securitygroups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "securitygroups"
}

func (r *Resource) Title() string {
	return "Security Groups"
}

func (r *Resource) Aliases() []string {
	return []string{"securitygroup", "secgroups", "secgroup", "sg"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		//{Key: "project_id", Title: "PROJECT ID", MinWidth: 40, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-r", Description: "Rules", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := groups.List(client, groups.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing security groups: %w", err)
	}

	items, err := groups.ExtractGroups(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting security groups: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, group := range items {
		rows = append(rows, resource.Row{
			ID: group.ID,
			Fields: map[string]string{
				"id":          group.ID,
				"name":        group.Name,
				"project_id":  group.ProjectID,
				"description": group.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-r":
		return r.navigateToRules(row)
	}

	return nil
}

package servergroups

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "servergroups"
}

func (r *Resource) Title() string {
	return "Server Groups"
}

func (r *Resource) Aliases() []string {
	return []string{"servergroup", "sgroups", "sgroup"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "policies", Title: "POLICIES", MinWidth: 20, Flex: 1},
		{Key: "members", Title: "MEMBERS", MinWidth: 40, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := servergroups.List(client, servergroups.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing server groups: %w", err)
	}

	items, err := servergroups.ExtractServerGroups(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting server groups: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, serverGroup := range items {
		rows = append(rows, resource.Row{
			ID: serverGroup.ID,
			Fields: map[string]string{
				"id":       serverGroup.ID,
				"name":     serverGroup.Name,
				"policies": strings.Join(serverGroup.Policies, ", "),
				"members":  strings.Join(serverGroup.Members, ", "),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

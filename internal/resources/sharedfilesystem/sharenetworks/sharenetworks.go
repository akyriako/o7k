package sharenetworks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasharenetworks "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharenetworks"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Networks"
}

func (r *Resource) Kind() string {
	return "share-networks"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-network",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 32, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	pages, err := manilasharenetworks.ListDetail(client, manilasharenetworks.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing share networks: %w", err)
	}

	items, err := manilasharenetworks.ExtractShareNetworks(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting share networks: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, network := range items {
		rows = append(rows, resource.Row{
			ID: network.ID,
			Fields: map[string]string{
				"id":          network.ID,
				"name":        network.Name,
				"description": network.Description,
				"project_id":  network.ProjectID,
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

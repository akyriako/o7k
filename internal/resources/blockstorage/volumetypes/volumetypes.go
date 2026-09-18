package volumetypes

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "volumetypes"
}

func (r *Resource) Title() string {
	return "Volume Types"
}

func (r *Resource) Aliases() []string {
	return []string{"volumetype", "types", "type"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
		{Key: "is_public", Title: "PUBLIC", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.BlockStorageV3()
	if err != nil {
		return nil, err
	}

	pages, err := volumetypes.List(client, volumetypes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing volume types: %w", err)
	}

	items, err := volumetypes.ExtractVolumeTypes(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting volume types: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, volumeType := range items {
		rows = append(rows, resource.Row{
			ID: volumeType.ID,
			Fields: map[string]string{
				"id":          volumeType.ID,
				"name":        volumeType.Name,
				"description": volumeType.Description,
				"is_public":   fmt.Sprintf("%t", volumeType.IsPublic),
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

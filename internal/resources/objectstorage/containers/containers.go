package containers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	swiftcontainers "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{context: openstackContext}
}

func (r *Resource) Kind() string {
	return "swift-containers"
}

func (r *Resource) Title() string {
	return "Swift Containers"
}

func (r *Resource) Aliases() []string {
	return []string{
		"swift-container",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "objects", Title: "OBJECTS", MinWidth: 10, Flex: 0},
		{Key: "bytes", Title: "BYTES", MinWidth: 14, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-o", Description: "Objects", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ObjectStorageV1()
	if err != nil {
		return nil, err
	}

	pages, err := swiftcontainers.List(client, nil).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing object containers: %w", err)
	}

	items, err := swiftcontainers.ExtractInfo(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting object containers: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.Name,
			Fields: map[string]string{
				"name":    item.Name,
				"objects": strconv.FormatInt(item.Count, 10),
				"bytes":   strconv.FormatInt(item.Bytes, 10),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-o":
		return r.objects(row)
	default:
		return nil
	}
}

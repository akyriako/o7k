package objects

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	swiftobjects "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{context: openstackContext}
}

func (r *Resource) Kind() string {
	return "swift-objects"
}

func (r *Resource) Title() string {
	return "Swift Objects"
}

func (r *Resource) Aliases() []string {
	return []string{
		"swift-object",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 32, Flex: 1},
		{Key: "bytes", Title: "BYTES", MinWidth: 14, Flex: 0},
		{Key: "content_type", Title: "CONTENT TYPE", MinWidth: 24, Flex: 0},
		{Key: "hash", Title: "HASH", MinWidth: 32, Flex: 0},
		{Key: "last_modified", Title: "LAST MODIFIED", MinWidth: 24, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "v", Description: "View Content"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	container := scope["container"]
	if container == "" {
		return nil, fmt.Errorf("objects require container scope")
	}

	client, err := r.context.ObjectStorageV1()
	if err != nil {
		return nil, err
	}

	pages, err := swiftobjects.List(
		client,
		container,
		nil,
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing objects: %w", err)
	}

	items, err := swiftobjects.ExtractInfo(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting objects: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.Name,
			Fields: map[string]string{
				"name":          item.Name,
				"bytes":         strconv.FormatInt(item.Bytes, 10),
				"content_type":  item.ContentType,
				"hash":          item.Hash,
				"last_modified": item.LastModified.Format(time.RFC3339),
				"container":     container,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "v":
		return r.view(row)
	default:
		return nil
	}
}

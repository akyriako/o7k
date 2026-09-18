package images

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "images"
}

func (r *Resource) Title() string {
	return "Images"
}

func (r *Resource) Aliases() []string {
	return []string{"image", "img"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "visibility", Title: "VISIBILITY", MinWidth: 12, Flex: 0},
		{Key: "disk_format", Title: "DISK FORMAT", MinWidth: 12, Flex: 0},
		{Key: "size", Title: "SIZE", MinWidth: 14, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ImageV2()
	if err != nil {
		return nil, err
	}

	pages, err := images.List(client, images.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing images: %w", err)
	}

	items, err := images.ExtractImages(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting images: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, image := range items {
		rows = append(rows, resource.Row{
			ID: image.ID,
			Fields: map[string]string{
				"id":          image.ID,
				"name":        image.Name,
				"status":      string(image.Status),
				"visibility":  string(image.Visibility),
				"disk_format": image.DiskFormat,
				"size":        strconv.FormatInt(image.SizeBytes, 10),
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

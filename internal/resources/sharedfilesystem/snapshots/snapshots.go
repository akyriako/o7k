package snapshots

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasnapshots "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/snapshots"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Snapshots"
}

func (r *Resource) Kind() string {
	return "share-snapshots"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-snapshot",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "size", Title: "SIZE", MinWidth: 8, Flex: 0},
		{Key: "share_id", Title: "SHARE ID", MinWidth: 40, Flex: 0},
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

	pages, err := manilasnapshots.ListDetail(client, manilasnapshots.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing share snapshots: %w", err)
	}

	items, err := manilasnapshots.ExtractSnapshots(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting share snapshots: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, snapshot := range items {
		rows = append(rows, resource.Row{
			ID: snapshot.ID,
			Fields: map[string]string{
				"id":       snapshot.ID,
				"name":     snapshot.Name,
				"status":   snapshot.Status,
				"size":     strconv.Itoa(snapshot.Size),
				"share_id": snapshot.ShareID,
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

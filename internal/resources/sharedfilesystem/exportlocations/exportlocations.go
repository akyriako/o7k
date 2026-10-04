package exportlocations

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilashares "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Export Locations"
}

func (r *Resource) Kind() string {
	return "share-export-locations"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-export-location",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "path", Title: "PATH", MinWidth: 40, Flex: 1},
		{Key: "preferred", Title: "PREFERRED", MinWidth: 10, Flex: 0},
		{Key: "is_admin_only", Title: "ADMIN ONLY", MinWidth: 10, Flex: 0},
		{Key: "share_instance_id", Title: "SHARE INSTANCE ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	shareID := scope["share_id"]

	if shareID == "" {
		return nil, fmt.Errorf("share ID is required")
	}

	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	result := manilashares.ListExportLocations(ctx, client, shareID)

	locations, err := result.Extract()
	if err != nil {
		return nil, fmt.Errorf("listing share export locations: %w", err)
	}

	rows := make([]resource.Row, 0, len(locations))

	for _, location := range locations {
		rows = append(rows, resource.Row{
			ID: location.ID,
			Fields: map[string]string{
				"id":                location.ID,
				"path":              location.Path,
				"preferred":         strconv.FormatBool(location.Preferred),
				"is_admin_only":     strconv.FormatBool(location.IsAdminOnly),
				"share_instance_id": location.ShareInstanceID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

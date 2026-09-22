package regions

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityregions "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/regions"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "regions"
}

func (r *Resource) Title() string {
	return "Regions"
}

func (r *Resource) Aliases() []string {
	return []string{"region"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 24, Flex: 1},
		{Key: "parent_region_id", Title: "PARENT REGION ID", MinWidth: 32, Flex: 1},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identityregions.List(client, identityregions.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing regions: %w", err)
	}

	items, err := identityregions.ExtractRegions(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting regions: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, region := range items {
		rows = append(rows, resource.Row{
			ID: region.ID,
			Fields: map[string]string{
				"id":               region.ID,
				"parent_region_id": region.ParentRegionID,
				"description":      region.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

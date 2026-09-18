package availabilityzones

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "availabilityzones"
}

func (r *Resource) Title() string {
	return "Availability Zones"
}

func (r *Resource) Aliases() []string {
	return []string{"availabilityzone", "azs", "az"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "available", Title: "AVAILABLE", MinWidth: 12, Flex: 0},
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

	pages, err := availabilityzones.List(client).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing availability zones: %w", err)
	}

	items, err := availabilityzones.ExtractAvailabilityZones(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting availability zones: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, zone := range items {
		rows = append(rows, resource.Row{
			ID: zone.ZoneName,
			Fields: map[string]string{
				"name":      zone.ZoneName,
				"available": strconv.FormatBool(zone.ZoneState.Available),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

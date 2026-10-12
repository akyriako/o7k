package aggregates

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "aggregates"
}

func (r *Resource) Title() string {
	return "Compute Host Aggregates"
}

func (r *Resource) Aliases() []string {
	return []string{"aggregate", "hostaggregates", "hostaggregate"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 8, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 2},
		{Key: "availability_zone", Title: "AVAILABILITY ZONE", MinWidth: 20, Flex: 1},
		{Key: "hosts", Title: "HOSTS", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-h", Description: "Hypervisors"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := aggregates.List(client).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing host aggregates: %w", err)
	}

	items, err := aggregates.ExtractAggregates(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting host aggregates: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, aggregate := range items {
		id := strconv.Itoa(aggregate.ID)

		rows = append(rows, resource.Row{
			ID: id,
			Fields: map[string]string{
				"id":                id,
				"name":              aggregate.Name,
				"availability_zone": aggregate.AvailabilityZone,
				"hosts":             strings.Join(aggregate.Hosts, ", "),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-h":
		return r.hypervisors(row)
	}

	return nil
}

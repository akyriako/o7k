package flavors

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeflavors "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "flavors"
}

func (r *Resource) Title() string {
	return "Flavors"
}

func (r *Resource) Aliases() []string {
	return []string{"flavor"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "vcpus", Title: "VCPUS", MinWidth: 8, Flex: 0},
		{Key: "ram", Title: "RAM", MinWidth: 10, Flex: 0},
		{Key: "disk", Title: "DISK", MinWidth: 10, Flex: 0},
		{Key: "ephemeral", Title: "EPHEMERAL", MinWidth: 12, Flex: 0},
		{Key: "swap", Title: "SWAP", MinWidth: 10, Flex: 0},
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

	pages, err := computeflavors.ListDetail(client, computeflavors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing flavors: %w", err)
	}

	items, err := computeflavors.ExtractFlavors(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting flavors: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, flavor := range items {
		rows = append(rows, resource.Row{
			ID: flavor.ID,
			Fields: map[string]string{
				"id":        flavor.ID,
				"name":      flavor.Name,
				"vcpus":     strconv.Itoa(flavor.VCPUs),
				"ram":       strconv.Itoa(flavor.RAM),
				"disk":      strconv.Itoa(flavor.Disk),
				"ephemeral": strconv.Itoa(flavor.Ephemeral),
				"swap":      strconv.Itoa(flavor.Swap),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

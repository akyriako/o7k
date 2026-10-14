package instances

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbinstances "github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "db-instances"
}

func (r *Resource) Title() string {
	return "Trove Instances"
}

func (r *Resource) Aliases() []string {
	return []string{"db-instance", "trove", "trove-instances"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 3},
		{Key: "datastore", Title: "DATASTORE", MinWidth: 18, Flex: 1},
		{Key: "version", Title: "VERSION", MinWidth: 16, Flex: 1},
		{Key: "addresses", Title: "ADDRESSES", MinWidth: 32, Flex: 1},
		{Key: "volume_used", Title: "USED (GB)", MinWidth: 9, Flex: 0},
		{Key: "volume_size", Title: "SIZE (GB)", MinWidth: 9, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-f", Description: "Flavor"},
		{Key: "shift-d", Description: "Databases"},
		{Key: "shift-u", Description: "Users"},
		{Key: "shift-b", Description: "Backups"},
		{Key: "ctrl+r", Description: "Restart", StatusLabel: "Restarting"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DBV1()
	if err != nil {
		return nil, err
	}

	pages, err := dbinstances.List(client).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing database instances: %w", err)
	}

	items, err := dbinstances.ExtractInstances(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting database instances: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, instance := range items {
		rows = append(rows, resource.Row{
			ID: instance.ID,
			Fields: map[string]string{
				"id":          instance.ID,
				"name":        instance.Name,
				"status":      instance.Status,
				"datastore":   instance.Datastore.Type,
				"version":     instance.Datastore.Version,
				"addresses":   instanceAddresses(instance.Addresses),
				"flavor_id":   instance.Flavor.ID,
				"volume_size": fmt.Sprintf("%d", instance.Volume.Size),
				"volume_used": fmt.Sprintf("%.2f", instance.Volume.Used),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-f":
		return r.flavor(row)
	case "ctrl+r":
		return r.restart(row)
	case "shift-d":
		return r.databases(row)
	case "shift-u":
		return r.users(row)
	case "shift-b":
		return r.backups(row)
	}

	return nil
}

func instanceAddresses(instanceAddresses []dbinstances.Address) string {
	addresses := make([]string, 0, len(instanceAddresses))
	for _, ia := range instanceAddresses {
		if address := strings.TrimSpace(ia.Address); address != "" {
			addresses = append(addresses, address)
		}
	}
	return strings.Join(addresses, ", ")
}

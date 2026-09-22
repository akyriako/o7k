package zones

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dnszones "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "DNS Zones"
}

func (r *Resource) Kind() string {
	return "zones"
}

func (r *Resource) Aliases() []string {
	return []string{
		"zone",
		"dns-zones",
		"dns-zone",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 32, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 30, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "type", Title: "TYPE", MinWidth: 12, Flex: 0},
		{Key: "ttl", Title: "TTL", MinWidth: 8, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-r", Description: "Recordsets"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DNSV2()
	if err != nil {
		return nil, err
	}

	pages, err := dnszones.List(client, dnszones.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing DNS zones: %w", err)
	}

	items, err := dnszones.ExtractZones(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting DNS zones: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, zone := range items {
		rows = append(rows, resource.Row{
			ID: zone.ID,
			Fields: map[string]string{
				"id":          zone.ID,
				"name":        zone.Name,
				"status":      zone.Status,
				"type":        zone.Type,
				"ttl":         strconv.Itoa(zone.TTL),
				"description": zone.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-r":
		return r.navigateToRecordsets(row)
	}

	return nil
}

package floatingips

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "floatingips"
}

func (r *Resource) Title() string {
	return "Floating IPs"
}

func (r *Resource) Aliases() []string {
	return []string{"floatingip", "fips", "fip"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "floating_ip", Title: "FLOATING IP", MinWidth: 14, Flex: 0},
		{Key: "fixed_ip", Title: "FIXED IP", MinWidth: 14, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "floating_network_id", Title: "FLOATING NETWORK ID", MinWidth: 40, Flex: 0},
		{Key: "port_id", Title: "PORT ID", MinWidth: 40, Flex: 0},
		{Key: "router_id", Title: "ROUTER ID", MinWidth: 40, Flex: 0},
		//{Key: "project_id", Title: "PROJECT ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-r", Description: "Router", Default: true},
		{Key: "shift-p", Description: "Port"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := floatingips.List(client, floatingips.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing floating IPs: %w", err)
	}

	items, err := floatingips.ExtractFloatingIPs(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting floating IPs: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, floatingIP := range items {
		rows = append(rows, resource.Row{
			ID: floatingIP.ID,
			Fields: map[string]string{
				"id":                  floatingIP.ID,
				"floating_ip":         floatingIP.FloatingIP,
				"fixed_ip":            floatingIP.FixedIP,
				"status":              floatingIP.Status,
				"floating_network_id": floatingIP.FloatingNetworkID,
				"port_id":             floatingIP.PortID,
				"router_id":           floatingIP.RouterID,
				"project_id":          floatingIP.ProjectID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-r":
		return r.navigateToRouter(row)
	case "shift-p":
		return r.navigateToPort(row)
	}

	return nil
}

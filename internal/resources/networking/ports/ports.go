package ports

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "ports"
}

func (r *Resource) Title() string {
	return "Ports"
}

func (r *Resource) Aliases() []string {
	return []string{"port"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "network_id", Title: "NETWORK ID", MinWidth: 40, Flex: 0},
		{Key: "mac_address", Title: "MAC ADDRESS", MinWidth: 20, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "device_owner", Title: "DEVICE OWNER", MinWidth: 28, Flex: 1},
		{Key: "device_id", Title: "DEVICE ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-n", Description: "Network", Default: true},
		{Key: "shift-f", Description: "Floating IPs"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := ports.List(client, ports.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing ports: %w", err)
	}

	items, err := ports.ExtractPorts(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting ports: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, port := range items {
		rows = append(rows, resource.Row{
			ID: port.ID,
			Fields: map[string]string{
				"id":           port.ID,
				"name":         port.Name,
				"network_id":   port.NetworkID,
				"mac_address":  port.MACAddress,
				"status":       port.Status,
				"device_owner": port.DeviceOwner,
				"device_id":    port.DeviceID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-n":
		return r.navigateToNetwork(row)
	case "shift-f":
		return r.navigateToFloatingIPs(row)
	}

	return nil
}

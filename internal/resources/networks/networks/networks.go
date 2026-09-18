package networks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "networks"
}

func (r *Resource) Title() string {
	return "Networks"
}

func (r *Resource) Aliases() []string {
	return []string{"network", "net"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 10, Flex: 0},
		{Key: "admin_state", Title: "ADMIN STATE", MinWidth: 12, Flex: 0},
		{Key: "shared", Title: "SHARED", MinWidth: 8, Flex: 0},
		{Key: "external", Title: "EXTERNAL", MinWidth: 8, Flex: 0},
		{Key: "mtu", Title: "MTU", MinWidth: 8, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-s", Description: "Subnets", Default: true},
		{Key: "shift-p", Description: "Ports"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := networks.List(client, networks.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing networks: %w", err)
	}

	type networkWithExtensions struct {
		networks.Network
		external.NetworkExternalExt
		mtu.NetworkMTUExt
	}

	var items []networkWithExtensions

	if err := networks.ExtractNetworksInto(pages, &items); err != nil {
		return nil, fmt.Errorf("extracting networks: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, network := range items {
		rows = append(rows, resource.Row{
			ID: network.ID,
			Fields: map[string]string{
				"id":          network.ID,
				"name":        network.Name,
				"status":      network.Status,
				"admin_state": strconv.FormatBool(network.AdminStateUp),
				"shared":      strconv.FormatBool(network.Shared),
				"external":    strconv.FormatBool(network.External),
				"mtu":         strconv.Itoa(network.MTU),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-s":
		return r.navigateToSubnets(row)
	case "shift-p":
		return r.navigateToPorts(row)
	}

	return nil
}

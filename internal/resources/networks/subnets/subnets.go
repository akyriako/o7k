package subnets

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "subnets"
}

func (r *Resource) Title() string {
	return "Subnets"
}

func (r *Resource) Aliases() []string {
	return []string{"subnet"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "network_id", Title: "NETWORK ID", MinWidth: 40, Flex: 0},
		{Key: "cidr", Title: "CIDR", MinWidth: 20, Flex: 1},
		{Key: "ip_version", Title: "IP VERSION", MinWidth: 12, Flex: 0},
		{Key: "gateway", Title: "GATEWAY", MinWidth: 20, Flex: 1},
		{Key: "dhcp", Title: "DHCP", MinWidth: 8, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := subnets.List(client, subnets.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing subnets: %w", err)
	}

	items, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting subnets: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, subnet := range items {
		rows = append(rows, resource.Row{
			ID: subnet.ID,
			Fields: map[string]string{
				"id":         subnet.ID,
				"name":       subnet.Name,
				"network_id": subnet.NetworkID,
				"cidr":       subnet.CIDR,
				"ip_version": strconv.Itoa(subnet.IPVersion),
				"gateway":    subnet.GatewayIP,
				"dhcp":       strconv.FormatBool(subnet.EnableDHCP),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

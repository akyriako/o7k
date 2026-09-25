package loadbalancers

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Load Balancers"
}

func (r *Resource) Kind() string {
	return "loadbalancers"
}

func (r *Resource) Aliases() []string {
	return []string{
		"loadbalancer",
		"lbs",
		"lb",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
		{Key: "vip_address", Title: "VIP ADDRESS", MinWidth: 20, Flex: 0},
		{Key: "vip_network_id", Title: "VIP NETWORK ID", MinWidth: 40, Flex: 0},
		{Key: "vip_subnet_id", Title: "VIP SUBNET ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	pages, err := loadbalancers.List(client, loadbalancers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing load balancers: %w", err)
	}

	items, err := loadbalancers.ExtractLoadBalancers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting load balancers: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
				"vip_address":         item.VipAddress,
				"vip_network_id":      item.VipNetworkID,
				"vip_subnet_id":       item.VipSubnetID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	}

	return nil
}

package pools

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Pools"
}

func (r *Resource) Kind() string {
	return "pools"
}

func (r *Resource) Aliases() []string {
	return []string{
		"pool",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "protocol", Title: "PROTOCOL", MinWidth: 12, Flex: 0},
		{Key: "lb_algorithm", Title: "ALGORITHM", MinWidth: 20, Flex: 0},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
		{Key: "healthmonitor_id", Title: "HEALTH MONITOR ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	scope := resource.Scope(ctx)

	opts := pools.ListOpts{
		LoadbalancerID: scope["loadbalancer_id"],
	}

	pages, err := pools.List(client, opts).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing pools: %w", err)
	}

	items, err := pools.ExtractPools(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting pools: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"protocol":            item.Protocol,
				"lb_algorithm":        item.LBMethod,
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
				"healthmonitor_id":    item.MonitorID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

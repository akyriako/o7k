package listeners

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Octavia Listeners"
}

func (r *Resource) Kind() string {
	return "listeners"
}

func (r *Resource) Aliases() []string {
	return []string{
		"listener",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "protocol", Title: "PROTO", MinWidth: 5, Flex: 0},
		{Key: "protocol_port", Title: "PORT", MinWidth: 5, Flex: 0},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 12, Flex: 0},
		//{Key: "default_pool_id", Title: "DEFAULT POOL ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show"},
		{Key: "shift-p", Description: "Pools", Default: true},
		{Key: "shift-o", Description: "L7 Policies"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	scope := resource.Scope(ctx)
	loadbalancerID := scope["loadbalancer_id"]

	if loadbalancerID == "" {
		return nil, fmt.Errorf("listener requires loadbalancer_id")
	}

	pages, err := listeners.List(client, listeners.ListOpts{LoadbalancerID: loadbalancerID}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing listeners: %w", err)
	}

	items, err := listeners.ExtractListeners(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting listeners: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"protocol":            item.Protocol,
				"protocol_port":       strconv.Itoa(item.ProtocolPort),
				"provisioning_status": item.ProvisioningStatus,
				"default_pool_id":     item.DefaultPoolID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "shift-o":
		return r.l7Policies(row)
	case "shift-p":
		return r.pools(row)
	}

	return nil
}

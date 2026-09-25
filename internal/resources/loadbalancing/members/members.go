package members

import (
	"context"
	"fmt"
	"strconv"

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
	return "Members"
}

func (r *Resource) Kind() string {
	return "members"
}

func (r *Resource) Aliases() []string {
	return []string{
		"member",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "address", Title: "ADDRESS", MinWidth: 20, Flex: 0},
		{Key: "protocol_port", Title: "PORT", MinWidth: 8, Flex: 0},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
		{Key: "subnet_id", Title: "SUBNET ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	poolID := scope["pool_id"]

	if poolID == "" {
		return nil, fmt.Errorf("member requires pool_id")
	}

	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	pages, err := pools.ListMembers(client, poolID, pools.ListMembersOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing pool members: %w", err)
	}

	items, err := pools.ExtractMembers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting pool members: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"address":             item.Address,
				"protocol_port":       strconv.Itoa(item.ProtocolPort),
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
				"subnet_id":           item.SubnetID,
				"pool_id":             poolID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

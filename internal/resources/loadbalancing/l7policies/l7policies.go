package l7policies

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "L7 Policies"
}

func (r *Resource) Kind() string {
	return "l7policies"
}

func (r *Resource) Aliases() []string {
	return []string{
		"l7policy",
		"l7-policies",
		"l7-policy",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "action", Title: "ACTION", MinWidth: 20, Flex: 0},
		{Key: "position", Title: "POSITION", MinWidth: 10, Flex: 0},
		{Key: "listener_id", Title: "LISTENER ID", MinWidth: 40, Flex: 0},
		{Key: "redirect_pool_id", Title: "REDIRECT POOL ID", MinWidth: 40, Flex: 0},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
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

	pages, err := l7policies.List(client, l7policies.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing L7 policies: %w", err)
	}

	items, err := l7policies.ExtractL7Policies(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting L7 policies: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"action":              item.Action,
				"position":            strconv.FormatInt(int64(item.Position), 10),
				"listener_id":         item.ListenerID,
				"redirect_pool_id":    item.RedirectPoolID,
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

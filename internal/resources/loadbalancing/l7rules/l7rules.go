package l7rules

import (
	"context"
	"fmt"

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
	return "L7 Rules"
}

func (r *Resource) Kind() string {
	return "l7rules"
}

func (r *Resource) Aliases() []string {
	return []string{
		"l7rule",
		"l7-rules",
		"l7-rule",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "type", Title: "TYPE", MinWidth: 16, Flex: 0},
		{Key: "compare_type", Title: "COMPARE", MinWidth: 16, Flex: 0},
		{Key: "key", Title: "KEY", MinWidth: 24, Flex: 1},
		{Key: "value", Title: "VALUE", MinWidth: 30, Flex: 1},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	policyID := scope["l7policy_id"]

	if policyID == "" {
		return nil, fmt.Errorf("l7 rule requires l7policy_id")
	}

	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	pages, err := l7policies.ListRules(
		client,
		policyID,
		l7policies.ListRulesOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing L7 rules: %w", err)
	}

	items, err := l7policies.ExtractRules(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting L7 rules: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"type":                item.RuleType,
				"compare_type":        item.CompareType,
				"key":                 item.Key,
				"value":               item.Value,
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
				"l7policy_id":         policyID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

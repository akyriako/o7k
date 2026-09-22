package securitygrouprules

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "securitygrouprules"
}

func (r *Resource) Title() string {
	return "Security Group Rules"
}

func (r *Resource) Aliases() []string {
	return []string{"securitygrouprule", "secrules", "secrule", "sgr"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "direction", Title: "DIRECTION", MinWidth: 12, Flex: 0},
		{Key: "ethertype", Title: "ETHERTYPE", MinWidth: 12, Flex: 0},
		{Key: "protocol", Title: "PROTOCOL", MinWidth: 12, Flex: 0},
		{Key: "port_range", Title: "PORT RANGE", MinWidth: 16, Flex: 0},
		{Key: "remote_ip_prefix", Title: "REMOTE IP PREFIX", MinWidth: 20, Flex: 1},
		{Key: "remote_group_id", Title: "REMOTE GROUP ID", MinWidth: 40, Flex: 0},
		{Key: "security_group_id", Title: "SECURITY GROUP ID", MinWidth: 40, Flex: 0},
		//{Key: "project_id", Title: "PROJECT ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-g", Description: "Security Group"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := rules.List(client, rules.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing security group rules: %w", err)
	}

	items, err := rules.ExtractRules(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting security group rules: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, rule := range items {
		rows = append(rows, resource.Row{
			ID: rule.ID,
			Fields: map[string]string{
				"id":                rule.ID,
				"direction":         rule.Direction,
				"ethertype":         rule.EtherType,
				"protocol":          rule.Protocol,
				"port_range":        formatPortRange(rule.PortRangeMin, rule.PortRangeMax),
				"remote_ip_prefix":  rule.RemoteIPPrefix,
				"remote_group_id":   rule.RemoteGroupID,
				"security_group_id": rule.SecGroupID,
				"project_id":        rule.ProjectID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-g":
		return r.navigateToSecurityGroup(row)
	}

	return nil
}

func formatPortRange(min, max int) string {
	switch {
	case min == 0 && max == 0:
		return ""
	case min == max:
		return fmt.Sprintf("%d", min)
	default:
		return fmt.Sprintf("%d-%d", min, max)
	}
}

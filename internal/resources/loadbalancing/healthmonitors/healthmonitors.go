package healthmonitors

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Health Monitors"
}

func (r *Resource) Kind() string {
	return "healthmonitors"
}

func (r *Resource) Aliases() []string {
	return []string{
		"healthmonitor",
		"monitors",
		"monitor",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 12, Flex: 0},
		{Key: "delay", Title: "DELAY", MinWidth: 8, Flex: 0},
		{Key: "timeout", Title: "TIMEOUT", MinWidth: 8, Flex: 0},
		{Key: "max_retries", Title: "MAX RETRIES", MinWidth: 12, Flex: 0},
		{Key: "provisioning_status", Title: "PROVISIONING", MinWidth: 20, Flex: 0},
		{Key: "operating_status", Title: "OPERATING", MinWidth: 16, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-p", Description: "Pool"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.LoadBalancerV2()
	if err != nil {
		return nil, err
	}

	pages, err := monitors.List(client, monitors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing health monitors: %w", err)
	}

	items, err := monitors.ExtractMonitors(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting health monitors: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":                  item.ID,
				"name":                item.Name,
				"type":                item.Type,
				"delay":               strconv.Itoa(item.Delay),
				"timeout":             strconv.Itoa(item.Timeout),
				"max_retries":         strconv.Itoa(item.MaxRetries),
				"provisioning_status": item.ProvisioningStatus,
				"operating_status":    item.OperatingStatus,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-p":
		return r.pool(row)
	default:
		return nil
	}
}

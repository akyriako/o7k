package routers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "routers"
}

func (r *Resource) Title() string {
	return "Routers"
}

func (r *Resource) Aliases() []string {
	return []string{"router"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "admin_state", Title: "ADMIN STATE", MinWidth: 12, Flex: 0},
		{Key: "distributed", Title: "DISTRIBUTED", MinWidth: 12, Flex: 0},
		{Key: "project_id", Title: "PROJECT ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-p", Description: "Ports"},
		{Key: "shift-f", Description: "Floating IPs"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.NetworkV2()
	if err != nil {
		return nil, err
	}

	pages, err := routers.List(client, routers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing routers: %w", err)
	}

	items, err := routers.ExtractRouters(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting routers: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, router := range items {
		rows = append(rows, resource.Row{
			ID: router.ID,
			Fields: map[string]string{
				"id":          router.ID,
				"name":        router.Name,
				"status":      router.Status,
				"admin_state": strconv.FormatBool(router.AdminStateUp),
				"distributed": strconv.FormatBool(router.Distributed),
				"project_id":  router.ProjectID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-p":
		return r.navigateToPorts(row)
	case "shift-f":
		return r.navigateToFloatingIPs(row)
	}

	return nil
}

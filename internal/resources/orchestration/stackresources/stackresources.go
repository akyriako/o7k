package stackresources

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackresources"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Stack Resources"
}

func (r *Resource) Kind() string {
	return "stack-resources"
}

func (r *Resource) Aliases() []string {
	return []string{
		"stack-resource",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 30, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 20, Flex: 0},
		{Key: "physical_id", Title: "PHYSICAL ID", MinWidth: 40, Flex: 1},
		{Key: "stack_id", Title: "STACK ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)

	stackName := scope["stack_name"]
	stackID := scope["stack_id"]

	if stackName == "" || stackID == "" {
		return nil, fmt.Errorf("stack resource requires stack_name and stack_id")
	}

	client, err := r.context.OrchestrationV1()
	if err != nil {
		return nil, err
	}

	pages, err := stackresources.List(client, stackName, stackID, nil).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing stack resources: %w", err)
	}

	items, err := stackresources.ExtractResources(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting stack resources: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.Name,
			Fields: map[string]string{
				"name":        item.Name,
				"type":        item.Type,
				"status":      item.Status,
				"physical_id": item.PhysicalID,
				"stack_id":    stackID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

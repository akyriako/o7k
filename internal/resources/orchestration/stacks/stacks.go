package stacks

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	orchestrationstacks "github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stacks"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Stacks"
}

func (r *Resource) Kind() string {
	return "stacks"
}

func (r *Resource) Aliases() []string {
	return []string{
		"stack",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 20, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.OrchestrationV1()
	if err != nil {
		return nil, err
	}

	pages, err := orchestrationstacks.List(client, orchestrationstacks.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing stacks: %w", err)
	}

	items, err := orchestrationstacks.ExtractStacks(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting stacks: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, stack := range items {
		rows = append(rows, resource.Row{
			ID: stack.ID,
			Fields: map[string]string{
				"id":          stack.ID,
				"name":        stack.Name,
				"status":      stack.Status,
				"description": stack.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

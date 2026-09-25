package stackevents

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/stackevents"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Stack Events"
}

func (r *Resource) Kind() string {
	return "stack-events"
}

func (r *Resource) Aliases() []string {
	return []string{
		"stack-event",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "resource", Title: "RESOURCE", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 20, Flex: 0},
		{Key: "reason", Title: "REASON", MinWidth: 50, Flex: 1},
		{Key: "physical_id", Title: "PHYSICAL ID", MinWidth: 40, Flex: 0},
		{Key: "time", Title: "TIME", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)

	stackName := scope["stack_name"]
	stackID := scope["stack_id"]

	if stackName == "" || stackID == "" {
		return nil, fmt.Errorf("stack event requires stack_name and stack_id")
	}

	client, err := r.context.OrchestrationV1()
	if err != nil {
		return nil, err
	}

	pages, err := stackevents.List(client, stackName, stackID, nil).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing stack events: %w", err)
	}

	items, err := stackevents.ExtractEvents(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting stack events: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		rows = append(rows, resource.Row{
			ID: item.ID,
			Fields: map[string]string{
				"id":          item.ID,
				"resource":    item.ResourceName,
				"status":      item.ResourceStatus,
				"reason":      item.ResourceStatusReason,
				"physical_id": item.PhysicalResourceID,
				"time":        item.Time.String(),
				"stack_id":    stackID,
				"stack_name":  stackName,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	}

	return nil
}

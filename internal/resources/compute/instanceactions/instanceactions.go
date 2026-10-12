package instanceactions

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeactions "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/instanceactions"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Compute Instance Actions"
}

func (r *Resource) Kind() string {
	return "instanceactions"
}

func (r *Resource) Aliases() []string {
	return []string{"instanceaction"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "request_id", Title: "REQUEST ID", MinWidth: 40, Flex: 0},
		{Key: "action", Title: "ACTION", MinWidth: 20, Flex: 2},
		{Key: "start_time", Title: "START TIME", MinWidth: 26, Flex: 1},
		{Key: "user_id", Title: "USER ID", MinWidth: 36, Flex: 1},
		{Key: "message", Title: "MESSAGE", MinWidth: 24, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	serverID := resource.Scope(ctx)["server_id"]
	if serverID == "" {
		return nil, fmt.Errorf("server ID is required")
	}

	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := computeactions.List(
		client,
		serverID,
		computeactions.ListOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing instance actions: %w", err)
	}

	actions, err := computeactions.ExtractInstanceActions(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting instance actions: %w", err)
	}

	rows := make([]resource.Row, 0, len(actions))

	for _, action := range actions {
		rows = append(rows, resource.Row{
			ID: action.RequestID,
			Fields: map[string]string{
				"request_id": action.RequestID,
				"action":     action.Action,
				"start_time": action.StartTime.Format("2006-01-02 15:04:05"),
				"user_id":    action.UserID,
				"message":    action.Message,
				"server_id":  serverID,
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

package softwaredeployments

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Heat Software Deployments"
}

func (r *Resource) Kind() string {
	return "software-deployments"
}

func (r *Resource) Aliases() []string {
	return []string{"software-deployment"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 36, Flex: 0},
		{Key: "config_id", Title: "CONFIG ID", MinWidth: 36, Flex: 0},
		{Key: "server_id", Title: "SERVER ID", MinWidth: 36, Flex: 0},
		{Key: "action", Title: "ACTION", MinWidth: 6, Flex: 2},
		{Key: "status", Title: "STATUS", MinWidth: 6, Flex: 2},
		{Key: "status_reason", Title: "REASON", MinWidth: 21, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "o", Description: "Outputs"},
	}
}

type Deployment struct {
	ID           string `json:"id"`
	ConfigID     string `json:"config_id"`
	ServerID     string `json:"server_id"`
	Action       string `json:"action"`
	Status       string `json:"status"`
	StatusReason string `json:"status_reason"`
	CreationTime string `json:"creation_time"`
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.OrchestrationV1()
	if err != nil {
		return nil, err
	}

	var result struct {
		Deployments []Deployment `json:"software_deployments"`
	}

	_, err = client.Get(ctx, client.ServiceURL("software_deployments"), &result, nil)
	if err != nil {
		return nil, fmt.Errorf("listing software deployments: %w", err)
	}

	rows := make([]resource.Row, 0, len(result.Deployments))
	for _, deployment := range result.Deployments {
		rows = append(rows, resource.Row{
			ID: deployment.ID,
			Fields: map[string]string{
				"id":            deployment.ID,
				"config_id":     deployment.ConfigID,
				"server_id":     deployment.ServerID,
				"action":        deployment.Action,
				"status":        deployment.Status,
				"status_reason": deployment.StatusReason,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "o":
		return r.outputs(row)
	}

	return nil
}

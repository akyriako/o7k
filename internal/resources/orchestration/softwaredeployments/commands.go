package softwaredeployments

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dc9-dev/promptveil"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		var result struct {
			Deployment Deployment `json:"software_deployment"`
		}

		_, err = client.Get(context.Background(), client.ServiceURL("software_deployments", row.ID), &result, nil)
		if err != nil {
			return resource.DetailsMsg{Err: fmt.Errorf("getting software deployment %q: %w", row.ID, err)}
		}

		return resource.DetailsMsg{
			ID:      result.Deployment.ID,
			Content: result.Deployment,
		}
	}
}

func (r *Resource) outputs(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		var result struct {
			SoftwareDeployment struct {
				OutputValues map[string]any `json:"output_values"`
			} `json:"software_deployment"`
		}

		_, err = client.Get(context.Background(), client.ServiceURL("software_deployments", row.ID), &result, nil)
		if err != nil {
			return resource.DetailsMsg{Err: fmt.Errorf("getting software deployment outputs %q: %w", row.ID, err)}
		}

		outputs := result.SoftwareDeployment.OutputValues

		content := fmt.Sprintf(
			"deploy_stdout:\n%s\n\ndeploy_stderr:\n%s\n\ndeploy_status_code: %v",
			outputs["deploy_stdout"],
			outputs["deploy_stderr"],
			outputs["deploy_status_code"],
		)

		redacted := promptveil.Redact(content)
		return resource.DetailsMsg{
			ID:          row.ID,
			Content:     redacted.Text,
			WrapContent: new(true),
		}
	}
}

func (r *Resource) config(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.OrchestrationV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		var deployment struct {
			SoftwareDeployment struct {
				ConfigID string `json:"config_id"`
			} `json:"software_deployment"`
		}

		_, err = client.Get(context.Background(), client.ServiceURL("software_deployments", row.ID), &deployment, nil)
		if err != nil {
			return resource.DetailsMsg{Err: fmt.Errorf("getting software deployment %q: %w", row.ID, err)}
		}

		configID := deployment.SoftwareDeployment.ConfigID
		if configID == "" {
			return resource.DetailsMsg{Err: fmt.Errorf("software deployment %q has no config_id", row.ID)}
		}

		var result struct {
			SoftwareConfig map[string]any `json:"software_config"`
		}

		_, err = client.Get(context.Background(), client.ServiceURL("software_configs", configID), &result, nil)
		if err != nil {
			return resource.DetailsMsg{Err: fmt.Errorf("getting software config %q: %w", configID, err)}
		}

		content, err := json.MarshalIndent(result.SoftwareConfig, "", "  ")
		if err != nil {
			return resource.DetailsMsg{Err: fmt.Errorf("marshalling software config %q: %w", configID, err)}
		}

		redacted := promptveil.Redact(string(content))
		return resource.DetailsMsg{
			ID:          configID,
			Content:     redacted.Text,
			WrapContent: new(true),
		}
	}
}

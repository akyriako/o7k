package plugins

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	plugin "github.com/akyriako/o7k/plugin"
	tea "github.com/charmbracelet/bubbletea"
)

type pluginResource struct {
	resource plugin.Resource
}

func newResource(r plugin.Resource) resource.Resource {
	return &pluginResource{resource: r}
}

func (r *pluginResource) Kind() string {
	return r.resource.Kind()
}

func (r *pluginResource) Title() string {
	return r.resource.Title()
}

func (r *pluginResource) Aliases() []string {
	return r.resource.Aliases()
}

func (r *pluginResource) Columns() []resource.Column {
	columns := r.resource.Columns()
	result := make([]resource.Column, 0, len(columns))

	for _, column := range columns {
		result = append(result, resource.Column{
			Key:      column.Key,
			Title:    column.Title,
			MinWidth: column.MinWidth,
			Flex:     column.Flex,
		})
	}

	return result
}

func (r *pluginResource) Commands() []resource.Command {
	commands := r.resource.Commands()
	result := make([]resource.Command, 0, len(commands))

	for _, command := range commands {
		result = append(result, resource.Command{
			Key:         command.Key,
			Description: command.Description,
			Default:     command.Default,
		})
	}

	return result
}

func (r *pluginResource) List(ctx context.Context) ([]resource.Row, error) {
	rows, err := r.resource.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]resource.Row, 0, len(rows))

	for _, row := range rows {
		result = append(result, resource.Row{
			ID:     row.ID,
			Fields: row.Fields,
		})
	}

	return result, nil
}

func (r *pluginResource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return func() tea.Msg {
		result, err := r.resource.Execute(context.Background(), plugin.Command{
			Key:         command.Key,
			Description: command.Description,
			Default:     command.Default,
		}, plugin.Row{
			ID:     row.ID,
			Fields: row.Fields,
		})
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		if result.Details != nil {
			var content any

			if err := json.Unmarshal(result.Details.Content, &content); err != nil {
				return resource.DetailsMsg{
					Err: fmt.Errorf("decoding plugin details: %w", err),
				}
			}

			return resource.DetailsMsg{
				ID:      result.Details.ID,
				Content: content,
			}
		}

		if result.Navigate != nil {
			if len(result.Navigate.Scope) > 0 {
				return resource.NavigateScopedMsg{
					Resource: result.Navigate.Resource,
					Scope:    result.Navigate.Scope,
				}
			}

			if result.Navigate.Field != "" {
				return resource.NavigateFilteredMsg{
					Resource: result.Navigate.Resource,
					Field:    result.Navigate.Field,
					Value:    result.Navigate.Value,
				}
			}

			return resource.NavigateMsg{
				Resource: result.Navigate.Resource,
				ID:       result.Navigate.ID,
			}
		}

		return nil
	}
}

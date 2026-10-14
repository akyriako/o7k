package backups

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.DBV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		var result struct {
			Backup Backup `json:"backup"`
		}

		_, err = client.Get(
			context.Background(),
			client.ServiceURL("backups", id),
			&result,
			nil,
		)
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting database backup %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      id,
			Content: result.Backup,
		}
	}
}

func (r *Resource) instance(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "db-instances",
			Field:    "id",
			Value:    row.Fields["instance_id"],
		}
	}
}

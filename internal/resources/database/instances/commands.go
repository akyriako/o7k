package instances

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dbinstances "github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.DBV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		instance, err := dbinstances.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting database instance %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      instance.ID,
			Content: instance,
		}
	}
}

func (r *Resource) flavor(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "flavors",
			Field:    "id",
			Value:    row.Fields["flavor_id"],
		}
	}
}

func (r *Resource) restart(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		if !strings.EqualFold(row.Fields["status"], "ACTIVE") {
			return resource.ErrorMsg{
				Err: fmt.Errorf(
					"cannot restart database instance in %s state",
					row.Fields["status"],
				),
			}
		}

		client, err := r.context.DBV1()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		err = dbinstances.Restart(context.Background(), client, row.ID).ExtractErr()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("restarting database instance %q: %w", row.ID, err),
			}
		}

		return resource.CommandCompletedMsg{
			RefreshOnCompletion: true,
		}
	}
}

func (r *Resource) databases(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "db-databases",
			Scope: map[string]string{
				"instance_id": row.ID,
			},
		}
	}
}

func (r *Resource) users(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "db-users",
			Scope: map[string]string{
				"instance_id": row.ID,
			},
		}
	}
}

func (r *Resource) backups(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "db-backups",
			Field:    "instance_id",
			Value:    row.ID,
		}
	}
}

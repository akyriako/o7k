package aggregates

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		aggregateID, err := strconv.Atoi(id)
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("invalid host aggregate ID %q: %w", id, err),
			}
		}

		aggregate, err := aggregates.Get(
			context.Background(),
			client,
			aggregateID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting host aggregate %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      strconv.Itoa(aggregate.ID),
			Content: aggregate,
		}
	}
}

func (r *Resource) hypervisors(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		aggregateID, err := strconv.Atoi(row.ID)
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("invalid host aggregate ID %q: %w", row.ID, err),
			}
		}

		aggregate, err := aggregates.Get(context.Background(), client, aggregateID).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting host aggregate %q: %w", row.ID, err),
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "hypervisors",
			Field:    "host",
			Values:   aggregate.Hosts,
		}
	}
}

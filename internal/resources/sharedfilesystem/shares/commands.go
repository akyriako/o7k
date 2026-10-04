package shares

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilashares "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.SharedFileSystemV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		share, err := manilashares.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting share %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      share.ID,
			Content: share,
		}
	}
}

func (r *Resource) snapshots(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "share-snapshots",
			Field:    "share_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) replicas(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "share-replicas",
			Field:    "share_id",
			Value:    row.ID,
		}
	}
}

func (r *Resource) shareNetwork(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "share-networks",
			Field:    "id",
			Value:    row.Fields["share_network_id"],
		}
	}
}

func (r *Resource) accessRules(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "share-access-rules",
			Scope: map[string]string{
				"share_id": row.ID,
			},
		}
	}
}

func (r *Resource) shareType(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "share-types",
			Field:    "id",
			Value:    row.Fields["share_type_id"],
		}
	}
}

func (r *Resource) exportLocations(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "share-export-locations",
			Scope: map[string]string{
				"share_id": row.ID,
			},
		}
	}
}

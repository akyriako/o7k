package endpoints

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityendpoints "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "endpoints"
}

func (r *Resource) Title() string {
	return "Endpoints"
}

func (r *Resource) Aliases() []string {
	return []string{"endpoint", "ep"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 36, Flex: 1},
		{Key: "service_id", Title: "SERVICE ID", MinWidth: 36, Flex: 1},
		{Key: "interface", Title: "INTERFACE", MinWidth: 10, Flex: 0},
		{Key: "region", Title: "REGION", MinWidth: 12, Flex: 1},
		{Key: "enabled", Title: "ENABLED", MinWidth: 10, Flex: 0},
		{Key: "url", Title: "URL", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identityendpoints.List(client, identityendpoints.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing endpoints: %w", err)
	}

	items, err := identityendpoints.ExtractEndpoints(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting endpoints: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, endpoint := range items {
		rows = append(rows, resource.Row{
			ID: endpoint.ID,
			Fields: map[string]string{
				"id":         endpoint.ID,
				"service_id": endpoint.ServiceID,
				"interface":  string(endpoint.Availability),
				"region":     endpoint.Region,
				"enabled":    strconv.FormatBool(endpoint.Enabled),
				"url":        endpoint.URL,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

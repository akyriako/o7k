package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityservices "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "services"
}

func (r *Resource) Title() string {
	return "Services"
}

func (r *Resource) Aliases() []string {
	return []string{"service", "svc"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 32, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 20, Flex: 1},
		{Key: "enabled", Title: "ENABLED", MinWidth: 10, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-e", Description: "Endpoints"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.IdentityV3()
	if err != nil {
		return nil, err
	}

	pages, err := identityservices.List(client, identityservices.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing services: %w", err)
	}

	items, err := identityservices.ExtractServices(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting services: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, service := range items {
		rows = append(rows, resource.Row{
			ID: service.ID,
			Fields: map[string]string{
				"id":          service.ID,
				"name":        service.Name,
				"type":        service.Type,
				"enabled":     strconv.FormatBool(service.Enabled),
				"description": service.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-e":
		return r.navigateToEndpoints(row)
	}

	return nil
}

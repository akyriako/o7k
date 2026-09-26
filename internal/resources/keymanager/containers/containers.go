package containers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/keymanager"
	tea "github.com/charmbracelet/bubbletea"
	barbicancontainers "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/containers"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{context: openstackContext}
}

func (r *Resource) Kind() string {
	return "secret-containers"
}

func (r *Resource) Title() string {
	return "Secret Containers"
}

func (r *Resource) Aliases() []string {
	return []string{
		"secret-container",
		"scontainers",
		"scontainer",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 16, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "secrets", Title: "SECRETS", MinWidth: 10, Flex: 0},
		{Key: "consumers", Title: "CONSUMERS", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.KeyManagerV1()
	if err != nil {
		return nil, err
	}

	pages, err := barbicancontainers.List(client, nil).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing secret containers: %w", err)
	}

	items, err := barbicancontainers.ExtractContainers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting secret containers: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		id := keymanager.GetReferenceID(item.ContainerRef)

		rows = append(rows, resource.Row{
			ID: id,
			Fields: map[string]string{
				"id":        id,
				"name":      item.Name,
				"type":      item.Type,
				"status":    item.Status,
				"secrets":   strconv.Itoa(len(item.SecretRefs)),
				"consumers": strconv.Itoa(len(item.Consumers)),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	default:
		return nil
	}
}

package secrets

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/keymanager"
	tea "github.com/charmbracelet/bubbletea"
	barbicansecrets "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/secrets"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{context: openstackContext}
}

func (r *Resource) Kind() string {
	return "secrets"
}

func (r *Resource) Title() string {
	return "Secrets"
}

func (r *Resource) Aliases() []string {
	return []string{"secret"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 16, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "algorithm", Title: "ALGORITHM", MinWidth: 12, Flex: 0},
		{Key: "bits", Title: "BITS", MinWidth: 8, Flex: 0},
		{Key: "mode", Title: "MODE", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "p", Description: "Payload"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.KeyManagerV1()
	if err != nil {
		return nil, err
	}

	pages, err := barbicansecrets.List(client, barbicansecrets.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	items, err := barbicansecrets.ExtractSecrets(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting secrets: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		id := keymanager.GetReferenceID(item.SecretRef)

		rows = append(rows, resource.Row{
			ID: id,
			Fields: map[string]string{
				"id":        id,
				"name":      item.Name,
				"type":      item.SecretType,
				"status":    item.Status,
				"algorithm": item.Algorithm,
				"bits":      strconv.Itoa(item.BitLength),
				"mode":      item.Mode,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "p":
		return r.payload(row)
	default:
		return nil
	}
}

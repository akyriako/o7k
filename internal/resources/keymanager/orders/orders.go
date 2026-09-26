package orders

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/keymanager"
	tea "github.com/charmbracelet/bubbletea"
	barbicanorders "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/orders"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{context: openstackContext}
}

func (r *Resource) Kind() string {
	return "orders"
}

func (r *Resource) Title() string {
	return "Orders"
}

func (r *Resource) Aliases() []string {
	return []string{"order"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "type", Title: "TYPE", MinWidth: 14, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "substatus", Title: "SUBSTATUS", MinWidth: 18, Flex: 0},
		{Key: "algorithm", Title: "ALGORITHM", MinWidth: 12, Flex: 0},
		{Key: "bits", Title: "BITS", MinWidth: 8, Flex: 0},
		{Key: "mode", Title: "MODE", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Secret"},
		{Key: "shift-c", Description: "Secret Container"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.KeyManagerV1()
	if err != nil {
		return nil, err
	}

	pages, err := barbicanorders.List(client, nil).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing orders: %w", err)
	}

	items, err := barbicanorders.ExtractOrders(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting orders: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, item := range items {
		id := keymanager.GetReferenceID(item.OrderRef)

		rows = append(rows, resource.Row{
			ID: id,
			Fields: map[string]string{
				"id":           id,
				"type":         item.Type,
				"name":         item.Meta.Name,
				"status":       item.Status,
				"substatus":    item.SubStatus,
				"algorithm":    item.Meta.Algorithm,
				"bits":         strconv.Itoa(item.Meta.BitLength),
				"mode":         item.Meta.Mode,
				"secret_id":    keymanager.GetReferenceID(item.SecretRef),
				"container_id": keymanager.GetReferenceID(item.ContainerRef),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row)
	case "shift-s":
		return r.secret(row)
	case "shift-c":
		return r.container(row)
	default:
		return nil
	}
}

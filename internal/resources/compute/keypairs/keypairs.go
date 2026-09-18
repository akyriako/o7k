package keypairs

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "keypairs"
}

func (r *Resource) Title() string {
	return "Keypairs"
}

func (r *Resource) Aliases() []string {
	return []string{"keypair", "keys", "key"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 12, Flex: 0},
		{Key: "fingerprint", Title: "FINGERPRINT", MinWidth: 48, Flex: 1},
		//{Key: "public_key", Title: "PUBLIC KEY", MinWidth: 40, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := keypairs.List(client, keypairs.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing keypairs: %w", err)
	}

	items, err := keypairs.ExtractKeyPairs(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting keypairs: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, keypair := range items {
		rows = append(rows, resource.Row{
			ID: keypair.Name,
			Fields: map[string]string{
				"name":        keypair.Name,
				"type":        keypair.Type,
				"fingerprint": keypair.Fingerprint,
				"public_key":  keypair.PublicKey,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	}

	return nil
}

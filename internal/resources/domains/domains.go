package domains

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identitydomains "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/domains"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "domains"
}

func (r *Resource) Title() string {
	return "Domains"
}

func (r *Resource) Aliases() []string {
	return []string{"domain"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "enabled", Title: "ENABLED", MinWidth: 10, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 30, Flex: 2},
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

	pages, err := identitydomains.List(client, identitydomains.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing domains: %w", err)
	}

	items, err := identitydomains.ExtractDomains(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting domains: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, domain := range items {
		rows = append(rows, resource.Row{
			ID: domain.ID,
			Fields: map[string]string{
				"id":          domain.ID,
				"name":        domain.Name,
				"enabled":     strconv.FormatBool(domain.Enabled),
				"description": domain.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

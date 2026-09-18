package projects

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	identityprojects "github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "projects"
}

func (r *Resource) Title() string {
	return "Projects"
}

func (r *Resource) Aliases() []string {
	return []string{"project", "prj"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 20, Flex: 1},
		{Key: "domain_id", Title: "DOMAIN ID", MinWidth: 40, Flex: 0},
		{Key: "parent_id", Title: "PARENT ID", MinWidth: 40, Flex: 0},
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

	pages, err := identityprojects.List(client, identityprojects.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}

	items, err := identityprojects.ExtractProjects(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting projects: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, project := range items {
		rows = append(rows, resource.Row{
			ID: project.ID,
			Fields: map[string]string{
				"id":          project.ID,
				"name":        project.Name,
				"domain_id":   project.DomainID,
				"parent_id":   project.ParentID,
				"enabled":     strconv.FormatBool(project.Enabled),
				"description": project.Description,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

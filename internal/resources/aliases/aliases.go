package aliases

import (
	"context"
	"slices"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	registry *resource.Registry
}

func New(r *resource.Registry) *Resource {
	return &Resource{registry: r}
}

func (r *Resource) Kind() string {
	return "commands"
}

func (r *Resource) Title() string {
	return "Commands"
}

func (r *Resource) Aliases() []string {
	return []string{"aliases", "alias", "command", "cmd"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "TITLE", MinWidth: 20, Flex: 1},
		{Key: "kind", Title: "KIND", MinWidth: 20, Flex: 1},
		{Key: "aliases", Title: "ALIASES", MinWidth: 20, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(context.Context) ([]resource.Row, error) {
	registeredResources := r.registry.Resources()
	slices.SortFunc(registeredResources, func(a, b resource.Resource) int {
		return strings.Compare(a.Kind(), b.Kind())
	})

	rows := make([]resource.Row, 0, len(registeredResources))

	for _, rr := range registeredResources {
		rows = append(rows, resource.Row{
			ID: rr.Kind(),
			Fields: map[string]string{
				"kind":    rr.Kind(),
				"name":    rr.Title(),
				"aliases": strings.Join(rr.Aliases(), ", "),
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

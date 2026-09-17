package servers

import (
	"context"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct{}

var _ resource.Resource = (*Resource)(nil)

func New() *Resource {
	return &Resource{}
}

func (r *Resource) Title() string {
	return "Servers"
}

func (r *Resource) Kind() string {
	return "servers"
}

func (r *Resource) Aliases() []string {
	return []string{
		"server",
		"srv",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "name", Title: "NAME", MinWidth: 16, Flex: 3},
		{Key: "status", Title: "STATUS", MinWidth: 10, Flex: 1},
		{Key: "flavor", Title: "FLAVOR", MinWidth: 12, Flex: 2},
		{Key: "ip", Title: "IP", MinWidth: 15, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	return []resource.Row{
		{
			ID: "server-1",
			Fields: map[string]string{
				"name":   "api-01",
				"status": "ACTIVE",
				"flavor": "m1.large",
				"ip":     "10.0.0.4",
			},
		},
		{
			ID: "server-2",
			Fields: map[string]string{
				"name":   "api-02",
				"status": "ACTIVE",
				"flavor": "m1.large",
				"ip":     "10.0.0.5",
			},
		},
		{
			ID: "server-3",
			Fields: map[string]string{
				"name":   "worker-01",
				"status": "SHUTOFF",
				"flavor": "m1.xlarge",
				"ip":     "10.0.0.8",
			},
		},
	}, nil
}

func (r *Resource) Execute(_ resource.Command, _ resource.Row) tea.Cmd {
	return nil
}

package networks

import (
	"context"

	"github.com/akyriako/o7k/internal/resource"
)

type Resource struct{}

var _ resource.Resource = (*Resource)(nil)

func New() *Resource {
	return &Resource{}
}

func (r *Resource) Title() string {
	return "Networks"
}

func (r *Resource) Kind() string {
	return "networks"
}

func (r *Resource) Aliases() []string {
	return []string{
		"network",
		"net",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{
			Key:      "name",
			Title:    "NAME",
			MinWidth: 20,
			Flex:     3,
		},
		{
			Key:      "status",
			Title:    "STATUS",
			MinWidth: 10,
			Flex:     1,
		},
		{
			Key:      "shared",
			Title:    "SHARED",
			MinWidth: 10,
			Flex:     1,
		},
		{
			Key:      "external",
			Title:    "EXTERNAL",
			MinWidth: 10,
			Flex:     1,
		},
	}
}

func (r *Resource) List(_ context.Context) ([]resource.Row, error) {
	return []resource.Row{
		{
			ID: "network-1",
			Fields: map[string]string{
				"name":     "public",
				"status":   "ACTIVE",
				"shared":   "true",
				"external": "true",
			},
		},
		{
			ID: "network-2",
			Fields: map[string]string{
				"name":     "private",
				"status":   "ACTIVE",
				"shared":   "false",
				"external": "false",
			},
		},
		{
			ID: "network-3",
			Fields: map[string]string{
				"name":     "storage",
				"status":   "ACTIVE",
				"shared":   "true",
				"external": "false",
			},
		},
	}, nil
}

package sharetypes

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	manilasharetypes "github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Types"
}

func (r *Resource) Kind() string {
	return "share-types"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-type",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "is_public", Title: "PUBLIC", MinWidth: 8, Flex: 0},
		{Key: "is_default", Title: "DEFAULT", MinWidth: 8, Flex: 0},
		{Key: "description", Title: "DESCRIPTION", MinWidth: 32, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	pages, err := manilasharetypes.List(client, manilasharetypes.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing share types: %w", err)
	}

	items, err := manilasharetypes.ExtractShareTypes(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting share types: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, shareType := range items {
		description := ""
		if shareType.Description != nil {
			description = *shareType.Description
		}

		isDefault := ""
		if shareType.IsDefault != nil {
			isDefault = strconv.FormatBool(*shareType.IsDefault)
		}

		rows = append(rows, resource.Row{
			ID: shareType.ID,
			Fields: map[string]string{
				"id":          shareType.ID,
				"name":        shareType.Name,
				"is_public":   strconv.FormatBool(shareType.IsPublic),
				"is_default":  isDefault,
				"description": description,
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

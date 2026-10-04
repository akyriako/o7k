package accessrules

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shareaccessrules"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Manila Share Access Rules"
}

func (r *Resource) Kind() string {
	return "share-access-rules"
}

func (r *Resource) Aliases() []string {
	return []string{
		"share-access-rule",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "access_type", Title: "TYPE", MinWidth: 10, Flex: 0},
		{Key: "access_to", Title: "ACCESS TO", MinWidth: 24, Flex: 1},
		{Key: "access_level", Title: "LEVEL", MinWidth: 10, Flex: 0},
		{Key: "state", Title: "STATE", MinWidth: 12, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	shareID := scope["share_id"]

	if shareID == "" {
		return nil, fmt.Errorf("share ID is required")
	}

	client, err := r.context.SharedFileSystemV2()
	if err != nil {
		return nil, err
	}

	result := shareaccessrules.List(ctx, client, shareID)

	items, err := result.Extract()
	if err != nil {
		return nil, fmt.Errorf("listing share access rules: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, rule := range items {
		rows = append(rows, resource.Row{
			ID: rule.ID,
			Fields: map[string]string{
				"id":           rule.ID,
				"share_id":     rule.ShareID,
				"access_type":  rule.AccessType,
				"access_to":    rule.AccessTo,
				"access_level": rule.AccessLevel,
				"state":        rule.State,
				"access_key":   rule.AccessKey,
			},
		})
	}

	return rows, nil
}
func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

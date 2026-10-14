package backups

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type Resource struct {
	context *openstack.Context
}

type Backup struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	InstanceID  string  `json:"instance_id"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
	Size        float64 `json:"size"`
	ParentID    string  `json:"parent_id"`
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "db-backups"
}

func (r *Resource) Title() string {
	return "Trove Backups"
}

func (r *Resource) Aliases() []string {
	return []string{"dbbackup", "trovebackups"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 3},
		{Key: "status", Title: "STATUS", MinWidth: 14, Flex: 0},
		{Key: "instance_id", Title: "INSTANCE ID", MinWidth: 36, Flex: 1},
		{Key: "size", Title: "SIZE (GB)", MinWidth: 12, Flex: 0},
		{Key: "created", Title: "CREATED", MinWidth: 24, Flex: 1},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DBV1()
	if err != nil {
		return nil, err
	}

	rows := make([]resource.Row, 0)
	nextURL := client.ServiceURL("backups")
	visited := make(map[string]bool)

	for nextURL != "" {
		if visited[nextURL] {
			return nil, fmt.Errorf("backup pagination loop detected: %s", nextURL)
		}
		visited[nextURL] = true

		var result struct {
			Backups []Backup `json:"backups"`
			Links   []struct {
				Href string `json:"href"`
				Rel  string `json:"rel"`
			} `json:"backups_links"`
		}

		_, err := client.Get(ctx, nextURL, &result, nil)
		if err != nil {
			return nil, fmt.Errorf("listing database backups: %w", err)
		}

		for _, backup := range result.Backups {
			rows = append(rows, resource.Row{
				ID: backup.ID,
				Fields: map[string]string{
					"id":          backup.ID,
					"name":        backup.Name,
					"description": backup.Description,
					"status":      backup.Status,
					"instance_id": backup.InstanceID,
					"size":        strconv.FormatFloat(backup.Size, 'f', 2, 64),
					"created":     backup.Created,
					"updated":     backup.Updated,
					"parent_id":   backup.ParentID,
				},
			})
		}

		currentURL, err := url.Parse(nextURL)
		if err != nil {
			return nil, fmt.Errorf("parsing backup page URL: %w", err)
		}

		nextURL = ""

		for _, link := range result.Links {
			if link.Rel != "next" {
				continue
			}

			target, err := url.Parse(link.Href)
			if err != nil {
				return nil, fmt.Errorf("parsing next backup page: %w", err)
			}

			resolved := currentURL.ResolveReference(target)

			if resolved.Scheme != currentURL.Scheme ||
				resolved.Host != currentURL.Host ||
				resolved.Path != currentURL.Path {
				return nil, fmt.Errorf("unexpected backup pagination URL: %s", resolved)
			}

			nextURL = resolved.String()
			break
		}
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

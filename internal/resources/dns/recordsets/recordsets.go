package recordsets

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	dnsrecordsets "github.com/gophercloud/gophercloud/v2/openstack/dns/v2/recordsets"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "DNS Recordsets"
}

func (r *Resource) Kind() string {
	return "recordsets"
}

func (r *Resource) Aliases() []string {
	return []string{
		"recordset",
		"records",
		"record",
	}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 30, Flex: 1},
		{Key: "type", Title: "TYPE", MinWidth: 8, Flex: 0},
		{Key: "records", Title: "RECORDS", MinWidth: 30, Flex: 2},
		{Key: "ttl", Title: "TTL", MinWidth: 8, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "zone_id", Title: "ZONE ID", MinWidth: 40, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "shift-z", Description: "Zone"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.DNSV2()
	if err != nil {
		return nil, err
	}

	pages, err := dnsrecordsets.ListAll(
		client,
		dnsrecordsets.ListOpts{},
	).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing DNS recordsets: %w", err)
	}

	items, err := dnsrecordsets.ExtractRecordSets(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting DNS recordsets: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, recordset := range items {
		rows = append(rows, resource.Row{
			ID: recordset.ID,
			Fields: map[string]string{
				"id":      recordset.ID,
				"name":    recordset.Name,
				"type":    recordset.Type,
				"records": strings.Join(recordset.Records, ", "),
				"ttl":     strconv.Itoa(recordset.TTL),
				"status":  recordset.Status,
				"zone_id": recordset.ZoneID,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "shift-z":
		return r.navigateToZone(row)
	}

	return nil
}

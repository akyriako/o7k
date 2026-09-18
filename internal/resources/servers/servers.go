package servers

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeflavors "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
	computeservers "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
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
		{Key: "id", Title: "ID", MinWidth: 40, Flex: 0},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12, Flex: 0},
		{Key: "flavor", Title: "FLAVOR", MinWidth: 15, Flex: 1},
		{Key: "image", Title: "IMAGE", MinWidth: 40, Flex: 0},
		{Key: "addresses", Title: "ADDRESSES", MinWidth: 40, Flex: 2},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-i", Description: "Image"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := computeservers.List(client, computeservers.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing servers: %w", err)
	}

	items, err := computeservers.ExtractServers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting servers: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	flavorPages, err := computeflavors.ListDetail(client, computeflavors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing flavors: %w", err)
	}

	flavors, err := computeflavors.ExtractFlavors(flavorPages)
	if err != nil {
		return nil, fmt.Errorf("extracting flavors: %w", err)
	}

	flavorNames := make(map[string]string, len(flavors))

	for _, flavor := range flavors {
		flavorNames[flavor.ID] = flavor.Name
	}

	for _, server := range items {
		flavorID, _ := server.Flavor["id"].(string)
		flavorName := flavorNames[flavorID]

		if flavorName == "" {
			flavorName = flavorID
		}

		imageID, _ := server.Image["id"].(string)

		image := imageID
		if name, ok := server.Image["name"].(string); ok && name != "" {
			image = name
		}

		rows = append(rows, resource.Row{
			ID: server.ID,
			Fields: map[string]string{
				"id":        server.ID,
				"name":      server.Name,
				"status":    server.Status,
				"flavor":    flavorName,
				"addresses": serverAddresses(server.Addresses),
				"image":     image,
				"image_id":  imageID,
			},
		})
	}

	return rows, nil
}

func serverAddresses(addresses map[string]any) string {
	values := make([]string, 0)

	for network, raw := range addresses {
		items, ok := raw.([]any)
		if !ok {
			continue
		}

		for _, rawAddress := range items {
			address, ok := rawAddress.(map[string]any)
			if !ok {
				continue
			}

			ip, ok := address["addr"].(string)
			if !ok || ip == "" {
				continue
			}

			values = append(values, network+"="+ip)
		}
	}

	return strings.Join(values, ", ")
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-i":
		return r.navigateToImage(row)
	}

	return nil
}

package hypervisors

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/hypervisors"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Kind() string {
	return "hypervisors"
}

func (r *Resource) Title() string {
	return "Compute Hypervisors"
}

func (r *Resource) Aliases() []string {
	return []string{"hypervisor", "hyp", "hyps"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "id", Title: "ID", MinWidth: 36, Flex: 0},
		{Key: "hostname", Title: "HOSTNAME", MinWidth: 24, Flex: 2},
		{Key: "type", Title: "TYPE", MinWidth: 12, Flex: 1},
		{Key: "host_ip", Title: "HOST IP", MinWidth: 15, Flex: 0},
		{Key: "running_vms", Title: "VMs", MinWidth: 6, Flex: 0},
		{Key: "vcpus_used", Title: "USED VCPUs", MinWidth: 11, Flex: 0},
		{Key: "vcpus", Title: "VCPUs", MinWidth: 7, Flex: 0},
		{Key: "memory_mb_used", Title: "USED RAM (MB)", MinWidth: 14, Flex: 0},
		{Key: "memory_mb", Title: "RAM (MB)", MinWidth: 10, Flex: 0},
		{Key: "state", Title: "STATE", MinWidth: 8, Flex: 0},
		{Key: "status", Title: "STATUS", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-s", Description: "Servers"},
	}
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ComputeV2()
	if err != nil {
		return nil, err
	}

	pages, err := hypervisors.List(client, hypervisors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing hypervisors: %w", err)
	}

	items, err := hypervisors.ExtractHypervisors(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting hypervisors: %w", err)
	}

	rows := make([]resource.Row, 0, len(items))

	for _, hypervisor := range items {
		rows = append(rows, resource.Row{
			ID: hypervisor.ID,
			Fields: map[string]string{
				"id":             hypervisor.ID,
				"hostname":       hypervisor.HypervisorHostname,
				"host":           hypervisor.Service.Host,
				"type":           hypervisor.HypervisorType,
				"state":          hypervisor.State,
				"status":         hypervisor.Status,
				"running_vms":    strconv.Itoa(hypervisor.RunningVMs),
				"vcpus":          strconv.Itoa(hypervisor.VCPUs),
				"vcpus_used":     strconv.Itoa(hypervisor.VCPUsUsed),
				"memory_mb":      strconv.Itoa(hypervisor.MemoryMB),
				"memory_mb_used": strconv.Itoa(hypervisor.MemoryMBUsed),
				"host_ip":        hypervisor.HostIP,
			},
		})
	}

	return rows, nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	switch command.Key {
	case "s":
		return r.show(row.ID)
	case "shift-s":
		return r.servers(row)
	}

	return nil
}

package quotas

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/quotasets"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/limits"
	networkquotas "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
)

type Resource struct {
	context *openstack.Context
}

func New(context *openstack.Context) *Resource {
	return &Resource{context: context}
}

func (r *Resource) Title() string {
	return "Project Quotas"
}

func (r *Resource) Kind() string {
	return "quotas"
}

func (r *Resource) Aliases() []string {
	return []string{"quota"}
}

func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{Key: "service", Title: "SERVICE", MinWidth: 12, Flex: 0},
		{Key: "resource", Title: "RESOURCE", MinWidth: 24, Flex: 1},
		{Key: "used", Title: "USED", MinWidth: 10, Flex: 0},
		{Key: "reserved", Title: "RESERVED", MinWidth: 10, Flex: 0},
		{Key: "limit", Title: "LIMIT", MinWidth: 10, Flex: 0},
	}
}

func (r *Resource) Commands() []resource.Command {
	return nil
}

func (r *Resource) Execute(command resource.Command, row resource.Row) tea.Cmd {
	return nil
}

func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	projectID := scope["project_id"]

	if projectID == "" {
		return nil, fmt.Errorf("project ID is required")
	}

	var rows []resource.Row
	var errors []error

	if err := r.appendComputeQuotas(ctx, projectID, &rows); err != nil {
		errors = append(errors, err)
	}

	if err := r.appendNetworkQuotas(ctx, projectID, &rows); err != nil {
		errors = append(errors, err)
	}

	if err := r.appendBlockStorageQuotas(ctx, projectID, &rows); err != nil {
		errors = append(errors, err)
	}

	if len(errors) == 3 {
		return nil, fmt.Errorf("getting project quotas: %v", errors)
	}

	return rows, nil
}

func appendRow(rows *[]resource.Row, service, name string, used, reserved, limit int) {
	*rows = append(*rows, resource.Row{
		ID: service + ":" + name,
		Fields: map[string]string{
			"service":  service,
			"resource": name,
			"used":     strconv.Itoa(used),
			"reserved": strconv.Itoa(reserved),
			"limit":    quotaLimit(limit),
		},
	})
}

func (r *Resource) appendComputeQuotas(ctx context.Context, projectID string, rows *[]resource.Row) error {
	client, err := r.context.ComputeV2()
	if err != nil {
		return fmt.Errorf("creating compute client: %w", err)
	}

	result, err := limits.Get(ctx, client, limits.GetOpts{TenantID: projectID}).Extract()
	if err != nil {
		return fmt.Errorf("getting compute quotas: %w", err)
	}

	absolute := result.Absolute

	appendRow(rows, "Nova", "Instances", absolute.TotalInstancesUsed, 0, absolute.MaxTotalInstances)
	appendRow(rows, "Nova", "Cores", absolute.TotalCoresUsed, 0, absolute.MaxTotalCores)
	appendRow(rows, "Nova", "RAM (MB)", absolute.TotalRAMUsed, 0, absolute.MaxTotalRAMSize)
	appendRow(rows, "Nova", "Floating IPs", absolute.TotalFloatingIpsUsed, 0, absolute.MaxTotalFloatingIps)
	appendRow(rows, "Nova", "Security Groups", absolute.TotalSecurityGroupsUsed, 0, absolute.MaxSecurityGroups)
	appendRow(rows, "Nova", "Server Groups", absolute.TotalServerGroupsUsed, 0, absolute.MaxServerGroups)

	return nil
}

func (r *Resource) appendNetworkQuotas(ctx context.Context, projectID string, rows *[]resource.Row) error {
	client, err := r.context.NetworkV2()
	if err != nil {
		return fmt.Errorf("creating network client: %w", err)
	}

	result, err := networkquotas.GetDetail(ctx, client, projectID).Extract()
	if err != nil {
		return fmt.Errorf("getting network quotas: %w", err)
	}

	appendRow(rows, "Neutron", "Networks", result.Network.Used, result.Network.Reserved, result.Network.Limit)
	appendRow(rows, "Neutron", "Subnets", result.Subnet.Used, result.Subnet.Reserved, result.Subnet.Limit)
	appendRow(rows, "Neutron", "Subnet Pools", result.SubnetPool.Used, result.SubnetPool.Reserved, result.SubnetPool.Limit)
	appendRow(rows, "Neutron", "Ports", result.Port.Used, result.Port.Reserved, result.Port.Limit)
	appendRow(rows, "Neutron", "Routers", result.Router.Used, result.Router.Reserved, result.Router.Limit)
	appendRow(rows, "Neutron", "Floating IPs", result.FloatingIP.Used, result.FloatingIP.Reserved, result.FloatingIP.Limit)
	appendRow(rows, "Neutron", "Security Groups", result.SecurityGroup.Used, result.SecurityGroup.Reserved, result.SecurityGroup.Limit)
	appendRow(rows, "Neutron", "Security Group Rules", result.SecurityGroupRule.Used, result.SecurityGroupRule.Reserved, result.SecurityGroupRule.Limit)
	appendRow(rows, "Neutron", "RBAC Policies", result.RBACPolicy.Used, result.RBACPolicy.Reserved, result.RBACPolicy.Limit)
	appendRow(rows, "Neutron", "Trunks", result.Trunk.Used, result.Trunk.Reserved, result.Trunk.Limit)

	return nil
}

func (r *Resource) appendBlockStorageQuotas(ctx context.Context, projectID string, rows *[]resource.Row) error {
	client, err := r.context.BlockStorageV3()
	if err != nil {
		return fmt.Errorf("creating block storage client: %w", err)
	}

	result, err := quotasets.GetUsage(ctx, client, projectID).Extract()
	if err != nil {
		return fmt.Errorf("getting block storage quotas: %w", err)
	}

	appendRow(rows, "Cinder", "Volumes", result.Volumes.InUse, result.Volumes.Reserved, result.Volumes.Limit)
	appendRow(rows, "Cinder", "Snapshots", result.Snapshots.InUse, result.Snapshots.Reserved, result.Snapshots.Limit)
	appendRow(rows, "Cinder", "Gigabytes", result.Gigabytes.InUse, result.Gigabytes.Reserved, result.Gigabytes.Limit)
	appendRow(rows, "Cinder", "Per Volume Gigabytes", result.PerVolumeGigabytes.InUse, result.PerVolumeGigabytes.Reserved, result.PerVolumeGigabytes.Limit)
	appendRow(rows, "Cinder", "Backups", result.Backups.InUse, result.Backups.Reserved, result.Backups.Limit)
	appendRow(rows, "Cinder", "Backup Gigabytes", result.BackupGigabytes.InUse, result.BackupGigabytes.Reserved, result.BackupGigabytes.Limit)
	appendRow(rows, "Cinder", "Groups", result.Groups.InUse, result.Groups.Reserved, result.Groups.Limit)

	return nil
}

func quotaLimit(value int) string {
	if value == -1 {
		return ""
	}

	return strconv.Itoa(value)
}

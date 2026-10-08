package servers

import (
	"context"
	"fmt"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	computeservers "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
)

func (r *Resource) image(row resource.Row) tea.Cmd {
	imageID := row.Fields["image_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "images",
			Field:    "id",
			Value:    imageID,
		}
	}
}

func (r *Resource) flavor(row resource.Row) tea.Cmd {
	flavor := row.Fields["flavor"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "flavors",
			Field:    "name",
			Value:    flavor,
		}
	}
}

// securityGroups violates the "no N+1 calls" rule, but this is a special case because we
// have to first load the ports and then find out the corresponding security group ids.
// this method should not be considered as the norm of jumping to resources
func (r *Resource) securityGroups(row resource.Row) tea.Cmd {
	return func() tea.Msg {

		client, err := r.context.NetworkV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		pages, err := ports.List(client, ports.ListOpts{
			DeviceID: row.ID,
		}).AllPages(context.Background())
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("listing ports for server %q: %w", row.ID, err),
			}
		}

		items, err := ports.ExtractPorts(pages)
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("extracting ports for server %q: %w", row.ID, err),
			}
		}

		uniqueSecurityGroups := make(map[string]struct{})
		securityGroupIDs := make([]string, 0)

		for _, port := range items {
			for _, securityGroupID := range port.SecurityGroups {
				if _, exists := uniqueSecurityGroups[securityGroupID]; exists {
					continue
				}

				uniqueSecurityGroups[securityGroupID] = struct{}{}
				securityGroupIDs = append(securityGroupIDs, securityGroupID)
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "securitygroups",
			Field:    "id",
			Values:   securityGroupIDs,
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		server, err := computeservers.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting server %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      server.ID,
			Content: server,
		}
	}
}

func (r *Resource) volumes(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		server, err := computeservers.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("getting server %q: %w", row.ID, err),
			}
		}

		volumeIDs := make([]string, 0, len(server.AttachedVolumes))

		for _, volume := range server.AttachedVolumes {
			if volume.ID != "" {
				volumeIDs = append(volumeIDs, volume.ID)
			}
		}

		if len(volumeIDs) == 0 {
			return resource.ErrorMsg{
				Err: fmt.Errorf("server %q has no attached volumes", row.ID),
			}
		}

		return resource.NavigateFilteredMultiMsg{
			Resource: "volumes",
			Field:    "id",
			Values:   volumeIDs,
		}
	}
}

func (r *Resource) start(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		if err := validatePowerAction("start", row.Fields["status"]); err != nil {
			return resource.ErrorMsg{Err: err}
		}

		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		err = computeservers.Start(
			context.Background(),
			client,
			row.ID,
		).ExtractErr()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("starting server %q: %w", row.ID, err),
			}
		}

		return resource.CommandCompletedMsg{RefreshOnCompletion: true}
	}
}

func (r *Resource) stop(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		if err := validatePowerAction("stop", row.Fields["status"]); err != nil {
			return resource.ErrorMsg{Err: err}
		}

		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		err = computeservers.Stop(
			context.Background(),
			client,
			row.ID,
		).ExtractErr()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("stopping server %q: %w", row.ID, err),
			}
		}

		return resource.CommandCompletedMsg{RefreshOnCompletion: true}
	}
}

func (r *Resource) reboot(row resource.Row, soft bool) tea.Cmd {
	return func() tea.Msg {
		if err := validatePowerAction("reboot", row.Fields["status"]); err != nil {
			return resource.ErrorMsg{Err: err}
		}

		client, err := r.context.ComputeV2()
		if err != nil {
			return resource.ErrorMsg{Err: err}
		}

		rebootType := computeservers.SoftReboot
		if !soft {
			rebootType = computeservers.HardReboot
		}
		err = computeservers.Reboot(
			context.Background(),
			client,
			row.ID,
			computeservers.RebootOpts{
				Type: rebootType,
			},
		).ExtractErr()
		if err != nil {
			return resource.ErrorMsg{
				Err: fmt.Errorf("rebooting server %q: %w", row.ID, err),
			}
		}

		return resource.CommandCompletedMsg{RefreshOnCompletion: true}
	}
}

func validatePowerAction(action, status string) error {
	status = strings.ToUpper(strings.TrimSpace(status))

	switch action {
	case "start":
		if status == "SHUTOFF" {
			return nil
		}
	case "stop":
		if status == "ACTIVE" {
			return nil
		}
	case "reboot":
		if status == "ACTIVE" {
			return nil
		}
	}

	return fmt.Errorf(
		"cannot %s server in %s state",
		action,
		status,
	)
}

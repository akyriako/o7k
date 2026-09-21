package contexts

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
)

type ActivatedMsg struct {
	Context *openstack.Context
	Err     error
}

func (r *Resource) activate(row resource.Row) tea.Cmd {
	cloudsPaths := r.cloudsPaths
	cloudName := row.ID

	return func() tea.Msg {
		clouds, err := openstack.LoadClouds(cloudsPaths)
		if err != nil {
			return ActivatedMsg{Err: err}
		}

		cloud, ok := clouds.Get(cloudName)
		if !ok {
			return ActivatedMsg{
				Err: fmt.Errorf("cloud not found: %s", cloudName),
			}
		}

		next := openstack.Context{
			Cloud:    cloud.Name,
			Identity: cloud.Identity,
			Region:   cloud.Region,
			Domain:   cloud.Domain,
			Project:  cloud.Project,
		}

		err = next.Connect(context.Background(), cloud.Path)

		return ActivatedMsg{
			Context: &next,
			Err:     err,
		}
	}
}

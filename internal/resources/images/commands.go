package images

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
)

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ImageV2()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		image, err := images.Get(context.Background(), client, id).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting image %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      image.ID,
			Content: image,
		}
	}
}

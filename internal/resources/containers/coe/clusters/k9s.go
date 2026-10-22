package clusters

import (
	"context"
	"fmt"
	"os"

	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/terminal"
	tea "github.com/charmbracelet/bubbletea"
)

func (r *Resource) openK9s(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		if err := terminal.CheckK9s(); err != nil {
			return resource.DetailsMsg{Err: err}
		}

		data, err := r.getKubeconfigData(context.Background(), row.ID)
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		f, err := os.CreateTemp("", "o7k-generated_kubeconfig-*.yaml")
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("creating temporary kubeconfig failed: %w", err),
			}
		}

		path := f.Name()
		if _, err := f.Write(data); err != nil {
			f.Close()
			os.Remove(path)
			return resource.DetailsMsg{
				Err: fmt.Errorf("writing temporary kubeconfig failed: %w", err),
			}
		}

		if err := f.Close(); err != nil {
			os.Remove(path)
			return resource.DetailsMsg{
				Err: fmt.Errorf("closing temporary kubeconfig failed: %w", err),
			}
		}

		if err := terminal.OpenK9s(path); err != nil {
			os.Remove(path)
			return resource.DetailsMsg{
				Err: fmt.Errorf("launching K9s failed: %w", err),
			}
		}

		return nil
	}
}

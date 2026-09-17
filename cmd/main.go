package main

import (
	"fmt"
	"os"

	"github.com/akyriako/o7k/internal/logging"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/networks"
	"github.com/akyriako/o7k/internal/resources/servers"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/akyriako/o7k/internal/ui"
)

func init() {

}

func main() {
	logger, logFile, err := logging.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing logging: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	logger.Info("starting o7k")

	registry := resource.NewRegistry()

	if err := registry.Register(servers.New()); err != nil {
		fmt.Fprintf(os.Stderr, "error registering servers resource: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Register(networks.New()); err != nil {
		fmt.Fprintf(os.Stderr, "error registering networks resource: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(
		ui.New(registry),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running o7k: %v\n", err)
		os.Exit(1)
	}
}

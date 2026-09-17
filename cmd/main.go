package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akyriako/o7k/internal/logging"
	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/contexts"
	"github.com/akyriako/o7k/internal/resources/networks"
	"github.com/akyriako/o7k/internal/resources/servers"
	"github.com/akyriako/o7k/internal/resources/services"
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

	cloudsPath, err := openstack.DiscoverCloudsFile("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error discovering clouds.yaml: %v\n", err)
		os.Exit(1)
	}

	availableClouds, err := openstack.LoadClouds(cloudsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading clouds.yaml: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Register(contexts.New(availableClouds)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering contexts resource: %v\n", err)
		os.Exit(1)
	}

	cloud := availableClouds.Items[0]

	openstackContext := openstack.Context{
		Cloud:    cloud.Name,
		Region:   cloud.Region,
		Project:  cloud.Project,
		Domain:   cloud.Domain,
		Identity: cloud.Identity,
	}

	if err := openstackContext.Connect(context.Background(), cloudsPath); err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to OpenStack: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Register(services.New(&openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering services resource: %v\n", err)
		os.Exit(1)
	}

	logger.Info("connected to OpenStack", "cloud", openstackContext.Cloud)

	p := tea.NewProgram(
		ui.New(
			registry,
			&openstackContext,
		),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running o7k: %v\n", err)
		os.Exit(1)
	}
}

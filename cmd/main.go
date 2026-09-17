package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akyriako/o7k/internal/logging"
	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/contexts"
	"github.com/akyriako/o7k/internal/resources/domains"
	"github.com/akyriako/o7k/internal/resources/endpoints"
	"github.com/akyriako/o7k/internal/resources/flavors"
	"github.com/akyriako/o7k/internal/resources/projects"
	"github.com/akyriako/o7k/internal/resources/regions"
	"github.com/akyriako/o7k/internal/resources/roles"
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

	if err := registry.Register(contexts.New(availableClouds.Path)); err != nil {
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

	logger.Info("connected to OpenStack", "cloud", openstackContext.Cloud)

	if err := openstackContext.Connect(context.Background(), cloudsPath); err != nil {
		fmt.Fprintf(os.Stderr, "error connecting to OpenStack: %v\n", err)
		os.Exit(1)
	}

	registerAll(registry, &openstackContext)

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

func registerAll(r *resource.Registry, openstackContext *openstack.Context) {
	if err := r.Register(services.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering services resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(endpoints.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering endpoints resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(regions.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering regions resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(projects.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering projects resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(domains.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering domains resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(roles.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering roles resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(servers.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering servers resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(flavors.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering flavors resource: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/akyriako/o7k/internal/logging"
	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
	"github.com/akyriako/o7k/internal/resources/blockstorage/backups"
	"github.com/akyriako/o7k/internal/resources/blockstorage/snapshots"
	"github.com/akyriako/o7k/internal/resources/blockstorage/volumes"
	"github.com/akyriako/o7k/internal/resources/blockstorage/volumetypes"
	"github.com/akyriako/o7k/internal/resources/compute/availabilityzones"
	"github.com/akyriako/o7k/internal/resources/compute/flavors"
	"github.com/akyriako/o7k/internal/resources/compute/keypairs"
	"github.com/akyriako/o7k/internal/resources/compute/servergroups"
	"github.com/akyriako/o7k/internal/resources/compute/servers"
	"github.com/akyriako/o7k/internal/resources/contexts"
	"github.com/akyriako/o7k/internal/resources/identity/catalog"
	"github.com/akyriako/o7k/internal/resources/identity/domains"
	"github.com/akyriako/o7k/internal/resources/identity/endpoints"
	"github.com/akyriako/o7k/internal/resources/identity/projects"
	"github.com/akyriako/o7k/internal/resources/identity/regions"
	"github.com/akyriako/o7k/internal/resources/identity/roles"
	"github.com/akyriako/o7k/internal/resources/identity/services"
	"github.com/akyriako/o7k/internal/resources/images"
	"github.com/akyriako/o7k/internal/resources/networking/floatingips"
	"github.com/akyriako/o7k/internal/resources/networking/networks"
	"github.com/akyriako/o7k/internal/resources/networking/ports"
	"github.com/akyriako/o7k/internal/resources/networking/routers"
	"github.com/akyriako/o7k/internal/resources/networking/securitygrouprules"
	"github.com/akyriako/o7k/internal/resources/networking/securitygroups"
	"github.com/akyriako/o7k/internal/resources/networking/subnets"
	"github.com/akyriako/o7k/internal/version"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/akyriako/o7k/internal/ui"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		info := version.GetBuildInfo()

		fmt.Printf("o7k %s\n", info.Version)
		fmt.Printf("commit: %s\n", info.Commit)
		fmt.Printf("built: %s\n", info.BuildDate)
		fmt.Printf("go: %s\n", info.GoVersion)
		fmt.Printf("modified: %s\n", info.Modified)
		return
	}

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
	if err := r.Register(catalog.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering catalog resource: %v\n", err)
		os.Exit(1)
	}

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

	if err := r.Register(networks.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering networks resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(subnets.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering subnets resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(ports.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering ports resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(routers.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering routers resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(floatingips.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering floating IPs resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(securitygroups.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering security groups resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(securitygrouprules.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering security group rules resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(images.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering images resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(volumes.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering volumes resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(snapshots.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering snapshots resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(volumetypes.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering volume types resource: %v\n", err)
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

	if err := r.Register(keypairs.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering keypairs resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(servergroups.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering server groups resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(availabilityzones.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering availability zones resource: %v\n", err)
		os.Exit(1)
	}

	if err := r.Register(backups.New(openstackContext)); err != nil {
		fmt.Fprintf(os.Stderr, "error registering volume backups resource: %v\n", err)
		os.Exit(1)
	}
}

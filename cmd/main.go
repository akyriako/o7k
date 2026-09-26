package main

import (
	"context"
	"errors"
	"fmt"
	"io"
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
	"github.com/akyriako/o7k/internal/resources/dns/recordsets"
	"github.com/akyriako/o7k/internal/resources/dns/zones"
	"github.com/akyriako/o7k/internal/resources/identity/catalog"
	"github.com/akyriako/o7k/internal/resources/identity/domains"
	"github.com/akyriako/o7k/internal/resources/identity/endpoints"
	"github.com/akyriako/o7k/internal/resources/identity/projects"
	"github.com/akyriako/o7k/internal/resources/identity/regions"
	"github.com/akyriako/o7k/internal/resources/identity/roles"
	"github.com/akyriako/o7k/internal/resources/identity/services"
	"github.com/akyriako/o7k/internal/resources/images"
	"github.com/akyriako/o7k/internal/resources/keymanager/containers"
	"github.com/akyriako/o7k/internal/resources/keymanager/orders"
	"github.com/akyriako/o7k/internal/resources/keymanager/secrets"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/healthmonitors"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/l7policies"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/l7rules"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/listeners"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/loadbalancers"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/members"
	"github.com/akyriako/o7k/internal/resources/loadbalancing/pools"
	"github.com/akyriako/o7k/internal/resources/networking/floatingips"
	"github.com/akyriako/o7k/internal/resources/networking/networks"
	"github.com/akyriako/o7k/internal/resources/networking/ports"
	"github.com/akyriako/o7k/internal/resources/networking/routers"
	"github.com/akyriako/o7k/internal/resources/networking/securitygrouprules"
	"github.com/akyriako/o7k/internal/resources/networking/securitygroups"
	"github.com/akyriako/o7k/internal/resources/networking/subnets"
	swiftcontainers "github.com/akyriako/o7k/internal/resources/objectstorage/containers"
	swiftobjects "github.com/akyriako/o7k/internal/resources/objectstorage/objects"
	"github.com/akyriako/o7k/internal/resources/orchestration/stackevents"
	"github.com/akyriako/o7k/internal/resources/orchestration/stackresources"
	"github.com/akyriako/o7k/internal/resources/orchestration/stacks"
	"github.com/akyriako/o7k/internal/version"
	tea "github.com/charmbracelet/bubbletea"

	"log/slog"

	"github.com/akyriako/o7k/internal/ui"
)

var (
	logger  *slog.Logger
	stderr  io.Writer = os.Stderr
	logFile io.Closer
)

func main() {
	info := version.GetBuildInfo()
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("o7k %s\n", info.Version)
		fmt.Printf("commit: %s\n", info.Commit)
		fmt.Printf("built: %s\n", info.BuildDate)
		fmt.Printf("go: %s\n", info.GoVersion)
		fmt.Printf("modified: %s\n", info.Modified)
		return
	}

	var err error

	logger, stderr, logFile, err = logging.New()
	if err != nil {
		fmt.Fprintf(stderr, "error initializing logging: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	slog.SetDefault(logger)

	logger.Info("starting o7k", "version", info.Version, "commit", info.Commit, "built", info.BuildDate)

	registry := resource.NewRegistry()

	cloudsPaths, err := openstack.DiscoverCloudsFiles("")
	if err != nil {
		fmt.Fprintf(stderr, "error discovering clouds.yaml: %v\n", err)
		os.Exit(1)
	}

	_, err = openstack.LoadClouds(cloudsPaths)
	if err != nil {
		fmt.Fprintf(stderr, "error loading clouds.yaml: %v\n", err)
		os.Exit(1)
	}

	if err := registry.Register(contexts.New(cloudsPaths)); err != nil {
		fmt.Fprintf(stderr, "error registering contexts resource: %v\n", err)
		os.Exit(1)
	}

	openstackContext := openstack.Context{}
	err = registerAll(registry, &openstackContext)
	if err != nil {
		if joined, ok := err.(interface{ Unwrap() []error }); ok {
			for _, err := range joined.Unwrap() {
				fmt.Fprintf(stderr, "registration error: %v\n", err)
			}
		}
		os.Exit(1)
	}

	ver := info.GetVersion()
	update, err := version.CheckForUpdate(context.Background())
	if err != nil {
		fmt.Fprintf(stderr, "error checking for updates: %v\n", err)
	} else {
		if update.Available {
			if info.Version == "" {
				info.Version = version.Version
			}
			ver = fmt.Sprintf("%s → [update available: %s]", info.GetVersion(), update.LatestVersion)
		}
	}

	p := tea.NewProgram(
		ui.New(
			registry,
			&openstackContext,
			ver,
		),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(stderr, "error running o7k: %v\n", err)
		os.Exit(1)
	}
}

func registerAll(r *resource.Registry, openstackContext *openstack.Context) (errs error) {
	if err := r.Register(catalog.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering catalog resource: %w", err))
	}

	if err := r.Register(services.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering services resource: %w", err))
	}

	if err := r.Register(endpoints.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering endpoints resource: %w", err))
	}

	if err := r.Register(regions.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering regions resource: %w", err))
	}

	if err := r.Register(projects.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering projects resource: %w", err))
	}

	if err := r.Register(domains.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering domains resource: %w", err))
	}

	if err := r.Register(roles.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering roles resource: %w", err))
	}

	if err := r.Register(networks.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering networks resource: %w", err))
	}

	if err := r.Register(subnets.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering subnets resource: %w", err))
	}

	if err := r.Register(ports.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering ports resource: %w", err))
	}

	if err := r.Register(routers.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering routers resource: %w", err))
	}

	if err := r.Register(floatingips.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering floating IPs resource: %w", err))
	}

	if err := r.Register(securitygroups.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering security groups resource: %w", err))
	}

	if err := r.Register(securitygrouprules.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering security group rules resource: %w", err))
	}

	if err := r.Register(images.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering images resource: %w", err))
	}

	if err := r.Register(volumes.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering volumes resource: %w", err))
	}

	if err := r.Register(snapshots.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering snapshots resource: %w", err))
	}

	if err := r.Register(volumetypes.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering volume types resource: %w", err))
	}

	if err := r.Register(servers.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering servers resource: %w", err))
	}

	if err := r.Register(flavors.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering flavors resource: %w", err))
	}

	if err := r.Register(keypairs.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering keypairs resource: %w", err))
	}

	if err := r.Register(servergroups.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering server groups resource: %w", err))
	}

	if err := r.Register(availabilityzones.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering availability zones resource: %w", err))
	}

	if err := r.Register(backups.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering volume backups resource: %w", err))
	}

	if err := r.Register(stacks.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering stacks resource: %w", err))
	}

	if err := r.RegisterNavigationOnly(stackresources.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering stacks resources: %w", err))
	}

	if err := r.RegisterNavigationOnly(stackevents.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering stacks resources: %w", err))
	}

	if err := r.Register(zones.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering dns zones resource: %w", err))
	}

	if err := r.Register(recordsets.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering dns recordsets resource: %w", err))
	}

	if err := r.Register(loadbalancers.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering load balancers resource: %w", err))
	}

	if err := r.Register(listeners.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering listeners resource: %w", err))
	}

	if err := r.Register(pools.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering pools resource: %w", err))
	}

	if err := r.RegisterNavigationOnly(members.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering members resource: %w", err))
	}

	if err := r.Register(healthmonitors.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering health monitors resource: %w", err))
	}

	if err := r.Register(l7policies.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering L7 policies resource: %w", err))
	}

	if err := r.RegisterNavigationOnly(l7rules.New(openstackContext)); err != nil {
		errs = errors.Join(errs, fmt.Errorf("registering L7 rules resource: %w", err))
	}

	if err := r.Register(swiftcontainers.New(openstackContext)); err != nil {
		return fmt.Errorf("registering object containers resource: %w", err)
	}

	if err := r.RegisterNavigationOnly(swiftobjects.New(openstackContext)); err != nil {
		return fmt.Errorf("registering objects resource: %w", err)
	}

	if err := r.Register(secrets.New(openstackContext)); err != nil {
		return fmt.Errorf("registering secrets resource: %w", err)
	}

	if err := r.Register(containers.New(openstackContext)); err != nil {
		return fmt.Errorf("registering secret containers resource: %w", err)
	}

	if err := r.Register(orders.New(openstackContext)); err != nil {
		return fmt.Errorf("registering orders resource: %w", err)
	}

	return errs
}

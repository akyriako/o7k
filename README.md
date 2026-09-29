# o7k
![Status](https://img.shields.io/badge/status-early%20beta-orange)

**o7k** (pronounced **oooo-sju-k** 🇸🇪) is a terminal UI for OpenStack, **heavily** inspired by [k9s](https://github.com/derailed/k9s) ❣️ 

It provides a fast way to inspect and navigate OpenStack resources directly from the terminal.

<img width="2544" height="1247" alt="image" src="https://github.com/user-attachments/assets/fd654afe-5813-4aeb-99a6-0cd79fa56a3a" />

## Installation

### Homebrew

Install **o7k** using the Homebrew tap:

```bash
brew install --cask akyriako/tap/o7k
```

Upgrade to the latest release:

```bash
brew upgrade --cask o7k
```

### Debian / Ubuntu

Add the **o7k** repository signing key:

```bash
curl -fsSL https://akyriako.github.io/o7k-apt/o7k-archive-keyring.gpg \
  | sudo tee /usr/share/keyrings/o7k-archive-keyring.gpg >/dev/null
```

Add the APT repository:

```bash
echo "deb [signed-by=/usr/share/keyrings/o7k-archive-keyring.gpg] https://akyriako.github.io/o7k-apt/ stable main" \
  | sudo tee /etc/apt/sources.list.d/o7k.list
```

Install **o7k**:

```bash
sudo apt update
sudo apt install o7k
```

Upgrade to the latest release:

```bash
sudo apt update
sudo apt install --only-upgrade o7k
```

### Fedora / RHEL / Rocky Linux / AlmaLinux

Import the repository signing key:

```bash
sudo rpm --import https://akyriako.github.io/o7k-rpm/o7k-rpm-signing-key.asc
```

Add the **o7k** repository:

```bash
sudo tee /etc/yum.repos.d/o7k.repo >/dev/null <<'EOF'
[o7k]
name=o7k
baseurl=https://akyriako.github.io/o7k-rpm/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://akyriako.github.io/o7k-rpm/o7k-rpm-signing-key.asc
EOF
```

Install **o7k**:

```bash
sudo dnf install o7k
```

Upgrade to the latest release:

```bash
sudo dnf upgrade --refresh o7k
```

### Verify the installation

```bash
o7k --version
```

## Usage

**o7k** uses the standard OpenStack `clouds.yaml` configuration and automatically discovers the 
configuration files from the following places:

* `/etc/openstack`
* `~/.config/openstack`
* the current working directory

> [!Note]
> All the discovered files will be automatically loaded as context in **o7k** 

### Global Controls

| Key | Action |
|---|---|
| `:` | Switch resource |
| `Enter` | Execute the default action for the selected resource |
| `r` | Refresh the current resource |
| `c` | Copy the current detail-view as JSON or table-view as tab-delimited text |
| `Esc` | Dismiss / return to the previous view |
| `Ctrl+X` | Quit |

> [!TIP] 
> **o7k** will automatically refresh the current resource every 30s, you don't need to explicitely press `r`.

To switch between resources, press `:` and enter a resource name or alias (**exactly** as you are being used to in k9s)

e.g.:

```text
:servers
:networks
:volumes
:images
```

Resource-specific commands are displayed in the _header_ while a resource is active.

## Resource Navigation

**o7k** supports navigation between related OpenStack resources.

Examples include:

```text
Server              -> Image

Network             -> Subnets
Network             -> Ports
Subnet              -> Network
Port                -> Network

Router              -> Ports			-> Networks			-> Subnets
Router              -> Floating IPs
Floating IP         -> Router
Floating IP         -> Port
Port                -> Floating IPs

Volume              -> Snapshots
Volume              -> Backups
Snapshot            -> Volume
Backup              -> Volume


Security Group      -> Security Group Rules
Security Group Rule -> Security Group
```

When navigating to a related collection, **o7k** filters the destination resource automatically.

<img width="2544" height="1247" alt="image" src="https://github.com/user-attachments/assets/ca91e438-06a9-4120-8753-c3b0dcd5ab4c" />

Press `Esc` to return to the previous resource and selection.

## Currently Supported Resources

| Service | Resources | Status |
|---|---|:-:|
| **Local / Auth** | contexts | Partial |
| **Keystone / Identity** | projects, users, groups, roles, domains | Partial |
| **Keystone / Catalog** | catalog, services, endpoints, regions | ✅ |
| **Cinder / Block Storage** | volumes, snapshots, volume types, volume backups | ✅ |
| **Glance / Image** | images | ✅ |
| **Swift / Object Storage** | swift containers, swift objects | ✅ |
| **Nova / Compute** | servers, flavors | ✅ |
| **Nova / Compute** | keypairs, server groups, availability zones | ✅ |
| **Neutron / Network** | networks, subnets, ports, routers, floating IPs, security groups | ✅ |
| **Neutron / Network** | security group rules | ✅ |
| **Neutron / Network** | router interfaces | ⬜ |
| **Designate / DNS** | zones, recordsets | ✅ |
| **Heat / Orchestration** | stacks, stack resources, stack events | ✅ |
| **Octavia / Load Balancing** | load balancers, listeners, pools, members | ✅ |
| **Octavia / Load Balancing** | health monitors, L7 policies, L7 rules | ✅ |
| **Barbican / Key Manager** | secrets, secret containers, orders | ✅ |
| **Trove / Databases** | instances | ⬜ |

> [!CAUTION]
> **o7k is early beta and under active development.** Expect bugs, incomplete features, and breaking changes. Behavior and configuration may change without notice.
>
> Use against production OpenStack environments at your own risk.

### Aliases

| Service | Resource | Command | Aliases |
|---|---|---|---|
| Command Palette | Aliases | `commands` | `command`, `cmd`, `aliases`, `alias` |
| Local / Auth | Contexts | `contexts` | `context`, `ctx`, `cloud`, `clouds` |
| Keystone / Catalog | Catalog | `catalog` | `cat` |
| Keystone / Catalog | Services | `services` | `service`, `svc` |
| Keystone / Catalog | Endpoints | `endpoints` | `endpoint` |
| Keystone / Catalog | Regions | `regions` | `region` |
| Keystone / Identity | Projects | `projects` | `project` |
| Keystone / Identity | Users | `users` | `user` |
| Keystone / Identity | Groups | `groups` | `group` |
| Keystone / Identity | Roles | `roles` | `role` |
| Keystone / Identity | Domains | `domains` | `domain` |
| Nova / Compute | Servers | `servers` | `server`, `srv` |
| Nova / Compute | Flavors | `flavors` | `flavor` |
| Nova / Compute | Keypairs | `keypairs` | `keypair`, `keys`, `key` |
| Nova / Compute | Server Groups | `servergroups` | `servergroup`, `sgroups`, `sgroup` |
| Nova / Compute | Availability Zones | `availabilityzones` | `availabilityzone`, `azs`, `az` |
| Neutron / Network | Networks | `networks` | `network`, `net` |
| Neutron / Network | Subnets | `subnets` | `subnet` |
| Neutron / Network | Ports | `ports` | `port` |
| Neutron / Network | Routers | `routers` | `router` |
| Neutron / Network | Floating IPs | `floatingips` | `floatingip`, `fips`, `fip` |
| Neutron / Network | Security Groups | `securitygroups` | `securitygroup`, `secgroups`, `secgroup`, `sg` |
| Neutron / Network | Security Group Rules | `securitygrouprules` | `securitygrouprule`, `secrules`, `secrule`, `sgr` |
| Cinder / Block Storage | Volumes | `volumes` | `volume`, `vol` |
| Cinder / Block Storage | Snapshots | `snapshots` | `snapshot`, `snap` |
| Cinder / Block Storage | Volume Types | `volumetypes` | `volumetype` |
| Cinder / Block Storage | Volume Backups | `backups` | `backup`, `volume-backups`, `volume-backup` |
| Glance / Image | Images | `images` | `image`, `img` |
| Swift / Object Storage | Swift Containers | `swift-containers` | `swift-container` |
| Swift / Object Storage | Swift Objects | `swift-objects` | `swift-object` |
| Heat / Orchestration | Stacks | `stacks` | `stack` |
| Heat / Orchestration | Stack Resources | `stack-resources` | `stack-resource` |
| Heat / Orchestration | Stack Events | `stack-events` | `stack-event` |
| Octavia / Load Balancing | Load Balancers | `loadbalancers` | `loadbalancer`, `lbs`, `lb` |
| Octavia / Load Balancing | Listeners | `listeners` | `listener` |
| Octavia / Load Balancing | Pools | `pools` | `pool` |
| Octavia / Load Balancing | Members | `members` | `member` |
| Octavia / Load Balancing | Health Monitors | `healthmonitors` | `healthmonitor`, `monitors`, `monitor` |
| Octavia / Load Balancing | L7 Policies | `l7policies` | `l7policy`, `l7-policies`, `l7-policy` |
| Octavia / Load Balancing | L7 Rules | `l7rules` | `l7rule`, `l7-rules`, `l7-rule` |
| Barbican / Key Manager | Secrets | `secrets` | `secret` |
| Barbican / Key Manager | Secret Containers | `secret-containers` | `secret-container` |
| Barbican / Key Manager | Orders | `orders` | `order` |
| Designate / DNS | Zones | `zones` | `zone`, `dns-zones`, `dns-zone` |
| Designate / DNS | Recordsets | `recordsets` | `recordset`, `records`, `record` |

## Development

### Adding a new OpenStack Resource

**o7k** resources are organized by OpenStack service and implement a common resource interface. Adding support for a new resource normally involves four areas:

1. Add or reuse the OpenStack service client in `openstack.Context`.
2. Implement the resource under `internal/resources/<service>/`.
3. Register the resource in `cmd/main.go`.
4. Update the README support matrix.

The examples below use a fictional `widgets` resource belonging to an OpenStack service called `example`.

#### 1. Add the OpenStack service client

OpenStack API access is centralized through `openstack.Context`. Resources should use a service client provided by the context rather than creating their own authenticated provider.

Add a helper for the service **if one does not already exist**.

For example:

```go
func (c *Context) ExampleV2() (*gophercloud.ServiceClient, error) {
	return example.NewExampleV2(
		c.Provider,
		gophercloud.EndpointOpts{
			Region: c.Region,
		},
	)
}
```

Use the appropriate Gophercloud constructor and API version for the OpenStack service.

> [!Important]
> If the service client already exists in `openstack.Context`, reuse it. **Do not create another client specifically for the new resource**.

A resource can then obtain its client with:

```go
client, err := r.context.ExampleV2()
if err != nil {
	return nil, fmt.Errorf("creating example client: %w", err)
}
```

#### 2. Create the resource package

Resources live below the OpenStack service they belong to.

For example:

```text
internal/resources/
├── compute/
│   ├── servers/
│   ├── flavors/
│   ├── keypairs/
│   └── servergroups/
├── networking/
│   ├── networks/
│   ├── subnets/
│   ├── ports/
│   └── routers/
├── blockstorage/
│   ├── volumes/
│   ├── snapshots/
│   ├── volumetypes/
│   └── backups/
└── example/
    └── widgets/
        ├── widgets.go
        └── commands.go
```

> [!Important]
> Keep the basic resource definition and listing logic in the resource file. Put commands and command execution in `commands.go` when the resource has commands.

#### 3. Implement the resource interface

Every resource implements the common **o7k** resource contract.

A typical resource looks like:

```go
package widgets

import (
	"context"

	"github.com/akyriako/o7k/internal/openstack"
	"github.com/akyriako/o7k/internal/resource"
)

type Resource struct {
	context *openstack.Context
}

func New(openstackContext *openstack.Context) *Resource {
	return &Resource{
		context: openstackContext,
	}
}

func (r *Resource) Kind() string {
	return "widgets"
}

func (r *Resource) Title() string {
	return "Widgets"
}

func (r *Resource) Aliases() []string {
	return []string{
		"widget",
		"wid",
	}
}
```

- `Kind()` is the canonical resource name used by **o7k**.  
- `Aliases()` provides alternative names accepted by the resource command:

```text
:widgets
:widget
:wid
```

> [!Important]
> Choose short aliases that do not conflict with existing resources.

#### 4. Define the table columns

Each resource declares the columns displayed in the table.

For example:

```go
func (r *Resource) Columns() []resource.Column {
	return []resource.Column{
		{
			Key:      "id",
			Title:    "ID",
			MinWidth: 40,
			Flex:     0,
		},
		{
			Key:      "name",
			Title:    "NAME",
			MinWidth: 20,
			Flex:     1,
		},
		{
			Key:      "status",
			Title:    "STATUS",
			MinWidth: 12,
			Flex:     0,
		},
	}
}
```
> [!Important]
> UUID/GUID columns are advised to use a fixed width of `40` and `Flex` of `0`:

```go
{
	Key:      "id",
	Title:    "ID",
	MinWidth: 40,
	Flex:     0,
}
```

- `MinWidth` is the minimum width of the column.
- `Flex` controls how additional terminal width is distributed. Use `Flex: 0` for columns that should remain fixed and a positive value for columns that may expand.

> [!Tip]
> Do not shrink UUID columns to make a resource fit into a small terminal. **o7k** supports horizontal table scrolling for tables wider than the available terminal.

#### 5. Implement `List()`

`List()` retrieves the OpenStack objects and converts them into generic `resource.Row` values understood by the UI.

Start by obtaining the appropriate service client from `openstack.Context`:

```go
func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	client, err := r.context.ExampleV2()
	if err != nil {
		return nil, err
	}

	// Use the appropriate Gophercloud List/Extract API here.

	return rows, nil
}
```

A typical conversion from an OpenStack object to an **o7k** row looks like:

```go
rows := make([]resource.Row, 0, len(items))

for _, item := range items {
	rows = append(rows, resource.Row{
		ID: item.ID,
		Fields: map[string]string{
			"id":     item.ID,
			"name":   item.Name,
			"status": item.Status,
		},
	})
}

return rows, nil
```
> [!Important]
> Every `Fields` key used by `Columns()` must be populated by `List()`.  
> For example, this column:
> ```go
> {
>	Key:      "status",
>	Title:    "STATUS",
>	MinWidth: 12,
>	Flex:     0,
> }
> ```
> 
> expects:
> 
> ```go
> Fields: map[string]string{
>	"status": item.Status,
> }
> ```

> [!Warning]
> - The `ID` field on `resource.Row` should contain the canonical identifier used for commands, navigation and cursor restoration.
> - Avoid performing an additional API request for every row just to populate a table column. If an API's list operation does not provide a value, consider leaving that information for `Show` rather than introducing a complex multi-service request pattern.

#### 6. Add resource commands

Commands are **optional**. Only expose commands that make sense for the resource.

For example:

```go
func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{
			Key:         "enter",
			Description: "Show",
			Default:     true,
		},
	}
}
```

Do not automatically add `Show` to every resource. `Show` support in **o7k** is opt-in.

Resources can also expose navigation commands to related resources. For example:

```go
func (r *Resource) Commands() []resource.Command {
	return []resource.Command{
		{
			Key:         "enter",
			Description: "Show",
			Default:     true,
		},
		{
			Key:         "shift-n",
			Description: "Networks",
		},
	}
}
```

>[!Caution]
> Keep global **o7k** shortcuts in mind when selecting keys.
>
> The following keys are reserved globally:
> 
> ```text
> ctrl+x   Quit
> :        Resource command
> r        Refresh
> c        Copy
> esc      Dismiss / navigate back
> enter    Default resource command
> ```

#### 7. Implement command execution

Resource-specific command handling belongs in the resource package.

A `Show` command normally retrieves the complete object and returns a generic `DetailsMsg` struct:

```go
return func() tea.Msg {
	item, err := widgets.Get(
		context.Background(),
		client,
		row.ID,
	).Extract()

	if err != nil {
		return resource.DetailsMsg{
			ID:  row.ID,
			Err: err,
		}
	}

	return resource.DetailsMsg{
		ID:      row.ID,
		Content: item,
	}
}
```

The UI handles `resource.DetailsMsg` generically and renders its JSON content in the details view.

#### 8. Add navigation to related resources

Relationships between OpenStack resources should use **o7k**'s generic navigation messages rather than implementing navigation directly in the UI.

To navigate directly to another resource:

```go
return func() tea.Msg {
	return resource.NavigateMsg{
		Resource: "images",
		ID:       imageID,
	}
}
```

Use this when the target object has a known ID and the destination table should select that object.

For one-to-many relationships, use filtered navigation:

```go
return func() tea.Msg {
	return resource.NavigateFilteredMsg{
		Resource: "snapshots",
		Field:    "volume_id",
		Value:    row.ID,
	}
}
```

This opens the target resource while filtering rows by the specified field.

Filtered navigation currently performs exact scalar matching:

```text
row.Fields[field] == value
```

The target resource therefore needs to expose the relationship in its row fields.

For example:

```go
Fields: map[string]string{
	"id":        snapshot.ID,
	"name":      snapshot.Name,
	"volume_id": snapshot.VolumeID,
}
```

>[!Tip]
> - Do not add special-case navigation logic to the UI for a particular OpenStack resource.
> - Prefer `NavigateFilteredMsg` over `NavigateMsg`, as it's visually more intuitive for the users.

##### Scoped navigation

Some OpenStack APIs cannot list a child resource without information about its parent. In these cases, use scoped navigation instead of client-side filtered navigation.

For example, listing members of an Octavia pool requires the pool ID:

```go
return func() tea.Msg {
	return resource.NavigateScopedMsg{
		Resource: "members",
		Scope: map[string]string{
			"pool_id": row.ID,
		},
	}
}
```

The destination resource retrieves the scope from the context in `List()`:

```go
func (r *Resource) List(ctx context.Context) ([]resource.Row, error) {
	scope := resource.Scope(ctx)
	poolID := scope["pool_id"]

	if poolID == "" {
		return nil, fmt.Errorf("member requires pool_id")
	}
}
```

Scope can contain multiple values when required by the API. Heat stack resources, for example, require both the stack name and stack ID:

```go
Scope: map[string]string{
	"stack_name": row.Fields["name"],
	"stack_id":   row.ID,
}
```

Scoped navigation is also appropriate when the OpenStack API supports an optional server-side filter. In that case the resource may still support normal unscoped listing:

```go
scope := resource.Scope(ctx)

opts := listeners.ListOpts{
	LoadbalancerID: scope["loadbalancer_id"],
}
```

With no scope, the resource lists globally. When reached through scoped navigation, the relationship is filtered by the OpenStack API.

> [!Important]
> Scope belongs to the navigation state. The UI preserves it across refreshes and restores the previous scope when navigating back with `Esc`.
>
> Do not encode list-valued relationships into scalar fields merely to make `NavigateFilteredMsg` work. Use the navigation mechanism that matches the OpenStack API and relationship.

#### 9. Register the resource

Creating the package is not enough. The resource must be registered so **o7k** knows it exists.

Add the resource package to the imports in `cmd/main.go`.

For example:

```go
import (
	"github.com/akyriako/o7k/internal/resources/example/widgets"
)
```

Then register it with the resource registry alongside the existing resources:

```go
registry.Register(
	widgets.New(openstackContext),
)
```

Registration makes the canonical resource name and its aliases available through the **o7k** resource command.

For example:

```text
:widgets
:widget
:wid
```

Keep resource registration in a logical service/resource order rather than adding new resources at arbitrary positions.

##### Navigation-only resources

Not every resource can be listed independently. Some OpenStack APIs require parent information before a child collection 
can be queried. Examples include Octavia members, which require a pool ID, and L7 rules, which require an L7 policy ID. 
These resources must still be registered so generic navigation can resolve them, but they should not be exposed as 
standalone `:<resource>` commands.

Register them with:

```go
registry.RegisterNavigationOnly(
	members.New(openstackContext),
)
```

A navigation-only resource:

- is registered in the resource registry
- can be resolved by generic navigation
- can have canonical names and aliases
- can receive scope through NavigateScopedMsg
- **cannot be opened** directly through `:<resource>`

For instance in Octavia, the following is valid because the selected pool supplies `pool_id`:

```text
:pools -> members
```

but accessing `members` directly is not valid because there is no parent pool from which to obtain the required scope.

Use normal registration, `registry.Registry`,  when a resource can list independently:

```go
registry.Register(
	listeners.New(openstackContext),
)
```

A resource may support both global listing and optional scoped navigation. Octavia listeners are an example: `:listeners` 
can list all listeners, while Load Balancer → Listeners can pass `loadbalancer_id` for server-side filtering. Such resources 
use normal `Register` and **not** `RegisterNavigationOnly`.

> [!Warning]
> Use `RegisterNavigationOnly` only when the resource fundamentally requires parent context to perform its list operation. 
**Do not make a resource navigation-only merely because it participates in a parent/child relationship.**

#### 10. Update the support matrix

Finally, update the supported-resource table in this README.

For example:

```markdown
| **Example Service** | widgets | ✅ |
```

Use the existing status convention:

```text
✅       Supported
Partial  Partially supported
⬜       Planned / not implemented
```

>[!Warning]
> A resource should only be marked as supported once its basic listing and intended commands work against an actual OpenStack environment.

### Add a Plugin for an OpenStack-based Provider

**o7k** can be extended with plugins that provide resources specific to an OpenStack-based provider without adding provider-specific code to the **o7k** codebase.

A **plugin** is a standalone executable that communicates with **o7k** through the public `pluginsdk` module. A single plugin can expose multiple resources across one or more OpenStack services.

Plugins run as separate processes and therefore do not share the authenticated OpenStack client used internally by **o7k**. Instead, a plugin receives information about the currently active **o7k** context and authenticates independently using the same `clouds.yaml` configuration.

The examples below use a fictional provider called `example` and a plugin named `o7k-plugin-example`.

#### 1. Create the plugin module

Create a separate Go module/project for the plugin and follow the same structure as the example plugin: e.g.:

```text
o7k-plugin-example/
├── go.mod
├── main.go
└── internal/
    └── plugin/
        ├── client.go
        ├── plugin.go
        └── resources/
            ├── compute/
            │   ├── commands.go
            │   └── servers.go
            └── blockstorage/
                ├── commands.go
                └── volumes.go
```


Initialize the module and add the **o7k** plugin SDK:

```bash
go mod init example.com/o7k-plugin-example
go get github.com/akyriako/o7k/pluginsdk
```
> [!Tip]
> Add the **OpenStack SDK used by the provider** as a dependency as well.
>
> The plugin is an independent Go application and must not import packages from `github.com/akyriako/o7k/internal/xxx`.  
> Communication between the plugin and o7k happens over gRPC, using [HashiCorp's go-plugin](https://github.com/hashicorp/go-plugin),
> to manage the plugin subprocess and RPC connection. The public `github.com/akyriako/o7k/pluginsdk` module defines the 
> interfaces and gRPC protocol shared by o7k and the plugin(s).
>
> The examples in `/examples/plugins/demo` uses two T Cloud Public (formerly known as Open Telekom Cloud) resources:
>
>```text
>Compute       -> ecs-servers
>Block Storage -> ecs-volumes
>```

#### 2. Implement the plugin

Create `internal/plugin/plugin.go`.

The plugin owns the host connection, the authenticated provider client and the resources it exposes:

```go
package plugin

import (
	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

type Plugin struct {
	host      pluginsdk.Host
	provider  *pluginsdk.ClientProvider[*golangsdk.ProviderClient]
	resources []pluginsdk.Resource
}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) SetHost(host pluginsdk.Host) {
	p.host = host
	p.provider = newClient(host)
}

func (p *Plugin) Register(resources ...pluginsdk.Resource) {
	p.resources = append(p.resources, resources...)
}

func (p *Plugin) Metadata() pluginsdk.Metadata {
	return pluginsdk.Metadata{
		Name:    "example",
		Version: "0.1.0",
		Color:   "#E20074",
	}
}

func (p *Plugin) Resources() []pluginsdk.Resource {
	return p.resources
}

func (p *Plugin) Host() pluginsdk.Host {
	return p.host
}

func (p *Plugin) Provider() *pluginsdk.ClientProvider[*golangsdk.ProviderClient] {
	return p.provider
}
```

`Metadata()` identifies the plugin to **o7k**:

- `Name` is the provider/plugin name displayed by **o7k**.
- `Version` is the plugin version.
- `Color` is used by **o7k** when rendering resources belonging to the plugin.

`Register()` collects the resources exposed by the plugin. A single provider plugin can register resources from multiple OpenStack services.

`SetHost()` is called by **o7k** when the plugin process is initialized. It stores the host connection and creates the provider client manager.

> [!Important]
> 1. Fill in the `Metadata()` and leave the rest intact.
> 2. Create the `pluginsdk.ClientProvider[T]`, where T corresponds to the target cloud provider.

#### 3. Create the provider client

Plugins run in a separate process from **o7k** and therefore cannot reuse the authenticated provider client owned by **o7k**. 
Instead, the plugin authenticates independently using the active context supplied by the host.

Create `internal/plugin/client.go`:

```go
package plugin

import (
	"context"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
)

func newClient(host pluginsdk.Host) *pluginsdk.ClientProvider[*golangsdk.ProviderClient] {
	return pluginsdk.NewClientProvider(host, connectClient)
}

func connectClient(_ context.Context, current pluginsdk.Context) (*golangsdk.ProviderClient, error) {
	data, err := os.ReadFile(current.CloudsPath)
	if err != nil {
		return nil, fmt.Errorf("reading clouds.yaml %q: %w", current.CloudsPath, err)
	}

	var config openstack.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing clouds.yaml %q: %w", current.CloudsPath, err)
	}

	cloud, ok := config.Clouds[current.Cloud]
	if !ok {
		return nil, fmt.Errorf("cloud %q not found in %q", current.Cloud, current.CloudsPath)
	}

	if current.Region != "" {
		cloud.RegionName = current.Region
	}

	provider, err := openstack.AuthenticatedClientFromCloud(&cloud)
	if err != nil {
		return nil, fmt.Errorf("authenticating cloud %q: %w", current.Cloud, err)
	}

	return provider, nil
}
```

> [!Important]
> The exact authentication implementation depends on the OpenStack Golang SDK used/developed by the target cloud provider. The example above 
uses Open Telekom Cloud Golang SDK. A provider plugin using another SDK should implement `connectClient` using that SDK's `clouds.yaml` and authentication support.

> [!Note]
> Always use both `Cloud` and `CloudsPath`. `CloudsPath` identifies the exact `clouds.yaml` file from which **o7k** loaded the active cloud.
>
> `pluginsdk.ClientProvider` caches the authenticated provider client for the current context generation. When the user activates another **o7k** context, the generation changes and the provider client is recreated automatically.
> 
> Plugin resources should, as we will see later, obtain the provider through:
>
>```go
>provider, err := r.plugin.Provider().Client(ctx)
>if err != nil {
>	return nil, fmt.Errorf("getting provider: %w", err)
>}
>```
>
>rather than maintaining their own authenticated provider clients.

#### 4. Add OpenStack service clients

Resources should obtain service clients from the plugin rather than creating a new service client every time `List()` or a command is executed. 
Add a helper for each OpenStack service used by the plugin in `internal/plugin/plugin.go`:

```go
func (p *Plugin) ComputeV2(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "compute", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewComputeV2(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}

func (p *Plugin) BlockStorageV3(ctx context.Context) (*golangsdk.ServiceClient, error) {
	return pluginsdk.GetServiceClient(ctx, p.provider, "block-storage", func(provider *golangsdk.ProviderClient, current pluginsdk.Context) (*golangsdk.ServiceClient, error) {
		return openstack.NewBlockStorageV3(provider, golangsdk.EndpointOpts{Region: current.Region})
	})
}
```

> [!Note]
> Use the service constructor provided by the OpenStack SDK used by the provider.
> A resource can then simply request the client it needs:
>
>```go
>client, err := r.plugin.ComputeV2(ctx)
>if err != nil {
>	return nil, fmt.Errorf("getting compute client: %w", err)
>}
>```
>
>`GetServiceClient` keeps service clients associated with the active **o7k** context generation, just like `ClientProvider` does for the provider client.

#### 5. Create a resource package

Resources are better organized below `internal/plugin/resources/` by service.

For example:

```text
internal/plugin/resources/
├── compute/
│   ├── commands.go
│   └── servers.go
└── blockstorage/
    ├── commands.go
    └── volumes.go
```

Keep the resource definition and listing logic in the resource file. Put command implementations in `commands.go`.

For example, create `internal/plugin/resources/compute/servers.go`:

```go
package compute

import (
	"context"
	"fmt"

	"example.com/o7k-plugin-example/internal/plugin"
)

type Resource struct {
	plugin *plugin.Plugin
}

func NewExampleServers(p *plugin.Plugin) *Resource {
	return &Resource{
		plugin: p,
	}
}

func (r *Resource) Service() string {
	return "example-compute"
}

func (r *Resource) Kind() string {
	return "example-servers"
}

func (r *Resource) Title() string {
	return "Example Servers"
}

func (r *Resource) Aliases() []string {
	return nil
}
```

`Service()` identifies the provider service the resource belongs to.

`Kind()` is the canonical resource name used by **o7k**:

```text
:example-servers
```

`Title()` is the human-readable resource name shown in the UI.

`Aliases()` can provide additional names accepted by the resource command:

```go
func (r *Resource) Aliases() []string {
	return []string{
		"example-server",
	}
}
```

#### 6. Define the table columns

Plugin resources use `pluginsdk.Column` to describe their table view:

```go
func (r *Resource) Columns() []pluginsdk.Column {
	return []pluginsdk.Column{
		{Key: "id", Title: "ID", MinWidth: 36},
		{Key: "name", Title: "NAME", MinWidth: 24, Flex: 1},
		{Key: "status", Title: "STATUS", MinWidth: 12},
	}
}
```

Each column key maps to a field returned by `List()`.

`MinWidth` specifies the minimum width of the column. `Flex` controls how additional terminal width is distributed.

Fields do not have to be visible columns. A resource may add additional values to `Row.Fields` for navigation or relationships between resources.

#### 7. Implement `List()`

`List()` obtains the service client from the plugin, retrieves the provider resources and converts them to `pluginsdk.Row` values:

```go
func (r *Resource) List(ctx context.Context) ([]pluginsdk.Row, error) {
	client, err := r.plugin.ComputeV2(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting compute client: %w", err)
	}

	pages, err := servers.List(client, servers.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("listing servers: %w", err)
	}

	allServers, err := servers.ExtractServers(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting servers: %w", err)
	}

	rows := make([]pluginsdk.Row, 0, len(allServers))

	for _, server := range allServers {
		rows = append(rows, pluginsdk.Row{
			ID: server.ID,
			Fields: map[string]string{
				"id":     server.ID,
				"name":   server.Name,
				"status": server.Status,
			},
		})
	}

	return rows, nil
}
```

`Row.ID` must contain the canonical identifier used by commands and navigation.

Every field referenced by `Columns()` must be present in `Fields`.

Additional fields may be included without defining a column. For example, a volume resource can expose the attached server ID for navigation without displaying it:

```go
Fields: map[string]string{
	"id":        volume.ID,
	"name":      volume.Name,
	"status":    volume.Status,
	"size":      strconv.Itoa(volume.Size),
	"server_id": serverID,
}
```

#### 8. Add resource commands

Commands are declared using `pluginsdk.Command`:

```go
func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
	}
}
```

`Default: true` makes the command the default action for the resource, which is executed by pressing `Enter`.

Command execution is handled by `Execute()`:

```go
func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row.ID)
	}

	return pluginsdk.Result{}, nil
}
```

Keep the implementation of individual commands in `commands.go`.

#### 9. Return resource details

A `Show` command retrieves the complete provider object, serializes it as JSON and returns `pluginsdk.Details`.

For example, in `internal/plugin/resources/compute/commands.go`:

```go
package compute

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/akyriako/o7k/pluginsdk"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
)

func (r *Resource) show(ctx context.Context, id string) (pluginsdk.Result, error) {
	client, err := r.plugin.ComputeV2(ctx)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting compute client: %w", err)
	}

	server, err := servers.Get(ctx, client, id).Extract()
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("getting server %q: %w", id, err)
	}

	content, err := json.Marshal(server)
	if err != nil {
		return pluginsdk.Result{}, fmt.Errorf("encoding server %q: %w", id, err)
	}

	return pluginsdk.Result{
		Details: &pluginsdk.Details{
			ID:      server.ID,
			Content: content,
		},
	}, nil
}
```

The JSON is transported back to **o7k** and rendered by the normal details view. The plugin does not implement any UI-related functionality.

#### 10. Navigate between plugin resources

Plugin resources can use `pluginsdk.Navigate` to participate in normal **o7k** resource navigation.

For example, a compute server can expose a command that opens only the volumes attached to that server:

```go
func (r *Resource) Commands() []pluginsdk.Command {
	return []pluginsdk.Command{
		{Key: "s", Description: "Show", Default: true},
		{Key: "shift-v", Description: "Volumes"},
	}
}
```

Handle the command in `Execute()`:

```go
func (r *Resource) Execute(ctx context.Context, command pluginsdk.Command, row pluginsdk.Row) (pluginsdk.Result, error) {
	switch command.Key {
	case "s":
		return r.show(ctx, row.ID)
	case "shift-v":
		return r.volumes(row)
	}

	return pluginsdk.Result{}, nil
}
```

The navigation command can then return:

```go
func (r *Resource) volumes(row pluginsdk.Row) (pluginsdk.Result, error) {
	return pluginsdk.Result{
		Navigate: &pluginsdk.Navigate{
			Resource: "example-volumes",
			Field:    "server_id",
			Value:    row.ID,
		},
	}, nil
}
```

The destination resource exposes `server_id` in its row fields:

```go
Fields: map[string]string{
	"id":        volume.ID,
	"name":      volume.Name,
	"status":    volume.Status,
	"server_id": serverID,
}
```

**o7k** then opens `example-volumes` and filters the rows using:

```text
row.Fields["server_id"] == selectedServer.ID
```

`pluginsdk.Navigate` also supports direct navigation by ID:

```go
Navigate: &pluginsdk.Navigate{
	Resource: "example-servers",
	ID:       serverID,
}
```

and scoped navigation for APIs that require parent information:

```go
Navigate: &pluginsdk.Navigate{
	Resource: "example-resource",
	Scope: map[string]string{
		"parent_id": row.ID,
	},
}
```

Use the navigation form that matches the relationship exposed by the provider API.

#### 11. Register the resources

Creating a resource does not automatically expose it to **o7k**.

Register the resources when constructing the plugin in `main.go`:

```go
package main

import (
	"example.com/o7k-plugin-example/internal/plugin"
	"example.com/o7k-plugin-example/internal/plugin/resources/blockstorage"
	"example.com/o7k-plugin-example/internal/plugin/resources/compute"
	"github.com/akyriako/o7k/pluginsdk"
)

func main() {
	p := plugin.New()

	p.Register(
		compute.NewServers(p),
		blockstorage.NewVolumes(p),
	)

	pluginsdk.Serve(p)
}
```

A single call to `Register()` may contain resources from multiple OpenStack services.

`pluginsdk.Serve()` starts the plugin process protocol used by **o7k**. No additional RPC or transport code is required 
in the provider plugin.

#### 12. Build and install the plugin locally

Build the plugin as a normal Go executable:

```bash
go build -o o7k-plugin-example .
```

Install the locally built executable:

```bash
o7k plugin install ./o7k-plugin-example
```

Start **o7k** and open the built-in plugin resource:

```text
:plugins
```

The plugin should be listed with its name, version and loading status.

Its registered resources are then available through the normal resource command:

```text
:example-servers
:example-volumes
```

If the plugin fails to initialize, inspect its error through `:plugins`.

When rebuilding the plugin during development, install the new executable again and **restart** **o7k**.

#### 13. Distribute the plugin

Provider plugins are distributed independently of **o7k**.

The provider project is responsible for building and publishing the plugin executable for the platforms it supports. 
A release asset, public object-storage URL or any other directly downloadable HTTP(S) location can be used.

Users install the published executable directly:

```bash
o7k plugin install https://example.com/releases/o7k-plugin-example
```

The downloaded executable is copied into the **o7k** plugin directory and is loaded the next time **o7k** starts.

To remove an installed plugin:

```bash
o7k plugin remove o7k-plugin-example
```

**Restart** **o7k** after removing the plugin.

> [!Caution]
> Plugins are executable programs running on the user's machine. Only install plugins from sources you trust.

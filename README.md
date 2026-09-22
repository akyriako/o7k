# o7k
![Status](https://img.shields.io/badge/status-early%20beta-orange)

**o7k** (pronounced **oooo-sju-k** 🇸🇪) is a terminal UI for OpenStack, **heavily** inspired by [k9s](https://github.com/derailed/k9s) ❣️ 

It provides a fast way to inspect and navigate OpenStack resources directly from the terminal.

<img width="2544" height="1247" alt="Screenshot from 2026-09-18 13-00-06" src="https://github.com/user-attachments/assets/55507289-02e6-4f91-8e5f-92ed6b14ab09" />  

> [!CAUTION]
> **o7k is early beta and under active development.** Expect bugs, incomplete features, and breaking changes. Behavior and configuration may change without notice.
>
> Use against production OpenStack environments at your own risk.

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
| `c` | Copy the current detail view |
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

Router              -> Ports
Router              -> Floating IPs
Floating IP         -> Router
Floating IP         -> Port
Port                -> Floating IPs

Security Group      -> Security Group Rules
Security Group Rule -> Security Group

Volume              -> Snapshots
Volume              -> Backups
Snapshot            -> Volume
Backup              -> Volume
```

When navigating to a related collection, **o7k** filters the destination resource automatically.

<img width="2544" height="1247" alt="Screenshot from 2026-09-18 13-01-45" src="https://github.com/user-attachments/assets/10652fed-f76e-4bd7-bdee-bafe20fb39a7" />

Press `Esc` to return to the previous resource and selection.

## Currently Supported Resources

| Service | Resources | Status |
|---|---|:-:|
| **Local / Auth** | contexts | Partial |
| **Keystone / Identity** | projects, users, groups, roles, domains | Partial |
| **Keystone / Catalog** | catalog, services, endpoints, regions | ✅ |
| **Cinder / Block Storage** | volumes, snapshots, volume types, volume backups | ✅ |
| **Glance / Image** | images | ✅ |
| **Nova / Compute** | servers, flavors | ✅ |
| **Nova / Compute** | keypairs, server groups, availability zones | ✅ |
| **Neutron / Network** | networks, subnets, ports, routers, floating IPs, security groups | ✅ |
| **Neutron / Network** | security group rules | ✅ |
| **Neutron / Network** | router interfaces | ⬜ |
| **Designate / DNS** | zones, recordsets | ✅ |
| **Heat / Orchestration** | stacks | Partial |
| **Trove / Databases** | instances | ⬜ |

### Aliases

| Service | Resource | Command | Aliases |
|---|---|---|---|
| Local / Auth | Contexts | `contexts` | `context`, `ctx`, `cloud`, `clouds` |
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
| Cinder / Block Storage | Volume Types | `volumetypes` | `volumetype`, `types`, `type` |
| Cinder / Block Storage | Volume Backups | `backups` | `backup`, `volume-backups`, `volume-backup` |
| Glance / Image | Images | `images` | `image`, `img` |
| Keystone / Identity | Projects | `projects` | `project` |
| Keystone / Identity | Users | `users` | `user` |
| Keystone / Identity | Groups | `groups` | `group` |
| Keystone / Identity | Roles | `roles` | `role` |
| Keystone / Identity | Domains | `domains` | `domain` |
| Keystone / Catalog | Services | `services` | `service`, `svc` |
| Keystone / Catalog | Endpoints | `endpoints` | `endpoint` |
| Keystone / Catalog | Regions | `regions` | `region` |
| Heat / Orchestration | Stacks | `stacks` | `stack` |
| Designate / DNS | Zones | `zones` | `zone`, `dns-zones`, `dns-zone` |
| Designate / DNS | Recordsets | `recordsets` | `recordset`, `records`, `record` |

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
sudo apt upgrade o7k
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
sudo dnf upgrade -- refresh o7k
```

### Verify the installation

```bash
o7k --version
```

## Development

### Adding a New OpenStack Resource

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

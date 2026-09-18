# o7k

`o7k` (pronounced **oooo-sju-k** 🇸🇪) is a terminal UI for OpenStack, **heavily** inspired by [k9s](https://github.com/derailed/k9s) ❣️ 

It provides a fast way to inspect and navigate OpenStack resources directly from the terminal.

<img width="2544" height="1247" alt="Screenshot from 2026-09-18 13-00-06" src="https://github.com/user-attachments/assets/55507289-02e6-4f91-8e5f-92ed6b14ab09" />

## Usage

`o7k` uses the standard OpenStack `clouds.yaml` configuration and automatically discovers the 
configuration files from the following places:

* `/etc/openstack`
* `~/.config/openstack`
* the current working directory

> [!Note]
> All the discovered files will be automatically loaded as context in `o7k` 

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
> `o7k` will automatically refresh the current resource every 30s, you don't need to explicitely press `r`.

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

`o7k` supports navigation between related OpenStack resources.

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

When navigating to a related collection, `o7k` filters the destination resource automatically.

<img width="2544" height="1247" alt="Screenshot from 2026-09-18 13-01-45" src="https://github.com/user-attachments/assets/10652fed-f76e-4bd7-bdee-bafe20fb39a7" />

Press `Esc` to return to the previous resource and selection.

## Currently Supported Resources

| Service | Resources | Status |
|---|---|---|
| **Local / Auth** | contexts | Partial |
| **Keystone / Identity** | projects, users, groups, roles, domains | Partial |
| **Keystone / Catalog** | services, endpoints, regions | ✅ |
| **Cinder / Block Storage** | volumes, snapshots, volume types, volume backups | ✅ |
| **Glance / Image** | images | ✅ |
| **Nova / Compute** | servers, flavors | ✅ |
| **Nova / Compute** | keypairs, server groups, availability zones | ✅ |
| **Neutron / Network** | networks, subnets, ports, routers, floating IPs, security groups | ✅ |
| **Neutron / Network** | security-group rules | ✅ |
| **Neutron / Network** | router interfaces | ⬜ |




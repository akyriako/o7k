# 07k

`o7k` (pronounced **oooo-sju-k** 🇸🇪) is a terminal UI for OpenStack, inspired by k9s. 

It provides a fast way to inspect and navigate OpenStack resources directly from the terminal.

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

Press `Esc` to return to the previous resource and selection.

## Currently Supported Resources

| Service | Resources | Status |
|---|---|---|
| **Local / Auth** | contexts | ✅ |
| **Nova / Compute** | servers, flavors | ✅ |
| **Neutron / Network** | networks, subnets, ports, routers, floating IPs, security groups | ✅ |
| **Cinder / Block Storage** | volumes, snapshots, volume types | ✅ |
| **Glance / Image** | images | ✅ |
| **Keystone / Identity** | projects, users, groups, roles, domains | Partial |
| **Keystone / Catalog** | services, endpoints, regions | ✅ |
| **Nova / Compute** | keypairs, server groups, availability zones | ✅ |
| **Neutron / Network** | security-group rules | ✅ |
| **Neutron / Network** | router interfaces | ⬜ |
| **Cinder / Block Storage** | volume backups | ✅ |


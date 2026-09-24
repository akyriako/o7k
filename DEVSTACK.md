# Running a Local OpenStack Cloud with DevStack

This guide sets up a local
[DevStack](https://docs.openstack.org/devstack/latest/) cloud for
developing and testing **o7k** against real OpenStack APIs.

The resulting environment provides the standard OpenStack services
needed for o7k development:

-   Keystone --- Identity
-   Nova --- Compute
-   Neutron --- Networking
-   Glance --- Image service
-   Cinder --- Block Storage
-   Placement --- Resource placement
-   Horizon --- Web dashboard

Optional OpenStack services can be added later as DevStack plugins.

> [!WARNING] 
> DevStack is a development and testing environment. Do not
> install it on a production machine. Use a dedicated or disposable VM.

## Requirements

A VM with the following resources is recommended:

| Resource       | Minimum            | Recommended                    |
|----------------|--------------------|--------------------------------|
| CPU            | 4 vCPU             | 6–8+ vCPU                      |
| RAM            | 16 GB               | 20–32 GB                       |
| Disk           | 60 GB              | 90+ GB                         |
| Virtualization | KVM with nested virtualization               | KVM with nested virtualization |
| OS             | Debian 13 (Trixie) | Debian 13 (Trixie)             |

The instructions above assume Debian 13.

### Verify the VM

Check the operating system:

``` bash
cat /etc/os-release
```

Check CPU virtualization support:

``` bash
lscpu | grep -E 'Virtualization|Model name'
```

Check that KVM is exposed to the VM:

``` bash
ls -la /dev/kvm
```

Check available resources:

``` bash
nproc
free -h
df -h /
```

For nested virtualization on Intel hosts:

``` bash
cat /sys/module/kvm_intel/parameters/nested
```

A result of `Y` means nested virtualization is enabled. You can also verify that the virtualization flags are exposed:

``` bash
grep -Eo 'vmx|svm' /proc/cpuinfo | sort | uniq -c
```

## 1. Install prerequisites

Update the package index:

``` bash
sudo apt-get update
```

Install the basic tools needed to bootstrap DevStack:

``` bash
sudo apt-get install -y git sudo curl
```

## 2. Create the DevStack user

DevStack must run as a regular user with passwordless sudo access.

Create a dedicated `stack` user:

``` bash
sudo useradd -s /bin/bash -d /opt/stack -m stack
```

Configure passwordless sudo:

``` bash
echo "stack ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/stack
sudo chmod 0440 /etc/sudoers.d/stack
```

Validate the sudo configuration:

``` bash
sudo visudo -c
```

### Allow access to KVM

Add the `stack` user to the `kvm` group:

``` bash
sudo usermod -aG kvm stack
```

Verify:

``` bash
id stack
```

The output should include the `kvm` group.

### Fix `/opt/stack` permissions

On Debian, the home directory created for `stack` may have mode `0700`.
OpenStack services run as separate users and must be able to traverse
`/opt/stack`.

Set the directory to `0755`:

``` bash
sudo chmod 0755 /opt/stack
```

Verify:

``` bash
stat -c '%A %a %U:%G %n' /opt /opt/stack
```

Expected:

``` text
drwxr-xr-x 755 root:root /opt
drwxr-xr-x 755 stack:stack /opt/stack
```

If this is not corrected, DevStack may terminate with:

``` text
Invalid path permissions
```

Do not work around the error with `SKIP_PATH_SANITY`. Fix the
permissions instead.

## 3. Verify the `stack` user

Test sudo and KVM access:

``` bash
sudo -iu stack bash -c '
echo "USER=$USER"
echo "HOME=$HOME"
id
sudo -n true && echo "sudo: OK"
test -r /dev/kvm -a -w /dev/kvm && echo "KVM: OK" || echo "KVM: FAILED"
'
```

The important results are:

``` text
sudo: OK
KVM: OK
```

## 4. Clone DevStack

Switch to the `stack` user:

``` bash
sudo -iu stack
```

The working directory should now be:

``` text
/opt/stack
```

Clone DevStack:

``` bash
git clone https://opendev.org/openstack/devstack
cd devstack
```

## 5. Determine the host IP address

DevStack needs the IP address through which the OpenStack APIs will be
reachable.

Run:

``` bash
ip -4 -br addr
ip route
```

For example:

``` text
eth0    UP    10.0.101.97/24
```

and:

``` text
default via 10.0.101.1 dev eth0
```

In this example, the DevStack host IP is:

``` text
10.0.101.97
```

Use the actual IP address of your VM in the configuration below.

> [!IMPORTANT] 
> The DevStack VM should have a stable address. If DHCP
> can assign a different address after reboot, use a DHCP reservation or
> configure a static address.

## 6. Configure DevStack

From `/opt/stack/devstack`, create:

``` text
local.conf
```

with:

``` ini
[[local|localrc]]

# Host address
HOST_IP=10.0.101.97

# Development credentials
ADMIN_PASSWORD=o7k-dev
DATABASE_PASSWORD=o7k-dev
RABBIT_PASSWORD=o7k-dev
SERVICE_PASSWORD=o7k-dev

# Logging
LOGFILE=/opt/stack/logs/stack.sh.log
LOGDAYS=2
```

Replace:

``` text
10.0.101.97
```

with the IP address of your DevStack VM.

The passwords above are deliberately simple development credentials.
Change them if the VM is accessible from an untrusted network.

## 7. Install OpenStack

Start the DevStack installation:

``` bash
./stack.sh
```

This can take some time. DevStack installs the required system packages,
OpenStack source repositories, Python dependencies, databases, message
queues, networking components, and systemd services.

A successful installation ends with output similar to:

``` text
This is your host IP address: 10.0.101.97
Horizon is now available at http://10.0.101.97/dashboard
Keystone is serving at http://10.0.101.97/identity/
The default users are: admin and demo
The password: o7k-dev

Services are running under systemd unit files.
```

## 8. Verify OpenStack authentication

DevStack creates an `openrc` file in `/opt/stack/devstack`.

Load the admin credentials:

``` bash
source openrc admin admin
```

Request a Keystone token:

``` bash
openstack token issue
```

If a token is returned, Keystone authentication is working.

## 9. Configure `clouds.yaml`

o7k uses the standard OpenStack `clouds.yaml` configuration.

Create the configuration directory:

``` bash
mkdir -p ~/.config/openstack
chmod 700 ~/.config/openstack
```

Create:

``` text
~/.config/openstack/clouds.yaml
```

with:

``` yaml
clouds:
  devstack:
    auth:
      auth_url: http://10.0.101.97/identity/v3
      username: admin
      password: o7k-dev
      project_name: admin
      user_domain_name: Default
      project_domain_name: Default
    region_name: RegionOne
    interface: public
    identity_api_version: 3
```

Again, replace `10.0.101.97` with the address of your DevStack VM.

Protect the credentials:

``` bash
chmod 600 ~/.config/openstack/clouds.yaml
```

Verify that OpenStackClient can read the configuration:

``` bash
openstack --os-cloud devstack configuration show
```

Then test authentication:

``` bash
openstack --os-cloud devstack token issue
```

## 10. Verify the OpenStack service catalog

List the registered services:

``` bash
openstack --os-cloud devstack service list
```

A standard installation should expose services similar to:

``` text
+-------------+----------------+
| Name        | Type           |
+-------------+----------------+
| cinder      | block-storage  |
| nova        | compute        |
| placement   | placement      |
| glance      | image          |
| neutron     | network        |
| nova_legacy | compute_legacy |
| keystone    | identity       |
+-------------+----------------+
```

List the API endpoints:

``` bash
openstack --os-cloud devstack endpoint list
```

The important service types for the initial o7k development environment
are:

| Service   | Type            |
|-----------|-----------------|
| Keystone  | `identity`      |
| Nova      | `compute`       |
| Neutron   | `network`       |
| Glance    | `image`         |
| Cinder    | `block-storage` |
| Placement | `placement`     |

## 11. Connect o7k

If o7k runs directly on the DevStack VM, the `clouds.yaml` created above
can be used immediately.

If o7k runs on another development machine, copy the cloud configuration
to:

``` text
~/.config/openstack/clouds.yaml
```

on that machine.

Before starting o7k, verify that the DevStack API is reachable:

``` bash
curl http://10.0.101.97/identity/
```

If OpenStackClient is installed locally, verify the entire
authentication path:

``` bash
openstack --os-cloud devstack token issue
```

Then select/use the `devstack` cloud in o7k.

o7k should now be able to authenticate through Keystone and interact
with the OpenStack services advertised by the service catalog.

## Useful commands

### Show services

``` bash
openstack --os-cloud devstack service list
```

### Show endpoints

``` bash
openstack --os-cloud devstack endpoint list
```

### Request a token

``` bash
openstack --os-cloud devstack token issue
```

### Show DevStack systemd units

``` bash
systemctl list-units 'devstack@*'
```

### View a service log

DevStack services run under systemd. For example:

``` bash
sudo journalctl -u devstack@n-api
```

Follow the log:

``` bash
sudo journalctl -f -u devstack@n-api
```

### Check disk usage

``` bash
df -h /
du -sh /opt/stack/*
```

## Troubleshooting

### `Invalid path permissions`

If `stack.sh` reports:

``` text
*** /opt/stack
*** appears to have 0700 permissions.
*** This is very likely to cause fatal issues for DevStack daemons.
...
Invalid path permissions
```

fix the directory permissions:

``` bash
sudo chmod 0755 /opt/stack
```

Then rerun:

``` bash
./stack.sh
```

### KVM is unavailable to `stack`

Check:

``` bash
ls -la /dev/kvm
id stack
```

If `/dev/kvm` belongs to the `kvm` group but `stack` is not a member:

``` bash
sudo usermod -aG kvm stack
```

Start a new login session for `stack` afterward:

``` bash
exit
sudo -iu stack
```

Verify:

``` bash
test -r /dev/kvm -a -w /dev/kvm && echo "KVM: OK"
```

### API works on the DevStack VM but not remotely

Verify that the remote development machine can reach the DevStack host:

``` bash
ping 10.0.101.97
curl http://10.0.101.97/identity/
```

Also verify that `HOST_IP` in `local.conf` contains an address reachable
from the machine running o7k.

### `clouds.yaml` is not found

The default per-user location is:

``` text
~/.config/openstack/clouds.yaml
```

Check:

``` bash
ls -l ~/.config/openstack/clouds.yaml
```

Then test explicitly:

``` bash
openstack --os-cloud devstack token issue
```

## Adding more OpenStack services

The baseline installation intentionally starts with the standard
DevStack services.

For broader o7k API development, additional OpenStack projects can be
enabled through [DevStack plugins](https://docs.openstack.org/devstack/latest/plugin-registry.html), including:

| Project    | API / Purpose            |
|------------|--------------------------|
| Swift      | Object Storage           |
| Heat       | Orchestration            |
| Designate  | DNS                      |
| Barbican   | Key Management           |
| Octavia    | Load Balancing           |
| Manila     | Shared Filesystems       |
| Trove      | Database as a Service    |
| Magnum     | Container Infrastructure |
| Ironic     | Bare Metal               |
| Aodh       | Alarming                 |
| CloudKitty | Rating                   |
| Blazar     | Reservations             |
| Zaqar      | Messaging                |


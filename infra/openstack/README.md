---
title: OpenStack
description: Two-node OpenStack deployment for the agent farm
tags:
  - ansible
  - kolla
  - openstack
---

This deployment runs OpenStack on two CentOS Stream 10 bare-metal hosts.
Kolla-Ansible runs the services in containers. The deployment does not use
service VMs or Ceph.

## Architecture

```text
master1.openstack1.dc1.alwaldend.com
  control plane, API, database, message bus, OVN gateway, Glance, and Cinder
  better boot NVMe: host OS, container state, databases, logs, and Glance
  2 x SATA SSD: Linux RAID1 -> LVM VG cinder-fast
  4 x DAS HDD: Linux RAID10 -> LVM VG cinder-bulk

compute1.openstack1.dc1.alwaldend.com
  Nova compute, libvirt, OVN controller, and agent instances
  boot NVMe: host OS, container images, and logs
  dedicated 2 TB NVMe: XFS mounted at /var/lib/nova
```

Nova stores local instance data below `/var/lib/nova/instances`. The dedicated
2 TB NVMe therefore holds VM root disks, ephemeral disks, image caches, and
agent workspaces. The compute boot drive does not hold normal VM disks.

The master should use the better boot drive. Its boot drive holds persistent
control-plane data and Glance images. The compute host sends most workload I/O
to its second NVMe.

The `fast` Cinder volume type uses the mirrored SATA SSDs. The `bulk` volume
type uses the four-disk RAID10 array. Use Cinder only for data that must remain
after instance replacement.

The master is one failure domain. RAID protects data from some disk failures.
RAID is not a backup, and it does not make the control plane highly available.

## Network

The two hosts use the same switch and a 1 GbE link. The link carries API,
overlay, image, migration, and Cinder traffic. It limits remote block storage
to approximately one gigabit per second.

The deployment uses OVN. The master provides centralized north-south routing.
The compute node does not need an external-network interface.

The master needs two logical interfaces:

- A management interface with one IPv4 address.
- A separate unnumbered external interface for `br-ex`.

A VLAN trunk can provide both interfaces through one physical NIC. Never set
the external interface to the addressed management interface. OVN can remove
management connectivity when it attaches that interface to `br-ex`.

Use MTU 1500 unless the complete path supports one larger MTU. Configure the
same MTU on both management interfaces and on every switch port in the path.

## Host preparation

Install CentOS Stream 10 on both hosts. Kolla uses Rocky Linux 10 container
images because the selected release does not publish CentOS Stream 10 images.

Create the `ansible` user on both hosts. Grant that user passwordless sudo.
Install the repository Vault SSH client CA as a trusted user certificate CA.
Record both host keys in the deployment account's `known_hosts` file.

Enable CPU virtualization and load KVM on the compute node. Configure working
time synchronization on both hosts.

The DAS must expose all four HDDs as separate Linux block devices. Do not use a
DAS mode that combines the drives behind an opaque RAID volume.

## Configure the inventory

Replace each `CHANGE_ME` value in these files:

- [`ansible/inventory.yaml`](./ansible/inventory.yaml)
- [`ansible/group_vars/all.yaml`](./ansible/group_vars/all.yaml)

Use stable `/dev/disk/by-id/...` paths. Do not use `/dev/sdX` or `/dev/nvmeXnY`
paths for a destructive storage operation.

Reserve one unused address on the management subnet for
`openstack_kolla_internal_vip_address`. Both hosts must have the VIP on-link.
Do not assign the address to another device.

The storage roles enforce this layout:

```text
fast members: 2 whole non-rotating disks
bulk members: 4 whole rotating disks
compute member: 1 whole NVMe disk, at least 1.8 TB
```

The roles reject duplicate paths, partitions, mounted member disks, unexpected
RAID members, degraded arrays, wrong volume-group backing devices, and a Nova
mount from the wrong disk.

## Storage management

Storage creation is disabled by default. In this mode, create and mount the
storage before you run the playbook:

```text
/dev/md/cinder-fast -> RAID1 -> VG cinder-fast
/dev/md/cinder-bulk -> RAID10 -> VG cinder-bulk
2 TB NVMe -> XFS -> /var/lib/nova
```

The playbook still validates every array, volume group, member disk, mount, and
fstab entry.

To let Ansible create the arrays and compute filesystem, set all four flags:

```yaml
openstack_manage_storage: true
openstack_storage_wipe: true
openstack_manage_compute_storage: true
openstack_compute_storage_wipe: true
```

**Verify every by-id path before you set these flags.** These settings permit
`wipefs`, filesystem creation, md array creation, and LVM initialization.

## Vault resources

The change adds the `src_infra_dc1_openstack1` AppRole to the Vault Terraform
configuration. It also grants the AppRole access to the Ansible SSH signing
role.

Review the Vault plan before you apply it:

```sh
bazel run //infra/vault/tf:tf.plan
```

Apply the Vault change through the normal repository process. Do not run the
OpenStack deployment before the AppRole exists.

Kolla needs one complete generated `passwords.yml` document. Store it as the
`passwords_yml` field at this Vault path:

```text
secrets/alwaldend.com/vault1/approles/src_infra_dc1_openstack1/kolla
```

One initialization procedure follows:

```sh
workdir="$(mktemp -d)"
python3 -m venv "${workdir}/venv"
"${workdir}/venv/bin/pip" install \
  'git+https://opendev.org/openstack/kolla-ansible.git@fe4f6adaf01e93af39cd28d3ad57d45b1db11884'
passwords_template="$("${workdir}/venv/bin/python" - <<'PY'
from kolla_ansible import utils

print(utils.get_data_files_path("etc_examples", "kolla", "passwords.yml"))
PY
)"
cp "${passwords_template}" "${workdir}/passwords.yml"
"${workdir}/venv/bin/kolla-genpwd" \
  --passwords "${workdir}/passwords.yml"
jq -Rs '{passwords_yml: .}' "${workdir}/passwords.yml" \
  >"${workdir}/passwords.json"
bazel run //infra/openstack:vault.kv_put \
  alwaldend.com/vault1/approles/src_infra_dc1_openstack1/kolla \
  @"${workdir}/passwords.json"
rm -rf "${workdir}"
```

The injector writes the password file with private permissions. The playbook
stages it only while Kolla runs, then removes the staged copy.

## Deploy

Run host preparation and Kolla prechecks first:

```sh
bazel run //infra/openstack/ansible:ansible.prechecks
```

Deploy the cloud and create the client configuration:

```sh
bazel run //infra/openstack/ansible:ansible.deploy
```

The default target also runs the deploy action:

```sh
bazel run //infra/openstack/ansible
```

Apply service configuration changes:

```sh
bazel run //infra/openstack/ansible:ansible.reconfigure
```

Regenerate client configuration and reconcile Cinder volume types:

```sh
bazel run //infra/openstack/ansible:ansible.post_deploy
```

The post-deploy checks require these active services:

```text
master@fast   cinder-volume enabled/up
master@bulk   cinder-volume enabled/up
compute       nova-compute  enabled/up
```

The checks also require `fast` and `bulk` volume types with matching backend
properties. Cinder uses `fast` as its default volume type.

Kolla state uses this directory pattern:

```text
~/.local/state/alwaldend/openstack/2026.1-<revision>
```

The Kolla virtual environment uses a revision-specific cache directory. The
playbook verifies the installed Git commit before it runs Kolla.

## Operations

Keep agent root disks and workspaces on local Nova storage. Attach Cinder
volumes only for persistent state. Large Cinder transfers can saturate the
1 GbE link and affect API or overlay traffic.

Back up important Cinder data to a separate failure domain. A disk mirror does
not protect against deletion, controller failure, host failure, or corruption.

Monitor these items:

- SMART health for all SSDs and HDDs.
- md array state and resync progress.
- Free space in both Cinder volume groups.
- Free space below `/var/lib/nova`.
- Master boot-drive health and free space.
- Packet loss and link saturation between both hosts.

## References

- https://docs.openstack.org/kolla-ansible/2026.1/
- https://docs.openstack.org/kolla-ansible/2026.1/user/quickstart.html
- https://docs.openstack.org/kolla-ansible/2026.1/reference/storage/cinder-guide.html
- https://docs.openstack.org/kolla-ansible/2026.1/reference/networking/neutron.html

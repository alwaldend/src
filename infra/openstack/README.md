---
title: OpenStack
description: Two-node OpenStack deployment for the agent farm
tags:
  - ansible
  - kolla
  - openstack
---

This deployment runs OpenStack directly on two CentOS Stream 10 or Rocky Linux
10 bare-metal hosts with Kolla-Ansible. It does not create service VMs and it
does not deploy Ceph.

## Architecture

```text
master1.openstack1.dc1.alwaldend.com
  control, API, database, message bus, networking, Glance and Cinder
  boot NVMe: host OS, control-plane data and Glance images
  2 x SATA SSD: md RAID1 -> LVM VG cinder-fast
  4 x DAS HDD: md RAID10 -> LVM VG cinder-bulk
  DAS requirement: expose all four disks individually to Linux

compute1.openstack1.dc1.alwaldend.com
  Nova compute, libvirt and agent instances
  boot NVMe: host OS
  dedicated 2 TB NVMe: XFS mounted at /var/lib/nova/instances
```

The master is a single failure domain. RAID protects Cinder data from a disk
failure, but it is not a backup and does not make the control plane highly
available.

The 1 GbE link is used for API, overlay, image and Cinder traffic. Agent roots
and workspaces should normally use the compute node's local NVMe. Use Cinder
only for state that must survive instance replacement.

## Network

The deployment uses OVN with centralized north-south routing on the master.
Distributed floating IPs and provider networks on the compute node are disabled,
so the compute node does not need an external-network interface.

Kolla still needs two logical interfaces on the master:

- a management interface with an IP address;
- an unnumbered external interface that OVN can attach to `br-ex`.

With one physical NIC, configure a VLAN trunk on the switch and use separate
VLAN interfaces. Do not point `openstack_external_interface` at the management
interface; doing so can remove the host's management connectivity.

## Configure

1. Replace all `CHANGE_ME` values in
   [`ansible/inventory.yaml`](./ansible/inventory.yaml) and
   [`ansible/group_vars/all.yaml`](./ansible/group_vars/all.yaml).
2. Use stable `/dev/disk/by-id/...` device paths. Never use `/dev/sdX` names for
   the destructive storage setup.
3. Configure DNS for the two inventory names, or set `ansible_host` for each
   host.
4. Reserve an unused management-subnet address for
   `openstack_kolla_internal_vip_address`.
5. Create the Vault AppRole `src_infra_dc1_openstack1` and permit it to read the
   Kolla secret path and request the configured SSH certificate.

Storage creation is disabled by default. For arrays that already exist, create
VGs named `cinder-fast` and `cinder-bulk`, mount the compute NVMe at
`/var/lib/nova/instances`, and leave the management flags disabled.

To let Ansible create the arrays and compute filesystem, set all four flags only
after verifying every device ID:

```yaml
openstack_manage_storage: true
openstack_storage_wipe: true
openstack_manage_compute_storage: true
openstack_compute_storage_wipe: true
```

Those settings authorize `wipefs` on the listed devices.

## Kolla passwords

Kolla's complete generated `passwords.yml` is stored as one Vault value named
`passwords_yml`; it is never committed. One way to initialize it is:

```sh
VENV="$(mktemp -d)/venv"
python3 -m venv "${VENV}"
"${VENV}/bin/pip" install \
  'git+https://opendev.org/openstack/kolla-ansible@stable/2026.1'
PASSWORDS_TEMPLATE="$("${VENV}/bin/python" - <<'PY'
from kolla_ansible import utils
print(utils.get_data_files_path("etc_examples", "kolla", "passwords.yml"))
PY
)"
cp "${PASSWORDS_TEMPLATE}" /tmp/openstack-passwords.yml
"${VENV}/bin/kolla-genpwd" --passwords /tmp/openstack-passwords.yml
jq -Rs '{passwords_yml: .}' /tmp/openstack-passwords.yml \
  >/tmp/openstack-passwords.json
bazel run //infra/openstack:vault.kv_put \
  alwaldend.com/vault1/approles/src_infra_dc1_openstack1/kolla \
  @/tmp/openstack-passwords.json
rm -f /tmp/openstack-passwords.yml /tmp/openstack-passwords.json
rm -rf "$(dirname "${VENV}")"
```

## Deployment

The Bazel execution host needs `/usr/bin/python3` with virtual-environment
support, Git, and the normal native build dependencies required by
Kolla-Ansible's Python packages.

Run host preparation and Kolla prechecks first:

```sh
bazel run //infra/openstack/ansible:ansible.prechecks
```

With both storage-management flags disabled, this validates pre-existing arrays,
volume groups and the Nova mount without wiping disks. When either management
flag is enabled, the same target also performs the explicitly authorized storage
provisioning before Kolla's checks.

Deploy and generate `clouds.yaml`:

```sh
bazel run //infra/openstack/ansible:ansible.deploy
```

The default target is equivalent to `ansible.deploy`:

```sh
bazel run //infra/openstack/ansible
```

Apply configuration changes without rebuilding the cloud:

```sh
bazel run //infra/openstack/ansible:ansible.reconfigure
```

Regenerate post-deployment client configuration and reconcile the `fast` and
`bulk` Cinder volume types:

```sh
bazel run //infra/openstack/ansible:ansible.post_deploy
```

Runtime Kolla state is kept below
`~/.local/state/alwaldend/openstack/2026.1`; the Kolla virtual environment is
kept below `~/.cache/alwaldend/openstack`. The Vault-rendered password file is
copied into the state directory with mode `0600`.

## References

- https://docs.openstack.org/kolla-ansible/2026.1/
- https://docs.openstack.org/kolla-ansible/2026.1/user/quickstart.html
- https://docs.openstack.org/kolla-ansible/2026.1/reference/storage/cinder-guide.html
- https://docs.openstack.org/kolla-ansible/2026.1/reference/networking/neutron.html

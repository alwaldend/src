# Proxmox infrastructure Specification

## Purpose

Specify the Proxmox cluster configuration, resource pools, and bootstrap test
VM owned by `infra/pve`. The baseline is checked-in source at revision
`550d7e79b1f5fdbc2b6017b75178471d6914082f`, observed on 2026-09-08.
It describes declared infrastructure and packaged operator entry points;
deployment and cluster health have not been observed for this baseline.

## Requirements

### Requirement: AppRole resource pools

The Terraform configuration SHALL resolve the members of Vault's `approles`
identity group and declare one Proxmox pool per resolved entity, keyed and
named by that entity's name. It SHALL also declare a separate `templates` pool.

Sources: [project documentation](../../../README.md) and
[resource pool definitions](../../../tf/pools.tf).

#### Scenario: Declare pools for known AppRoles

- **WHEN** Terraform resolves the entities in the Vault `approles` group
- **THEN** the desired configuration includes a pool with each entity's name
  and the dedicated `templates` pool.

### Requirement: Packaged host and snippet configuration

The PVE Ansible entry point SHALL package the inventory and deployment playbook
and apply the shared host and PVE roles to the `pve1` inventory group. It SHALL
package the canonical shared cloud-init files under `files/cloud_init.yaml`
and `files/cloud_init_min.yaml`, preserving the existing snippet update flow.

Sources: [deployment target](../../../ansible/BUILD.bazel),
[playbook](../../../ansible/playbook_deploy.yaml), and
[snippet variables](../../../ansible/group_vars/all.yaml).

#### Scenario: Assemble the PVE deployment package

- **WHEN** the Ansible binary's source package is built
- **THEN** it includes the playbook, inventory, shared role collection, and
  both shared cloud-init files under their established aliases.

### Requirement: Cloud-init test virtual machine

The Terraform configuration SHALL declare its cloud-init test VM through the
shared `pve_vm_qemu` module, derive its hostname and address from the project's
DNS configuration, and place it in the `src_infra_dc1_pve1` AppRole pool. The
baseline declaration SHALL use VM ID 100, a 5G `scsi0` disk, and the `test` tag.

Sources: [VM definition](../../../tf/vms.tf) and
[DNS configuration](../../../dnsconfig.json).

#### Scenario: Resolve the test VM's desired configuration

- **WHEN** Terraform evaluates `module.vm_cloudinit_test`
- **THEN** its name and `/24` address derive from `cloudinit_test` DNS fields,
  and its pool resolves to the declared PVE AppRole pool.

### Requirement: Repository-managed execution inputs

PVE Terraform commands SHALL use the project's AL configuration and the
declared Terraform, PVE-login, and Vault-environment plugins. The package SHALL
include the lockfile, Terraform source, shared backend and VM modules, and
required Vault helper targets. PVE Ansible commands SHALL likewise use the
project AL configuration and packaged injector.

Sources: [Terraform wrapper](../../../tf/BUILD.bazel) and
[Ansible wrapper](../../../ansible/BUILD.bazel).

#### Scenario: Construct a PVE Terraform invocation

- **WHEN** a repository Terraform command target is invoked with authorization
  for its operation
- **THEN** the wrapper supplies the declared modules and helper inputs and
  selects the `tf`, `pve_login`, and `vault_env=default_default` plugins.

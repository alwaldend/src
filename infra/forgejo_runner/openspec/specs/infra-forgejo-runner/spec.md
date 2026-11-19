# Forgejo runner infrastructure specification

## Purpose

Describe the Forgejo Actions runner VM and deployment scaffold owned by
`infra/forgejo_runner`. Its `tf_setup` and `ansible` packages are stages of
this owner. The runner role is commented out in the baseline playbook;
the checked-in configuration therefore does not establish a deployed or
registered Actions worker. No live runner state was observed.

Baseline revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observed: 2026-09-08. Sources:
[owner README](../../../README.md),
[BUILD](../../../BUILD.bazel), and the implementation
links below.

## Requirements

### Requirement: Dedicated runner VM allocation

Setup SHALL declare `runner1` as Proxmox VM 1200 in the
`src_infra_forgejo_runner` pool, with eight cores, 16 GiB memory, a 20 GiB
boot disk, and a 300 GiB runner disk on `ceph-ec` storage. Its hostname and
address SHALL derive from the owner's DNS declaration.

Source: [VM definition](../../../tf_setup/vms.tf).

#### Scenario: Inspect the runner capacity definition

- **WHEN** a maintainer evaluates the runner VM configuration
- **THEN** the source declares the dedicated runner pool, compute allocation,
  and separate boot and runner disks
- **AND** DNS source supplies the hostname and address

### Requirement: Explicit host-only deployment baseline

The baseline deployment playbook SHALL apply the shared host role to the
`forgejo_runner` inventory group. The specification MUST identify the
commented-out Forgejo runner role as inactive and MUST NOT interpret the
available runner variables as evidence that the worker is installed or
registered.

Sources: [deployment playbook](../../../ansible/playbook_deploy.yaml)
and [runner variables](../../../ansible/group_vars/all.yaml).

#### Scenario: Determine the effect of the current playbook

- **WHEN** the current playbook's active roles are inspected
- **THEN** only the shared host role is enabled
- **AND** runner installation and registration remain outside the behavior
  established by this source baseline

### Requirement: Packaged setup and deployment configuration

The owner SHALL expose Terraform setup and packaged Ansible entry points
through its `al` configuration and Vault injection dependencies. The runner
configuration SHALL reference `https://git.alwaldend.com`, `/dev/sdb`, and
an environment-supplied `FORGEJO_RUNNER_TOKEN` rather than a checked-in
registration token.

Sources: [setup BUILD](../../../tf_setup/BUILD.bazel),
[Ansible BUILD](../../../ansible/BUILD.bazel), and
[runner variables](../../../ansible/group_vars/all.yaml).

#### Scenario: Resolve the registration credential source

- **WHEN** a maintainer inspects the runner configuration scaffold
- **THEN** its token value is an environment lookup and its instance URL is
  the canonical Forgejo service
- **AND** the token reference alone does not enable the inactive runner role

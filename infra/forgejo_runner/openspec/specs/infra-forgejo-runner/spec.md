# Forgejo runner infrastructure specification

## Purpose

Define the XCP-ng runner allocation, Vault-backed registration, and packaged
host/dev_vm/runner deployment owned by `infra/forgejo_runner`.

Sources: [owner README](../../../README.md),
[Terraform setup](../../../tf_setup/README.md),
[Ansible deployment](../../../ansible/README.md), and
[deployment evidence](../../changes/archive/2026-09-15-deploy-secure-runner/design.md).

## Requirements

### Requirement: Dedicated runner VM allocation

Setup SHALL provision an extensible map of XCP-ng runner VMs, initially only
`secure` with eight vCPUs, 16 GiB memory and a 500 GiB root disk. Its hostname
and address SHALL derive from the owner's DNS declaration.

#### Scenario: Inspect the runner capacity definition

- **WHEN** setup uses its checked-in defaults
- **THEN** exactly one XCP-ng runner is provisioned with the requested resources
- **AND** DNS supplies its hostname and static address

### Requirement: Explicit host-only deployment baseline

The deployment SHALL actively configure the shared host, dev_vm and Forgejo
runner roles. Repeated deployment SHALL retain registration and ensure the
runner service is enabled and healthy.

#### Scenario: Determine the effect of the current playbook

- **WHEN** the packaged Ansible deployment runs against the provisioned VM
- **THEN** the host receives development tools and an enabled runner service
- **AND** registration credentials are excluded from task output

### Requirement: Packaged setup and deployment configuration

The owner SHALL expose Terraform and Ansible commands through its existing
AL and Vault integration. The runner registration credential SHALL remain in
Vault and arrive through injection; CI SHALL receive no provisioning identity.

#### Scenario: Resolve the registration credential source

- **WHEN** the registration playbook retrieves the repository-scoped credential
- **THEN** it preserves other Vault config fields and uses compare-and-set
- **AND** neither the Forgejo controller token nor Vault token is installed on
  the VM

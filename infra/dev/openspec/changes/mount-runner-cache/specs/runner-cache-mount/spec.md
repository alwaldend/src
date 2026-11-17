## Purpose

Define dedicated NAS-backed action-cache storage for the secure runner while preserving existing application and boot data.

## ADDED Requirements

### Requirement: NAS-backed guest cache disk

NAS-owned source under infra/nas/approles/src_infra_dev SHALL declare cache storage
through the reusable dataset module without input variables. That AppRole project
SHALL own its namespace and cache in independent state with its own plan/apply
targets, using NAS operator authentication. The VM root SHALL contain no NAS provider or credentials
and SHALL attach an 800-GiB disk from the provisioned repository. Terraform SHALL publish its actual guest disk device and capacity to the development AppRole's Vault namespace. The packaged Ansible wrapper SHALL inject the record automatically without a caller-supplied output or extra-vars step. Ansible SHALL consume those values, reject ambiguous or unsafe devices, preserve existing filesystems, mount it persistently and require that mount before the worker starts. Bazel SHALL use its action cache with 700-GiB GC while keeping output bases local.

#### Scenario: Review desired state

- **WHEN** the configuration is validated without deployment
- **THEN** The boot disk and dev1 configuration remain unchanged; the guest receives neither SMB credentials nor network-filesystem clients.

#### Scenario: Storage is unavailable

- **WHEN** the separately authorized deployment cannot establish the required storage
- **THEN** the worker does not start with an unmounted local cache

#### Scenario: Metadata is absent or unsafe

- **WHEN** metadata is missing, malformed or names an invalid device or capacity
- **THEN** Ansible fails before formatting or mounting any disk

### Requirement: CI storage diagram

The development documentation SHALL include a CI diagram using the development-host
diagram's site style. Its runner subgraph SHALL contain the cache disk, with
edges from the runner to the disk and from the disk to NAS. External dependency
services SHALL be represented by one node labeled External deps.
The VM SHALL connect to the VPN/proxy node, which connects to External deps.
A dark WebP SHALL be available for download.

#### Scenario: Read CI architecture

- **WHEN** the documentation assets are built
- **THEN** both light and dark variants show the runner's NAS-backed disk and a single External deps node with readable labels

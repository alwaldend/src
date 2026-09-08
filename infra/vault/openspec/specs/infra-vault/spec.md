# Infrastructure Vault

## Purpose

Describe the Vault host configuration, authentication and certificate services,
and packaged recovery entry points. This baseline concerns checked-in desired
state and supported wrapper structure; it does not establish live unsealed
state, successful backups, quorum, or certificate freshness.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [operator documentation](../../../README.md),
[operational wrappers](../../../BUILD.bazel),
[Ansible entry points](../../../ansible/BUILD.bazel),
[Vault host template](../../../ansible/files/vault.hcl),
[authentication backends](../../../tf/auth.tf),
[storage engines](../../../tf/mounts.tf),
[server PKI](../../../tf/pki_servers.tf),
[client PKI](../../../tf/pki_clients.tf), and
[public CA outputs](../../../tf/outputs.tf).

## Requirements

### Requirement: Separate host provisioning from Vault configuration

The documented workflow SHALL separate `tf_setup` VM provisioning, Ansible host
setup, and the `tf` Vault configuration stage. Ansible SHALL expose combined,
VM-only, and bare-metal-only setup entry points using the shared host and Vault
roles.

#### Scenario: Prepare an authorized host setup

- **WHEN** an operator selects `ansible.vm` or `ansible.bm`
- **THEN** the wrapper selects the matching checked-in playbook and inventory
- **AND** Vault's API configuration remains owned by the separate Terraform stage

### Requirement: Configure TLS endpoints and Raft storage

The host template SHALL configure TLS API and cluster endpoints on ports 8200
and 8201 and Raft storage under `/opt/vault/raft`. It SHALL derive each Raft node
identifier from the inventory hostname and emit retry-join addresses for hosts
in the Vault inventory group.

#### Scenario: Render a Vault host configuration

- **WHEN** Ansible renders `vault.hcl` for a member of the Vault group
- **THEN** its API and cluster URLs use its inventory hostname
- **AND** Raft receives that hostname as node ID and the group's HTTPS join URLs

### Requirement: Declare authentication and secret-engine boundaries

Terraform SHALL declare `userpass`, `approle`, and `cert` authentication
backends, a version-two KV engine at `secrets`, and a transit engine at
`transit/default`. Client and server certificate authorities SHALL remain
distinct, and server PKI ACME SHALL require external account binding with the
default directory policy set to `forbid`.

#### Scenario: Inspect the configured credential services

- **WHEN** the main Terraform configuration is evaluated
- **THEN** authentication, KV, transit, client PKI, and server PKI are distinct
  declared resources
- **AND** server ACME is enabled with `eab_policy = "always-required"`
- **AND** published CA output files contain certificate or public-key material

### Requirement: Package recovery and certificate operations explicitly

The package SHALL expose separate wrappers for backup, ordinary unseal,
standalone unseal, client certificate generation, and root-token generation.
Standalone unseal SHALL use the `default_no_auth` environment without selecting
the ordinary unseal plugin; these wrappers SHALL NOT imply that recovery has
been executed successfully.

#### Scenario: Select the standalone unseal workflow

- **WHEN** an authorized operator invokes `//infra/vault:unseal_standalone`
- **THEN** the wrapper selects the packaged unseal utility and declared Vault
  endpoint with its no-auth environment
- **AND** success remains an operational result to be observed separately

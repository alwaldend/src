# Flux infrastructure specification

## Purpose

Define the Flux DNS, Vault identities, minimal K3s bootstrap, and Terraform
workload configuration.

## Requirements

### Requirement: Standalone Flux DNS management

The Flux infrastructure project SHALL expose its site-local DNS declarations
through a standalone Terraform package authenticated with the existing
`src_infra_flux` Vault AppRole and its dedicated DNS state backend. The
project SHALL manage the Flux-specific AppRoles, policies, identity groups,
PKI, OIDC, and Transit configuration from a separate local root.

#### Scenario: Select the Flux DNS package

- **WHEN** a maintainer selects the Flux DNS Terraform package
- **THEN** it consumes the existing Flux DNS declaration and uses the
  component AppRole and DNS backend
- **AND** no Flux VM, host deployment, or cluster reconciliation package
  is required

#### Scenario: Reconcile a shared grant

- **WHEN** shared Vault reconciles a Forgejo, Harbor, or global grant for Flux
- **THEN** it resolves the Flux identities by their Vault names
- **AND** the Flux AppRole root owns the identities and component policies

### Requirement: Flux VM provisioning

Flux SHALL provision its guest through Xen Orchestra with an invocation-scoped
session issued for the Flux Vault AppRole.

#### Scenario: Provision the guest

- **WHEN** an operator applies the VM root
- **THEN** the provider uses the Flux resource set and resolves the guest address from the Flux DNS declaration

### Requirement: Minimal K3s bootstrap

Flux SHALL bootstrap a K3s host with the repository-pinned binary and
configure Vault OIDC for the API. The K3s auto-deploy directory SHALL
contain only the project RBAC binding that authorizes Terraform to
administer the cluster.

#### Scenario: Configure the Flux VM

- **WHEN** an operator runs the Flux bootstrap playbook
- **THEN** it installs K3s from the packaged binary and binds the Flux Vault administrator group

### Requirement: Terraform workload lifecycle

Terraform SHALL manage the Flux Operator and FluxInstance as ordered releases
using a dedicated state backend.

#### Scenario: Deploy to a bootstrapped cluster

- **WHEN** an operator applies the workload root
- **THEN** Terraform installs the mirrored operator and its CRDs before submitting the FluxInstance
- **AND** the launcher opens its SSH tunnel and supplies the Vault OIDC token and server CA
- **AND** the instance selects the mirrored Flux controller images

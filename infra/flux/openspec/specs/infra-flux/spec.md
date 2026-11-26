# Flux infrastructure specification

## Purpose

Define the site-local Flux DNS records and project-owned Vault configuration,
including how shared Vault preserves existing identities during state transfer.

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

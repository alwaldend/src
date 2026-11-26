## ADDED Requirements

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

## REMOVED Requirements

### Requirement: Dedicated K3s management host

**Reason**: Flux infrastructure is retained for DNS and Vault authentication.

**Migration**: Consult the prior revision for historical VM and K3s source.

### Requirement: Git reconciliation for Flux and Harbor

**Reason**: The Flux management cluster's reconciliation tree is retired.

**Migration**: Harbor's own configuration remains under `infra/harbor`.

### Requirement: Packaged operator and controller configuration

**Reason**: The Flux cluster package is retired.

**Migration**: Consult the prior revision for historical operator commands.

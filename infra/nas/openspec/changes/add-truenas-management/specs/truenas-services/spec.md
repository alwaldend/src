## Purpose

Manage the existing NAS storage and Nextcloud lifecycle through simple Terraform
source, with state and credentials independent of DNS management.

## ADDED Requirements

### Requirement: Independent service operations

The NAS owner SHALL expose packaged plan, import, and interactive apply commands
with independent state, Vault-injected credentials, and verified TLS.

#### Scenario: Select service configuration

- **WHEN** an operator selects a TrueNAS command
- **THEN** it uses the NAS service backend and connection entry
- **AND** it does not select DNS credentials

### Requirement: Direct declarations for adopted resources

The source SHALL directly declare the existing pool, seven filesystem datasets,
and Nextcloud lifecycle metadata without a generic variable interface or
separate configuration document. It SHALL preserve their imported resource
addresses and deletion protection. App configuration and encryption secrets
SHALL remain outside source.

#### Scenario: Simplify the adopted configuration

- **WHEN** the source is simplified
- **THEN** the nine imported resource addresses remain unchanged
- **AND** a live plan reports no resource changes or state moves

### Requirement: Separate source validation from live operations

Offline checks SHALL validate formatting and the pinned provider schema without
NAS or Vault access. Live validation SHALL use the owning normal-text plan.

#### Scenario: Validate source offline

- **WHEN** formatting and schema checks run
- **THEN** no live credentials or appliance connection are required

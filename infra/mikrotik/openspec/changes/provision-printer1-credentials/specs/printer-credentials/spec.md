## Purpose

Define printer account provisioning and the ownership and credential guarantees required by the printer infrastructure migration.

## ADDED Requirements

### Requirement: Printer account provisioning

The accounts root SHALL create a dedicated printers RouterOS user and group with read, write, API and REST API access, generate a stable sensitive password, and publish its credentials in owner-controlled Vault storage. Export snapshots SHALL remain documentation inputs.

#### Scenario: Inspect declared ownership

- **WHEN** the account root is provisioned through its owning workflow
- **THEN** the printers identity and its Vault credential agree
- **AND** repeated plans retain the password unless deliberately replaced

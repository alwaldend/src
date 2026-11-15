# Development DNS specification

## Purpose

Define one DNS deployment for both development hosts and dev1's services.

## Requirements

### Requirement: One owner-local DNS root

The dns root SHALL consume the complete canonical dnsconfig.json through the
shared DNS module, preserving all declared views. It SHALL use src_infra_dev
and the alwaldend.com/vault1/approles/src_infra_dev/tf_backend/dns backend.
There SHALL be no separate runner DNS root or record filter.

#### Scenario: Evaluate the declaration

- **WHEN** the DNS root is evaluated
- **THEN** secure1, dev1 and dev1's service records belong to the same root
- **AND** provider credentials are injected through the shared owner identity

### Requirement: Explicit ownership handoff

The owner SHALL document the pending transfer from both legacy states and
preserve legacy access until separately authorized reconciliation is complete.
Adoption SHALL bind existing DNS records without destructive recreation or
concurrent ownership. Unrelated legacy resources SHALL NOT be transferred or
destroyed by implication. After adoption, unchanged declarations SHALL produce
a no-change plan, and temporary import configuration SHALL be removed.

#### Scenario: Review source before adoption

- **WHEN** a maintainer inspects the merged root
- **THEN** the legacy backend locations and logical-key changes are documented
- **AND** source validation makes no claim of live ownership or deployment

### Requirement: Offline source checks

The package SHALL expose a format check and package both DNS provider binaries,
the canonical declaration and shared module without live authentication.

#### Scenario: Validate source

- **WHEN** offline checks run
- **THEN** they validate packaged source without contacting Vault or DNS providers

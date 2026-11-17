## Purpose

Provide secure1 with NAS-backed virtual disk storage while keeping storage
credentials and protocol handling outside the guest and its build jobs.

## ADDED Requirements

### Requirement: Dedicated cache storage repository

The configuration SHALL bootstrap a dedicated SMB repository in secure1's pool
and grant the development resource set access. It SHALL use the hypervisor AppRole's own Vault-backed NAS account. Access SHALL
be delegated on the cache dataset without reading another AppRole's password. Conflicting or duplicate identities SHALL fail without
mutation. Existing matching storage SHALL be preserved on repeat invocation.

#### Scenario: Create and reuse storage

- **WHEN** the authorized bootstrap runs against an absent repository
- **THEN** it creates and verifies one owned SMB repository
- **AND** a repeated invocation preserves the existing repository and disks

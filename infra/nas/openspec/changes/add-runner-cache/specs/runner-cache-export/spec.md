## Purpose

Provide private encrypted AppRole storage with independently managed lifecycles.

## ADDED Requirements

### Requirement: AppRole accounts

The shared NAS root SHALL create one SMB account for each identity in Vault's
truenas_admins and truenas_user groups. Only the administrator group SHALL grant
FULL_ADMIN. Passwords SHALL contain 64 generated characters and one rotation
version SHALL update generation and NAS together. Each AppRole SHALL read only
its own published NAS credential through the existing shared-read policy.

#### Scenario: Add a NAS group member

- **WHEN** a new identity belongs to either canonical NAS group
- **THEN** the shared root provisions its account and publishes its credential
- **AND** namespace provisioning remains owned by that AppRole's storage project

### Requirement: Independently managed AppRole storage

Each infra/nas/approles/<approle> project SHALL own its namespace, ACLs, shares
and child datasets with its own Terraform state and plan/apply targets. Roots
SHALL use NAS operator authentication without exposing its API credential to
the development AppRole. The shared root SHALL retain appliance-wide settings
and the encrypted parent, with no AppRole namespace loop or runner_storage call.
AppRole roots SHALL have no Terraform input variables or remote-state dependency.
Adopted appliance resource addresses SHALL remain unchanged.

#### Scenario: Plan development storage independently

- **WHEN** the development NAS project's plan target is selected
- **THEN** it selects only its own backend and owns its namespace and runner_cache
- **AND** the shared NAS root and other AppRole states are outside that plan

#### Scenario: Use the cache account

- **WHEN** XCP-ng accesses the development namespace using its own NAS account
- **THEN** its named grant permits traversal to the cache and modification there
- **AND** it cannot connect to the namespace share without an explicit MODIFY grant
- **AND** SMB encryption is required by the shared NAS transport settings

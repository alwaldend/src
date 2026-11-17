## ADDED Requirements

### Requirement: Independent AppRole state

Selected Vault AppRole projects SHALL have separate state and owning targets.
Core group memberships SHALL use their existing identity IDs through name lookups.
Migration SHALL preserve all live objects, policy contents and permissions.

#### Scenario: Move an existing AppRole

- **WHEN** its resources move to an independent state
- **THEN** normal resource plans propose no creation, replacement or deletion
- **AND** global memberships retain the same entity/group IDs

### Requirement: Guarded state recovery

Recovery SHALL use the preserved state and native Terraform state moves/pushes
without force. Resources SHALL have complete disjoint ownership after migration.
The recovery SHALL NOT restart or reload Vault or change the storage format.

#### Scenario: A target state already exists or core progress changed

- **WHEN** either condition is observed before migration
- **THEN** migration stops for inspection without overwriting remote progress

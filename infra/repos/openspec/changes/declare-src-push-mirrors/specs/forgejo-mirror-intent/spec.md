## ADDED Requirements

### Requirement: Shared mirror intent and identity

Each repository catalog file SHALL declare its native Forgejo destinations
under `forgejo.push_mirrors`. Each destination SHALL identify a username and
Vault credential reference.

#### Scenario: Add another Forgejo repository mirror

- **WHEN** the catalog declares a mirror and the corresponding repository entries
- **THEN** the mirror Terraform root declares a resource with a stable catalog key
- **AND** the bootstrap uses that destination's repository identity and dedicated Vault path

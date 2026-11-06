## Purpose

Compose multiple Vast.ai teams and project keys without duplicated resource declarations.

## ADDED Requirements

### Requirement: Team and project key module

By default, the module SHALL create one named team with zero credit transfer and any declared
project API keys. Key outputs SHALL be sensitive. The caller SHALL own Vault
authentication and storage paths. The module SHALL declare no compute rentals.

#### Scenario: Multiple teams

- **WHEN** a caller declares two module instances with different team names and key maps
- **THEN** each creates its own team and keys, unchanged plans are empty, and destruction revokes their keys and deletes only their owned teams.

### Requirement: Existing-team key composition

When its team resource is imported, the module SHALL create the declared
keys using its selected provider and imported team ID. It SHALL preserve the existing team and leave
Vault authentication and storage to the caller. Managed-team upgrades SHALL
preserve existing team and key identities.

#### Scenario: Existing funded team

- **WHEN** a caller imports an existing team and declares a full-access key map
- **THEN** no team is created or destroyed and the generated keys are sensitive outputs
- **AND** key-only cleanup revokes only those declared keys and team destruction is refused

#### Scenario: Upgrade a managed-team module

- **WHEN** an existing managed-team module upgrades to the import-capable module
- **THEN** its resource addresses remain unchanged without team or key replacement

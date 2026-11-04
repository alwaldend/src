## Purpose

Compose multiple Vast.ai teams and project keys without duplicated resource declarations.

## ADDED Requirements

### Requirement: Team and project key module

The module SHALL create one named team with zero credit transfer and any declared
project API keys. Key outputs SHALL be sensitive. The caller SHALL own Vault
authentication and storage paths. The module SHALL declare no compute rentals.

#### Scenario: Multiple teams

- **WHEN** a caller declares two module instances with different team names and key maps
- **THEN** each creates its own team and keys, unchanged plans are empty, and destruction revokes their keys and deletes only their owned teams.

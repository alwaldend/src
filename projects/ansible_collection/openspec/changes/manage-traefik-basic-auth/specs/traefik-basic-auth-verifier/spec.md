## ADDED Requirements

### Requirement: Traefik Basic Auth verifier

The Traefik role SHALL install a caller-supplied Basic Auth verifier with
mode `0640` and include it in service activation and rollback.

#### Scenario: Rotate the verifier

- **WHEN** the caller supplies a changed verifier
- **THEN** the role installs it, activates the updated proxy, and restores the previous verifier if activation fails.

#### Scenario: Configure Traefik without a verifier

- **WHEN** the caller supplies no verifier
- **THEN** the role installs no verifier file.

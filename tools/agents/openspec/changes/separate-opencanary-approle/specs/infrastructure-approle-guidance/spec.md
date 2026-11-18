## ADDED Requirements

### Requirement: Standalone AppRole guidance

The infrastructure skill SHALL require a dedicated Terraform root and backend for each new infrastructure AppRole.

#### Scenario: Add an infrastructure component

- **WHEN** a component needs a new Vault AppRole
- **THEN** the skill directs the author to a standalone project under `infra/vault/approles/` and explains how the shared stage resolves its identity

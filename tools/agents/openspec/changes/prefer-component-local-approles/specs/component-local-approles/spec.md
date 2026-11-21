## ADDED Requirements

### Requirement: Component-local placement

The infrastructure skill SHALL place new component-specific AppRole roots
in the owning infrastructure project and use the shared Vault
administrative configuration.

#### Scenario: Add a component AppRole

- **WHEN** an infrastructure component needs a new Vault AppRole
- **THEN** its Terraform root lives under the component and uses a dedicated Vault HTTP backend

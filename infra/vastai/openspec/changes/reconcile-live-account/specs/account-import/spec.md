## Purpose

Allow existing Vast.ai resources to be imported through the component's
established authentication and secure Terraform backend.

## ADDED Requirements

### Requirement: Authenticated resource import

The infrastructure root SHALL expose a Terraform import command using the same
credentials, providers, configuration inputs, and backend as its plan command.
Imports SHALL remain limited to resource types supported by the provider.

#### Scenario: Import an existing supported resource

- **WHEN** an operator imports a supported resource with its declaration and identifier
- **THEN** Terraform records it in the infrastructure root's existing Vault-backed state
- **AND** the command does not require an API credential in its arguments

### Requirement: Personal-account HTTPS setting

The infrastructure root SHALL manage the personal-account ENABLE_HTTPS
environment variable with the string value `true`.

#### Scenario: Adopt the existing HTTPS setting

- **WHEN** the existing variable is imported and its reviewed value update is applied
- **THEN** its identity is preserved and a subsequent plan reports no changes

### Requirement: Preserve existing team identities

The root SHALL declare the three existing teams by observed name, retain their
imported account IDs in state, and prevent destruction without generating keys.
Numeric account IDs SHALL NOT be duplicated in resource configuration.
All team definitions SHALL use a single catalog of per-team configuration objects.
Import blocks SHALL be removed after successful adoption. Operator account
keys SHALL remain outside Terraform ownership and SHALL NOT be revoked.

#### Scenario: Adopt existing teams and remove an unmanaged team key

- **WHEN** the observed teams are reconciled and the authorized unmanaged team API key is revoked
- **THEN** all declared teams retain their account IDs and a subsequent plan reports no changes
- **AND** account keys, automatic owner contexts, and Terraform-owned keys are preserved

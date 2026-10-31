## Purpose

Provide this project's independently packaged Terraform DNS stage while preserving
its canonical record declarations and explicit migration activation boundary.

## ADDED Requirements

### Requirement: Project-owned DNS stage

The project SHALL expose a Terraform stage consuming its existing DNS declaration
and using its own Vault authentication and Terraform state.

#### Scenario: Package the project stage

- **WHEN** the project's Terraform command is built
- **THEN** its runfiles include the project declaration, shared DNS module, and
  the configured authentication, state, and provider injection plugins.

### Requirement: Explicit ownership activation

The DNS stage SHALL declare no managed records until explicitly enabled for an
authorized adoption. Its public-only integration SHALL require only Cloudflare,
and operational commands SHALL select the existing stage's provider injection.

#### Scenario: Inspect the prepared source defaults

- **WHEN** the project stage uses its checked-in source defaults
- **THEN** record ownership is disabled and the zone input is optional, while
  operational provider setup and authentication still require their Vault prerequisites.

#### Scenario: Adopt existing public DNS records

- **WHEN** an authorized adoption freezes the old central writers and enables
  this project's record ownership
- **THEN** the stage uses its project state and Cloudflare configuration to
  import existing records without a RouterOS dependency.

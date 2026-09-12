# project-dns Specification

## Purpose

Provide this project's independently packaged Terraform DNS stage while preserving
its canonical record declarations and explicit migration activation boundary.

## Requirements

### Requirement: Project-owned DNS stage

The project SHALL expose a Terraform stage consuming its existing DNS declaration
and using its own Vault authentication and Terraform state.

#### Scenario: Package the project stage

- **WHEN** the project's Terraform command is built
- **THEN** its runfiles include the project declaration, shared DNS module, and
  the configured authentication, state, and provider injection plugins.

### Requirement: Explicit ownership activation

The DNS stage SHALL retain enabled record ownership after authorized adoption
into the project's state. Its public-only integration SHALL require only
Cloudflare, and operational commands SHALL select the existing stage's provider
injection.

#### Scenario: Inspect the prepared source defaults

- **WHEN** the project's prepared stage has not yet completed authorized adoption
- **THEN** record ownership stays disabled until its adoption revision, while
  operational provider setup and authentication require the Vault credentials.

#### Scenario: Adopt existing public DNS records

- **WHEN** an authorized adoption freezes the old central writers and enables
  this project's record ownership
- **THEN** the stage uses its project state and Cloudflare configuration to
  import existing records without a RouterOS dependency.

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted project stage uses its checked-in source defaults
- **THEN** record ownership is enabled and the zone input remains optional.

#### Scenario: Reconcile existing public DNS records

- **WHEN** the project plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
  and does not require a RouterOS provider.

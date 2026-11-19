## MODIFIED Requirements

### Requirement: Reusable execution and repository integration

The reusable module SHALL own Terraform execution and named command maps while
remaining independent of parent repository labels. Maps SHALL support optional
caller-supplied command wrappers. Repository consumers SHALL select generic AL
wrappers explicitly to preserve configuration, plugin lifecycle, backend
injection, named operations, and the saved-plan apply guard. Shared provider
pins SHALL be owned by third_party/terraform.

#### Scenario: Saved-plan apply

- **WHEN** a targeted apply wrapper receives a missing plan or extra arguments
- **THEN** the runner rejects it before Terraform initialization.

#### Scenario: Repository migration

- **WHEN** repository Terraform commands and tests are analyzed
- **THEN** they load maps from rules_terraform, select shared provider labels,
  and require neither tools/terraform nor a checked-in .terraform.lock.hcl file.

#### Scenario: Optional wrapper

- **WHEN** a command map supplies a wrapper
- **THEN** the wrapper receives the declared Terraform invocation and runfiles
  without requiring a framework dependency in the reusable module.

#### Scenario: Direct tests

- **WHEN** a test map needs no external command setup
- **THEN** it runs the Terraform test rule without an AL dependency.

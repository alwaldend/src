## Purpose

Maintain repository dependency boundaries through central Gazelle policy
configuration and Bazel analysis rather than repeated manual package edits.

## ADDED Requirements

### Requirement: Explicit path-based policy generation

Gazelle SHALL read `tools/transitive_visibility/transitive_visibility.toml` in
its workspace, or an explicitly selected configuration file. The `[paths]`
table SHALL map canonical workspace-relative directory paths to absolute
package-group labels. The most specific directory-prefix match SHALL apply.
Gazelle MUST preserve unrelated attributes and comments and MUST NOT invent
policy for unmatched paths or create otherwise empty packages.

#### Scenario: Descendant package inherits a tree policy

- **WHEN** a path has a mapping and a descendant has no more specific mapping
- **THEN** Gazelle applies that path's group to the descendant package

#### Scenario: Package ordering and repeat generation

- **WHEN** Gazelle updates a managed package declaration in an existing file
- **THEN** it appears after loads and before rules
- **AND** another run with unchanged configuration and sources changes nothing

#### Scenario: Newly created BUILD file

- **WHEN** upstream Gazelle creates a BUILD file containing a package declaration
- **THEN** its first-run ordering is retained without upstream patches
- **AND** a later run moves the declaration after loads and before rules

#### Scenario: Invalid configuration

- **WHEN** TOML, paths, labels, or unknown fields are invalid, or an explicitly
  selected configuration is missing
- **THEN** Gazelle fails before writing BUILD files

#### Scenario: Independent workspace

- **WHEN** Gazelle runs in a nested workspace without its own configuration
- **THEN** it neither reads the parent configuration nor changes manual policy

### Requirement: Native dependency enforcement

Checked-in repository Bazel configuration SHALL enable experimental transitive
visibility enforcement. Each repository tree SHALL own a package group selected by a
central path mapping. First-party module compatibility SHALL match the repository's
Bazel 9.3.1 baseline.

#### Scenario: Forbidden direct artifact dependency

- **WHEN** a project build depends directly on an infrastructure artifact
  excluded by its package's consumer policy
- **THEN** Bazel analysis rejects that dependency

#### Scenario: Native indirect-enforcement limitation

- **WHEN** an intermediary declares its own policy in Bazel 9.3.1
- **THEN** documentation records that inherited restrictions can be masked,
  without claiming complete transitive enforcement

#### Scenario: No host configuration mutation

- **WHEN** this feature is enabled
- **THEN** enforcement is configured in tracked repository files and requires
  no changes to host-global Bazel configuration

### Requirement: Data-owned public trust material

Public CA certificates and SSH CA public keys consumed by project builds SHALL
be checked-in data assets. Their data targets MUST NOT depend on infra targets.
Each owning data directory SHALL have a separate Terraform refresh root that
reads public values through Vault data sources without reading infrastructure Terraform state.
Infrastructure Terraform MUST NOT write those data assets.

#### Scenario: Reusable collection build

- **WHEN** the Ansible collection is built
- **THEN** it uses the data snapshots without fetching live material or
  depending on infra targets

#### Scenario: Explicit snapshot refresh

- **WHEN** an operator applies one data-owned Terraform root
- **THEN** that root updates only its public snapshot files, retaining state
  outside tracked source and selecting only the public certificate or key field

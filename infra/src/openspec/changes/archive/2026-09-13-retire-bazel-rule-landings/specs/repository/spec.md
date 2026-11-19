## ADDED Requirements

### Requirement: Bazel rule modules are tool components

First-party standalone Bazel rule modules SHALL live under `tools/<module>/`.
Their relocation MUST preserve module names, supported public rule and macro
interfaces, standalone workspace resolution, generated rule documentation,
and owner-local specifications. Repository consumers and discovery mechanisms
MUST resolve their new source locations through the owning declarations.

#### Scenario: Consume a relocated rule module

- **WHEN** a root or nested workspace loads a relocated rule module through its
  supported module name and public rule interface
- **THEN** the module resolves from its new tool-owned location
- **AND** its supported build behavior remains available without changing the
  public module identity

#### Scenario: Discover a relocated owner's specifications and rules

- **WHEN** repository documentation, skill discovery, or OpenSpec validation
  traverses the relocated modules
- **THEN** it selects the modules' new owner paths
- **AND** generated rule documentation and existing component specifications
  remain included in their declared outputs and checks

### Requirement: Rule tools do not publish dedicated project landings

The relocated rule modules SHALL have no dedicated project landing-site
configuration. After authorized retirement, their selected landing repository
and DNS declarations MUST be absent from the owning infrastructure, and
shared-catalog consumers MUST exclude those retired repositories. Retirement
MUST preserve unrelated repositories, sites, access grants, and authentication
resources.

#### Scenario: Build or import the remaining repository catalog

- **WHEN** the catalog is projected for GitHub or GitLab after rule-landing
  retirement
- **THEN** the retired landing repositories are absent
- **AND** remaining repositories retain their catalog identities, names, and
  configured access behavior

#### Scenario: Retire a selected landing site

- **WHEN** the owning declarative workflows execute the user's explicitly
  authorized retirement scope
- **THEN** the selected Pages repository and matching DNS record are removed
- **AND** the reviewed plans and postconditions show no unrelated deletion or
  replacement

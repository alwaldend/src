## Purpose

Define shared organization, repository identity, naming, and named access
configuration for the GitHub, GitLab, and Forgejo infrastructure consumers.

## ADDED Requirements

### Requirement: Authoritative organization and repository inventory

The repository catalog SHALL store shared defaults and schema version at its
root, one configuration file per organization, and one configuration file per
repository beneath its organization. Those records SHALL own the configured
organizations, repositories, named administrators, and named developers.
Each forge consumer SHALL use
that catalog for its selected repositories and named roles, preserving
forge-specific settings without maintaining another copy of the inventory.
Authentication and service-specific grants SHALL remain with their existing
owners.

#### Scenario: Consume one organization definition

- **WHEN** an organization assigns named administrators and developers and
  selects repositories for multiple forges
- **THEN** each selected consumer uses the same catalog assignments and
  repository identities
- **AND** forge-specific settings and service identities remain with their
  declared owners

### Requirement: Stable names follow ownership and original upstream

First-party repository names SHALL remain the same across forges, including
when copied or synchronized. External forks and mirrors SHALL belong to a
catalog organization and use a name derived from the original upstream's
reversed hostname followed by its full owner and repository path, in
lowercase with punctuation replaced by underscores. Cross-forge copies of
external repositories SHALL retain that name rather than naming the
intermediate copy as their upstream.

#### Scenario: Copy an organization-owned repository

- **WHEN** the first-party `alwaldend/src` repository is copied from GitHub to
  GitLab or Forgejo
- **THEN** its destination remains `alwaldend/src`

#### Scenario: Fork an external repository

- **WHEN** `https://gitlab.com/fdroid/fdroiddata` is forked into `alwaldend`
- **THEN** its organization-owned name is
  `alwaldend/com_gitlab_fdroid_fdroiddata`
- **AND** subsequent copies on another forge retain that name

### Requirement: Named developer access excludes protected default writes

Consumers SHALL grant catalog administrators administrative access and
catalog developers repository reads, feature-branch writes, and pull or
merge request creation. Catalog developers SHALL NOT receive direct push
or merge access to protected default branches through their named roles or
overlapping retained grants. Existing service-specific permissions SHALL be
preserved and distinguished from named developer assignments.

#### Scenario: Develop on a feature branch

- **WHEN** a catalog developer accesses a selected repository
- **THEN** the developer can read its contents, create a feature branch, and
  open a pull or merge request
- **AND** the developer cannot push or merge to its protected default branch

#### Scenario: Deploy a landing repository

- **WHEN** a GitHub landing repository is adopted
- **THEN** its protected default branch is `master`
- **AND** its Pages source remains `pages` with its existing custom domain
- **AND** the developer can continue deploying to `pages` without receiving
  push or merge access to `master`

### Requirement: Adopt existing resources without recreation

Consumers SHALL adopt existing remote resources through imports and preserve
their remote identities. Existing managed resource addresses SHALL remain
stable or use explicit state-preserving moves. Existing GitHub repository
settings and Pages configuration SHALL remain intact. Default branches SHALL
remain unchanged except for the requested landing-repository migration to
`master`; a missing destination branch SHALL be created from existing
repository contents before changing the default. No deletion or replacement
SHALL execute without separate user approval for its concrete scope.

#### Scenario: Adopt an existing repository

- **WHEN** a catalog repository already exists remotely or in Terraform state
- **THEN** its import or state move preserves the existing remote repository
  ID and contents
- **AND** the reviewed plan does not recreate it to resolve a name collision

#### Scenario: Encounter a destructive plan

- **WHEN** any reviewed plan proposes a deletion or replacement
- **THEN** execution of that operation stops pending the user's explicit
  approval for the identified operation and scope

#### Scenario: Change a landing default without changing Pages

- **WHEN** a landing repository has no `master` branch before adoption
- **THEN** the migration creates it from existing repository contents before
  selecting it as the protected default branch
- **AND** the existing `pages` branch and Pages configuration remain intact

### Requirement: GitLab copies and fork preserve their source identity

GitLab SHALL receive a one-time Git import of every catalog-selected GitHub
repository and an organization-owned fork of the selected F-Droid metadata
upstream. Imports SHALL preserve the source default branch, and the fork
SHALL retain its upstream fork relationship. This change SHALL configure no
ongoing synchronization.

#### Scenario: Complete the initial GitLab population

- **WHEN** the GitHub default-branch migration completes and the initial
  GitLab imports and fork finish
- **THEN** every selected GitHub repository has its intended GitLab copy and
  source default branch
- **AND** the metadata repository retains its F-Droid upstream fork identity
- **AND** later upstream changes are not promised to synchronize

### Requirement: Adopt GitLab default-branch protections without removal

GitLab default-branch access SHALL be established through group defaults
before new repositories are populated. Existing automatically created
protections SHALL be verified and imported before standalone protection
management is applied. Adoption SHALL preserve the existing protections;
provider limitations SHALL be reported rather than bypassed by unprotecting
or recreating them.

#### Scenario: Adopt an automatically protected default branch

- **WHEN** an imported GitLab repository has an existing default-branch rule
- **THEN** Terraform adopts that rule by import before applying permitted
  in-place access changes
- **AND** adoption does not invoke unprotect-and-recreate behavior

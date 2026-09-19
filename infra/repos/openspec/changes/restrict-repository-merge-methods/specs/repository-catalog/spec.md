## ADDED Requirements

### Requirement: Catalog GitHub repositories permit merge commits only

The catalog SHALL permit merge commits on its GitHub repositories and SHALL
disable squash merges and rebase merges. The shared defaults SHALL express
this policy, and each GitHub repository record SHALL inherit or restate it
without contradiction. A repository record SHALL NOT enable squash or rebase
merges. This policy SHALL NOT change GitLab or Forgejo merge settings, which
remain outside this catalog's ownership.

#### Scenario: Apply the shared merge policy

- **WHEN** the catalog projects its GitHub repositories for a consumer
- **THEN** every projected repository permits merge commits
- **AND** every projected repository disables squash merges and rebase merges

#### Scenario: Reject a contradicting repository record

- **WHEN** a GitHub repository record enables squash merges or rebase merges
- **THEN** catalog validation fails before any consumer plan is produced
- **AND** the reported error identifies the repository and the disallowed method

#### Scenario: Preserve the reviewed commit through a merge

- **WHEN** a reviewed feature branch is accepted into a protected default branch
- **THEN** the default branch records a merge commit whose parents include the
  reviewed feature commit
- **AND** the reviewed feature commit remains reachable and unmodified

## MODIFIED Requirements

### Requirement: Adopt existing resources without recreation

Consumers SHALL adopt existing remote resources through imports and preserve
their remote identities. Existing managed resource addresses SHALL remain
stable or use explicit state-preserving moves. Existing GitHub repository
settings and Pages configuration SHALL remain intact, except for the
catalog-owned GitHub merge-method settings, which the catalog owns and may
change in place on existing repositories. Default branches SHALL remain
unchanged. No deletion or replacement SHALL execute without separate user
approval for its concrete scope.

#### Scenario: Adopt an existing repository

- **WHEN** a catalog repository already exists remotely or in Terraform state
- **THEN** its import or state move preserves the existing remote repository
  ID and contents
- **AND** the reviewed plan does not recreate it to resolve a name collision

#### Scenario: Encounter a destructive plan

- **WHEN** any reviewed plan proposes a deletion or replacement
- **THEN** execution of that operation stops pending the user's explicit
  approval for the identified operation and scope

#### Scenario: Apply owned merge settings to an existing repository

- **WHEN** an adopted repository's merge settings differ from the catalog's policy
- **THEN** the reviewed plan updates those settings in place without deleting
  or replacing the repository
- **AND** settings the catalog does not own remain unchanged

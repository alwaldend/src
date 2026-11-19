## ADDED Requirements

### Requirement: Attribute signed pull request merges to the merging account

Forgejo SHALL use the collaborator trust model by default and sign pull
request merges with the instance SSH signing key.

#### Scenario: Merge a pull request in a repository that inherits the default

- **WHEN** a user merges a pull request and Forgejo creates a merge commit
- **THEN** the commit records the merging account as its Git committer
- **AND** the commit carries a signature from the instance SSH signing key

## ADDED Requirements

### Requirement: Repository CI generation and lint checks

Repository CI SHALL check Gazelle-generated BUILD files in the root workspace for drift and run
the semantic lint profile. The Forgejo summary SHALL show the lint result.

#### Scenario: A generation or lint check fails

- **WHEN** Gazelle finds BUILD file drift or the semantic lint profile fails
- **THEN** repository CI reports failure and preserves the diagnostics
- **AND** CI attempts its other phases

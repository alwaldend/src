## Purpose

Define how trusted Forgejo jobs check out the event revision and run shared
repository CI, including reuse of a stable Bazel workspace.

## ADDED Requirements

### Requirement: Optional stable CI checkout

The repository CI action SHALL default to the job checkout. When explicitly
enabled, it SHALL prepare the event revision in a stable directory owned by the
serial runner, remove prior job files, and run CI from that directory.

#### Scenario: Two serial jobs build different revisions

- **WHEN** two serial jobs enable the stable checkout
- **THEN** both run CI from the same directory and each sees its own revision,
  repository history, tags, and materialized LFS files
- **AND** untracked and ignored leftovers from the preceding job are removed

### Requirement: Reusable CI workflow

The repository SHALL expose its checkout and CI steps as a callable workflow.

#### Scenario: A trusted workflow calls repository CI

- **WHEN** a trusted workflow in the repository calls the shared workflow
- **THEN** it checks out the event revision and invokes the same CI action

## Purpose

Allow agents to validate and publish a prepared candidate without incidental environment drift and retain useful evidence when final verification is interrupted.

## ADDED Requirements

### Requirement: Combined validation and publication

The validation command SHALL accept explicit publication in the same process, preserving receipt, candidate, environment, log, and remote guards. Failed validation MUST NOT publish. Git-only publication MUST require explicit publication.

#### Scenario: Publish passing checks

- **WHEN** the caller runs `validate --publish` with a valid receipt and passing plan
- **THEN** the command validates and continues publication in one inherited environment

#### Scenario: Validation fails

- **WHEN** a selected check fails
- **THEN** the command reports the failed validation and performs no push or PR mutation

### Requirement: Retain publication evidence

When final verification fails after push and PR synchronization, the command SHALL retain the known head, tree, remote head, and PR identity in a structured report. It MUST mark verification incomplete and preserve verification-only continuation without repeating publication.

#### Scenario: Concurrent review update

- **WHEN** review inventory changes during final verification
- **THEN** the command reports `publication_verification_pending`, the known PR and candidate, and an error
- **AND** continuation verifies without pushing or rerunning checks

## Purpose

Provide reproducible x go sdk through explicit Bazel targets with declared inputs and outputs.

## ADDED Requirements

### Requirement: X Go SDK

The schema dependency SHALL own the authenticated schema update command. The library SHALL provide a Go client generated from the checksum-pinned mirrored X schema as a Bazel build artifact. Generated client source MUST NOT be committed. Failed, redirected, oversized or invalid schema responses MUST preserve the previous downloaded snapshot and MUST NOT disclose credentials.

#### Scenario: Successful generation

- **WHEN** a developer generates from the declared schema input
- **THEN** the pinned generator produces a compilable Go client without network access or credentials

#### Scenario: Partial-success response

- **WHEN** an HTTP 200 response includes successful user data and a Problem
- **THEN** the client retains the user data and selects the problem variant through its schema discriminator

#### Scenario: Schema update rejected

- **WHEN** the schema endpoint rejects authentication or returns an invalid response
- **THEN** the command fails without modifying the declared schema input or printing secrets

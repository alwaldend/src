## Purpose

Allow the Forgejo role to deploy a caller-supplied header through its controlled
service activation workflow.

## ADDED Requirements

### Requirement: Activate optional Forgejo header with startup inputs

When configured with a custom header source, the role SHALL compare it by
content, stage it before downtime, and install it under Forgejo's custom
template directory before starting the service. For existing installations the
role SHALL archive the current managed startup inputs before replacement.

#### Scenario: Add a header to an existing installation

- **WHEN** a source is supplied and the installed header is missing or different
- **THEN** the role activates the staged header through the guarded workflow
- **AND** an unchanged subsequent run does not restart the service

#### Scenario: Run without a supplied header

- **WHEN** no header source is configured
- **THEN** the role leaves the custom header path unmanaged

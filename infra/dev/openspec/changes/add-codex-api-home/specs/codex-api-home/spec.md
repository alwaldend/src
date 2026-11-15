## Purpose

Provide an isolated Codex home for API-key authentication while preserving existing homes.

## ADDED Requirements

### Requirement: Separate API home

Host Bot SHALL manage `~/.codex-api` with owner-only directory and file permissions, shared Codex settings and instructions, and API-only login. It SHALL NOT copy existing credentials or conversation state.

#### Scenario: Select API home

- **WHEN** the managed configuration is deployed and the user selects `~/.codex-api`
- **THEN** Codex uses API login and separate state with the managed settings and instructions
- **AND** the default home remains unchanged

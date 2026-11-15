# openspec-cli Specification

## Purpose

Keep Host Bot's user-facing OpenSpec CLI installation reproducible through its checked-in Ansible deployment.

## Requirements

### Requirement: Managed OpenSpec CLI installation

Host Bot's Ansible deployment SHALL install `@fission-ai/openspec` at an exact
inventory-configured version under the managed user's `~/.local` prefix, with
Node.js and npm provisioned through the existing host package backend.
Repeated deployment with the same version SHALL preserve the installed version.

#### Scenario: Provision the CLI

- **WHEN** the Host Bot deployment runs with OpenSpec absent
- **THEN** it installs the configured version and exposes `~/.local/bin/openspec` to the managed user

#### Scenario: Repeat the deployment

- **WHEN** the configured OpenSpec version is already installed
- **THEN** the npm task leaves the package unchanged

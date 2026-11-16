## Purpose

Make the repository-pinned Forgejo CLI available to the development host's
managed account through reproducible Ansible provisioning.

## ADDED Requirements

### Requirement: Install the pinned CLI for the managed account

The development deployment SHALL install the repository-pinned fj executable
as `.local/bin/fj` in the managed user's home, with executable permissions.
The installation SHALL be selectable without unrelated host changes. A local
entry point SHALL execute the same packaged installation tasks without requiring
Vault or SSH credentials. Reapplying unchanged source SHALL report no changes.

#### Scenario: Requested local installation

- **WHEN** the authorized operator runs the local fj installation entry point on dev1
- **THEN** only the managed user's CLI directory and executable are reconciled
- **AND** the installed command reports the pinned version

#### Scenario: Repeat installation

- **WHEN** the same installation runs again with unchanged pinned input
- **THEN** it reports no changes

### Requirement: Keep login separate from provisioning

Ansible SHALL NOT provision a Forgejo password, token or OAuth session. Native
authentication SHALL be initiated only under explicit login authority. If the
current host lacks a configured OAuth client, token import SHALL use stdin and
SHALL NOT expose the credential in command arguments or output.

#### Scenario: Browser login requested

- **WHEN** the user asks to start login after installation and an instance OAuth client is configured
- **THEN** native fj authentication opens its browser flow
- **AND** no password is included in repository source or command arguments

#### Scenario: Token fallback selected

- **WHEN** browser login reports a missing instance client and the user supplies a token privately
- **THEN** native fj imports the token through stdin and verifies the expected account
- **AND** the one-use credential handoff is deleted without displaying its contents

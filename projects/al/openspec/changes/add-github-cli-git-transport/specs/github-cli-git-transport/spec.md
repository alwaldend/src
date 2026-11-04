## Purpose

Allow repository operators to use GitHub CLI credentials for guarded Git publication when SSH authentication is unavailable.

## ADDED Requirements

### Requirement: Explicit GitHub CLI transport

Delivery SHALL support explicit `--git-transport gh` for canonical GitHub SSH and HTTPS remotes without editing persistent Git configuration. SSH SHALL remain the default and SHALL NOT automatically fall back.

#### Scenario: Select gh with an SSH remote

- **WHEN** an authenticated operator selects gh transport for a canonical GitHub SSH remote
- **THEN** Git fetch, push and verification use HTTPS with gh credentials while retaining the captured repository identity and publication guards.

#### Scenario: Select gh with an HTTPS remote

- **WHEN** an authenticated operator selects gh transport for a canonical GitHub HTTPS remote
- **THEN** delivery accepts that endpoint without importing global credential helpers.

### Requirement: Isolated credentials and transport binding

The transport MUST use only the selected GitHub CLI credential helper for the captured GitHub host, verify TLS, reject redirects and unsupported endpoints, disable prompts, and keep credentials out of AL arguments and output. Preparation receipts MUST bind transport mode; legacy receipts MUST imply SSH.

#### Scenario: Change transport after preparation

- **WHEN** publication selects a different transport from the preparation receipt
- **THEN** delivery refuses publication.

#### Scenario: Authentication fails

- **WHEN** gh cannot authenticate the captured GitHub remote
- **THEN** delivery fails without falling back, publishing, or exposing credentials.

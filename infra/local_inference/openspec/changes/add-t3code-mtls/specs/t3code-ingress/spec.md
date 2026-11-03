## Purpose

Expose the local inference host's T3 Code interface through authenticated
HTTPS while retaining local Codex inference and local execution boundaries.

## ADDED Requirements

### Requirement: Client-certificate-protected browser access

The component SHALL declare a DNS hostname and an HTTPS route to T3 Code.
Traefik MUST require and verify client certificates against the repository's
client CA. The application and inference ports MUST remain loopback-only.

#### Scenario: Connect without a trusted client certificate

- **WHEN** a browser connects without a valid trusted client certificate
- **THEN** the TLS handshake is rejected before reaching T3 Code

### Requirement: Local Codex integration

T3 Code SHALL use the managed local Codex executable and local model by
default, with Auto approvals through the same local provider and Full access
available. It MUST run normally as a dedicated service account, without
enabling hosted Codex tools or external inference. Installation MUST NOT require
package downloads during ordinary service startup.

#### Scenario: Start a coding session

- **WHEN** T3 Code launches Codex for a local workspace
- **THEN** managed requirements continue to pin loopback inference and
  reject OpenAI-hosted tools while local command execution and T3's
  authenticated loopback MCP remain available

#### Scenario: Redirect the approved MCP identity

- **WHEN** a workspace changes the `t3-code` MCP URL or adds another server
- **THEN** managed requirements reject the unapproved identity

#### Scenario: Invoke a namespaced local tool

- **WHEN** Codex advertises T3 or subagent functions inside a Responses namespace
- **THEN** the loopback bridge presents every function to llama.cpp, restores
  its original namespace in streamed and JSON replies, and preserves call
  identity when Codex submits the tool result on the following request

#### Scenario: Select Full access

- **WHEN** the operator selects Full access instead of Auto
- **THEN** managed policy permits it without changing the local model provider

### Requirement: Declarative secrets and ingress wiring

The component SHALL use T3 Code's native pairing/session authentication
and a component-scoped Vault PKI enrollment role. Live secret values MUST
NOT appear in checked-in configuration or validation fixtures.

#### Scenario: Review the prepared deployment

- **WHEN** an operator inspects the IaC
- **THEN** certificate enrollment uses declared Vault references, native
  application authentication is retained, and ingress agrees with DNS

### Requirement: Offline validation does not deploy

Validation MUST exercise rendered configuration and package checks without
provisioning cloud resources or enrolling live certificates.

#### Scenario: Validate the prepared extension

- **WHEN** the source checks run
- **THEN** no Terraform apply, inventory-host connection, or Vault write runs

## Purpose

Install an optional instance SSH commit-signing identity before Forgejo starts
while preserving its long-lived private key across repeated deployments.

## ADDED Requirements

### Requirement: Guard optional signing-file installation

The role SHALL accept an optional matching SSH pair, install it under its
backed-up persistence root before activation, and restrict private access to the
service account. An existing different private key SHALL cause failure before
replacement or service downtime. Tasks carrying secrets SHALL suppress output.

#### Scenario: Repeated installation

- **WHEN** the supplied pair matches the installed private key
- **THEN** the identity remains unchanged and file metadata and public content
  can be repaired without restarting the service

#### Scenario: Reject a different private key

- **WHEN** the supplied private key differs from the installed one
- **THEN** deployment fails without replacing the key or stopping Forgejo

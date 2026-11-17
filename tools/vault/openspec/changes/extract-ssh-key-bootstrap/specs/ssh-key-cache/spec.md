## Purpose

Provide a reusable repository tool that creates a durable SSH key pair in Vault
without coupling its lifecycle to a particular service or exposing private data.

## ADDED Requirements

### Requirement: Consumer-neutral immutable SSH key cache

The tool SHALL accept a Vault KV v2 mount and path through neutral inputs,
generate an Ed25519 pair only for an unused record, and reuse a complete matching
pair. Existing history, corrupt data, access errors, or conflicting writes SHALL
fail without overwriting the identity. Private values SHALL remain outside
arguments, output, local files, and Terraform attributes.

#### Scenario: Reuse a key for an independent consumer

- **WHEN** a caller selects a mount and path with a valid cached pair
- **THEN** the command reuses that pair without service-specific inputs

#### Scenario: Preserve an existing identity

- **WHEN** the record is corrupt, inaccessible, deleted, or concurrently created
- **THEN** the command fails instead of overwriting it

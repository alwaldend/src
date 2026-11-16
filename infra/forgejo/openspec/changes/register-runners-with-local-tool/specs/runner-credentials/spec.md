## Purpose

Bootstrap repository-scoped Forgejo runner credentials through the owning
Terraform stage while retaining secret values only in Vault.

## ADDED Requirements

### Requirement: Local apply-time credential bootstrap

The service Terraform SHALL invoke a packaged local tool during runner resource
creation. It SHALL use the invocation's existing Forgejo and Vault credentials,
store runner credentials directly in Vault KV v2, and emit only a non-secret
completion report. Ordinary plans SHALL NOT register runners through this tool.

#### Scenario: Register a new runner

- **WHEN** the declared runner and its Vault secret do not exist
- **THEN** an authorized apply creates one repository-scoped runner and stores
  its identity and confidential token in the declared Vault path
- **AND** no credential is printed or returned as a Terraform attribute

### Requirement: Preserve registration ownership on retries

Terraform SHALL retain a stable creation-only bootstrap resource without
replacement triggers. The tool SHALL reuse a complete matching Vault record
without another Forgejo request and reject conflicting or inaccessible records.
Execution SHALL be serialized through the owning Terraform backend lock.
When credentials are absent, an existing same-name runner SHALL prevent
additional creation or automatic deletion.

#### Scenario: Repeat a completed bootstrap

- **WHEN** a complete matching Vault record exists
- **THEN** an explicit command repeat reuses it without a Forgejo request or Vault write

#### Scenario: Preserve the creation-only resource

- **WHEN** an ordinary plan or apply runs after registration, including after
  tool source changes
- **THEN** Terraform does not rerun the tool or invalidate its credentials

#### Scenario: Interrupt registration

- **WHEN** creation succeeded but its response or credential persistence failed
- **THEN** a later invocation detects the same-name runner and stops for
  reconciliation without creating duplicates

#### Scenario: Reject an existing unmanaged runner

- **WHEN** a same-name runner is found on any runner-list page without a matching
  completed Vault record
- **THEN** the command fails without adopting, deleting, or recreating that runner

# project-dns Specification

## Purpose

Give `rules_template` an independent Terraform state for its existing DNS declarations
while retaining the project's reusable Bazel module boundary.

## Requirements

### Requirement: Own project DNS through one Terraform root

The project's `tf` root SHALL consume its canonical `dnsconfig.json` and use a
project-specific state and Vault/AppRole configuration. Record ownership and
adoption SHALL follow the
[shared migration contract](../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/specs/infra-dns/spec.md).
After authorized adoption into the project's state, the root SHALL retain
enabled record ownership in its checked-in source defaults.

#### Scenario: The prepared root has not been adopted

- **WHEN** the project root is packaged with its initial configuration
- **THEN** DNS record creation is disabled pending authorized adoption.

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted project root uses its checked-in source defaults
- **THEN** `dns_enabled` is true and record ownership remains enabled.

#### Scenario: Reconcile existing public DNS records

- **WHEN** the project plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions.

### Requirement: Use the normal authenticated Terraform flow

The project SHALL select DNS provider credential injection with `tf=main`.
Its AppRole and required Vault credentials SHALL be deployed before its
authenticated root commands run. The `dns_enabled` setting SHALL control DNS
record creation without bypassing authentication or credential injection.

#### Scenario: A root command runs before DNS adoption

- **WHEN** an authenticated Terraform root command runs with `dns_enabled=false`
- **THEN** the normal AppRole, backend, and provider credential flow still applies.

### Requirement: Keep reusable module consumers independent of operations

The project SHALL export canonical Terraform source for the repository's
operational wrapper without requiring consumers of its reusable Bazel rules
to depend on the parent monorepo.

#### Scenario: The parent workspace packages project DNS operations

- **WHEN** the root-workspace operational wrapper is built
- **THEN** it consumes the project's exported Terraform source and DNS declarations without copying their definitions.

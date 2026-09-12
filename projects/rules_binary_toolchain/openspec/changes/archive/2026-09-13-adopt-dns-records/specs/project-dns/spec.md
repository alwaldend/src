## MODIFIED Requirements

### Requirement: Own project DNS through one Terraform root

The project's `tf` root SHALL consume its canonical `dnsconfig.json` and use a
project-specific state and Vault/AppRole configuration. Record ownership and
adoption SHALL follow the
[shared migration contract](../../../../../../../../infra/dns/openspec/changes/archive/2026-09-13-migrate-project-dns-to-terraform/specs/infra-dns/spec.md).

After authorized adoption, the root SHALL retain enabled record ownership in
its checked-in defaults and reconcile unchanged declarations without record
additions, changes, replacements, or deletions.

#### Scenario: The prepared root has not been adopted

- **WHEN** the project root is packaged with its initial configuration
- **THEN** DNS record creation is disabled pending authorized adoption.

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted project root uses its checked-in source defaults
- **THEN** DNS record ownership is enabled in its project-specific state.

#### Scenario: Reconcile existing public DNS records

- **WHEN** the project plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions.

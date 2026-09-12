# owned-dns Specification

## Purpose

Define Host Bot DNS management through the owner's `tf` root,
including canonical declarations, scoped credentials, and offline source checks
before adopting live records.

## Requirements

### Requirement: Owner-local DNS configuration

The `tf` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled until the shared migration adopts them
and SHALL retain enabled ownership after authorized adoption into the owner's
state. Reconciliation of adopted records with unchanged declarations SHALL
propose no record additions, changes, replacements, or deletions.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before authorized
  adoption
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted host root is evaluated with its checked-in source defaults
- **THEN** record ownership is enabled for the owner's declared records and views
- **AND** the package retains the canonical declaration and shared module inputs

#### Scenario: Reconcile existing DNS records

- **WHEN** the host root plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
- **AND** records outside this owner's declarations remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrapper SHALL select `src_users_simeonwarren_host_bot` through the repository AL flow and keep
secret values in injected variables. DNS injection SHALL select the root's existing
Terraform stage label. Real DNS credentials and Vault policy grants SHALL be
prerequisites for operational Terraform calls. RouterOS DNS credentials SHALL be isolated
from unrelated RouterOS resources. The package SHALL expose a format test that
does not authenticate to Vault or contact DNS providers.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** the wrapper selects its Terraform stage labels
- **THEN** it selects DNS credential injection through the owner AppRole
- **AND** disabled DNS resources do not imply offline provider configuration

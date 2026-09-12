# owned-dns Specification

## Purpose

Define Vault service DNS management through the owner's `tf_setup` root,
including canonical declarations, scoped credentials, and offline source checks
before adopting live records.

## Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled until the authorized adoption revision
and SHALL retain enabled ownership after existing records are imported into the
owner's state. Reconciliation against unchanged declarations and adopted state
SHALL propose no record additions, changes, replacements, or deletions.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before its
  authorized adoption revision
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Adopt existing DNS records

- **WHEN** the authorized adoption completes the shared writer audit and
  applicable recovery prerequisites and enables the owner's record management
- **THEN** the owning root imports existing records by exact provider ID into
  its state while preserving their declared identities and views

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted root uses its checked-in source defaults
- **THEN** DNS ownership remains enabled and the package retains the module
  and canonical declaration inputs

#### Scenario: Reconcile existing DNS records

- **WHEN** the owning root plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions

### Requirement: Scoped DNS execution and offline checks

The `dns.plan`, `dns.show`, and `dns.apply` entrypoints SHALL select
`src_infra_dc1_vault` through the repository AL flow using `dns=1` and operate
on `module.dns` in the existing `tf_setup` root and Vault HTTP backend. The
apply entrypoint SHALL require a saved plan. Ordinary setup and service
wrappers SHALL retain their stage labels, authentication, and backend paths.
Secret values SHALL remain in injected variables. Real DNS credentials and
Vault policy grants SHALL be prerequisites for operational DNS calls.
RouterOS DNS credentials SHALL remain isolated from unrelated RouterOS
resources. The package SHALL expose a format test that does not authenticate
to Vault or contact DNS providers.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** a scoped DNS entrypoint selects `dns=1`
- **THEN** it selects DNS credential injection through the owner AppRole and
  existing setup backend, with plan and show scoped to `module.dns` and apply
  requiring a saved plan
- **AND** disabled DNS resources do not imply offline provider configuration
- **AND** ordinary setup and service authentication retains its original labels

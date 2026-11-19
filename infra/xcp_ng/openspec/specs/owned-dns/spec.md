# owned-dns Specification

## Purpose

Define XCP-ng infrastructure DNS management through the owner's `tf` root,
including canonical declarations, scoped execution, and offline source checks
through preparation, adoption, and ongoing reconciliation.

## Requirements

### Requirement: Owner-local DNS configuration

The `tf` root SHALL consume this owner's canonical `dnsconfig.json` through
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

- **WHEN** the authorized adoption audits competing writers, controls those
  identified within its observed scope, and enables the owner's record management
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

The DNS-only plan, show, and saved-plan apply commands SHALL select `dns=1`
through the repository AL flow, using `src_infra_xcp_ng` through its named
`xcp_ng` Vault authentication. They SHALL target `module.dns` in the existing
`tf` root and backend and keep secret values in injected variables. The ordinary
service wrappers SHALL retain their existing `xcp_ng_tf`, `xoa`, and Vault
environment labels and authentication behavior. Real DNS credentials and Vault
policy grants SHALL be prerequisites for operational Terraform calls. RouterOS
DNS credentials SHALL be isolated from unrelated RouterOS resources. The package
SHALL expose a format test that does not authenticate to Vault or contact DNS
providers. Successful DNS-only execution SHALL establish only the scoped DNS
result, not full service health.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** the DNS wrapper selects `dns=1`
- **THEN** it selects DNS credential injection through the owner AppRole while
  preserving the existing root and backend without starting XO authentication
- **AND** disabled DNS resources do not imply offline provider configuration

#### Scenario: Review and apply only DNS resources

- **WHEN** the operator reviews and applies a saved plan through the DNS-only
  plan, show, and apply commands
- **THEN** the operation is limited to `module.dns` and its dependencies
- **AND** its success does not establish the health of the other services in the root

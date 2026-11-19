## MODIFIED Requirements

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

- **WHEN** authorized adoption records the shared writer audit and control
  coverage and enables the owner's record management
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

The DNS wrappers SHALL select `src_infra_ingress` through the repository AL flow
and keep secret values in injected variables. `dns.plan`, `dns.show`, and
`dns.apply` SHALL retain the owning root and backend, select `dns=1`, and scope
planning to `module.dns`. Apply SHALL require the reviewed saved plan. This
workflow SHALL preserve ordinary service authentication and SHALL NOT establish
overall service health. Required DNS credentials and Vault grants SHALL remain
operational prerequisites, and RouterOS DNS credentials SHALL remain isolated
from unrelated RouterOS resources. The package SHALL expose a format test
without Vault authentication or DNS provider access.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** the operator selects the owning DNS wrappers
- **THEN** they select DNS credential injection through the owner AppRole with `dns=1`
- **AND** their scoped plan validates DNS and its dependencies
- **AND** disabled DNS resources do not imply offline provider configuration

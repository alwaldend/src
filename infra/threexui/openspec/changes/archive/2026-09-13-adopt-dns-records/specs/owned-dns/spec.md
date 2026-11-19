## MODIFIED Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled before the adoption revision and SHALL
retain enabled ownership after authorized adoption into the owner's state.
Reconciliation against adopted state and unchanged declarations SHALL propose
no record additions, changes, replacements, or deletions.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the preparatory root is evaluated with default inputs before its
  adoption revision
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted root is evaluated with its checked-in default inputs
- **THEN** DNS ownership is enabled and the root reads the canonical declaration
- **AND** the package contains the module and declaration inputs

#### Scenario: Adopt existing DNS records

- **WHEN** authorized adoption coordinates the previous writers and enables
  this owner's record ownership
- **THEN** the owner imports existing provider identities into its own state
  through the owning scoped DNS wrappers and reviewed saved-plan import blocks
- **AND** existing record values and declared DNS views are preserved

#### Scenario: Reconcile adopted DNS records

- **WHEN** the owner plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
- **AND** unrelated records remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrappers SHALL select `src_infra_threexui` through the repository AL flow
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

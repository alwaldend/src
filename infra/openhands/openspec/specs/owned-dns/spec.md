# owned-dns Specification

## Purpose

Define OpenHands DNS management through the owner's `tf_setup` root,
including canonical declarations, scoped credentials, and offline source checks
before adopting live records.

## Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled before adoption and SHALL default to
enabled in the reviewed adoption revision. Adoption SHALL bind existing records
to their exact provider identities and SHALL preserve their declared attributes.
Existing-record adoption SHALL contain no DNS record additions, changes,
replacements, or deletions. Verified missing declarations MAY be provisioned
through an exact additions-only DNS plan that preserves every existing record
and excludes non-DNS managed-resource changes. After deployment, default DNS
reconciliation SHALL retain the bindings and preserve unrelated provider records.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before adoption
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect the adopted default

- **WHEN** the reviewed adoption revision is evaluated with default inputs after adoption
- **THEN** it reads this owner's declaration and enables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Adopt existing records

- **WHEN** the owner's existing records are adopted through its scoped DNS wrappers
- **THEN** Terraform binds the records to their exact provider identities
- **AND** their declared values, attributes, and views remain unchanged
- **AND** the reviewed adoption plan contains no DNS additions, changes,
  replacements, or deletions

#### Scenario: Reconcile without changes after adoption

- **WHEN** the adopted root plans with default inputs and unchanged declarations
  and provider records
- **THEN** the plan retains the adopted bindings and contains no DNS additions,
  changes, replacements, or deletions
- **AND** unrelated provider records remain unchanged

#### Scenario: Provision missing declarations

- **WHEN** fresh complete provider inventories prove the owner's declared records
  are absent and conflict-free
- **THEN** the reviewed DNS plan creates only those declared missing records
- **AND** existing provider identities and attributes remain unchanged
- **AND** a subsequent DNS plan contains no changes

### Requirement: Scoped DNS execution and offline checks

The `dns.plan`, `dns.show`, and `dns.apply` wrappers SHALL select `dns=1` through
the repository AL flow using `src_infra_openhands`, and operate on `module.dns`
and its dependencies in the existing `tf_setup` root and Vault HTTP backend.
Apply SHALL consume the reviewed saved plan. Ordinary VM setup SHALL retain its
`tf=setup` and `xoa_login=1` labels, XO authentication, and existing backend paths.
Secret values SHALL remain in injected variables. Real DNS credentials and
Vault policy grants SHALL be prerequisites for operational DNS calls. RouterOS
DNS credentials SHALL remain isolated from unrelated RouterOS resources.
Successful scoped DNS execution SHALL not establish VM or service health.
The package SHALL expose a format test that does not authenticate to Vault or
contact DNS providers.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** the scoped DNS wrapper selects `dns=1`
- **THEN** it selects DNS credential injection through the owner AppRole and
  existing backend without starting XO login
- **AND** the targeted plan and saved-plan apply cover DNS and its dependencies,
  without establishing VM or service health
- **AND** disabled DNS resources do not imply offline provider configuration

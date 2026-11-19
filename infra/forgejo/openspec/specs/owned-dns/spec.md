# owned-dns Specification

## Purpose

Define Forgejo DNS management through the owner's `tf_setup` root,
including canonical declarations, scoped credentials, and offline source checks
before adopting live records.

## Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled before adoption and SHALL remain enabled
by default after authorized adoption. Adoption SHALL bind existing provider records
to the owner's state without adding, changing, replacing, or deleting DNS records.
Subsequent reconciliation of unchanged declarations SHALL preserve owned and
unrelated records.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before adoption
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect the adopted configuration

- **WHEN** the adopted root is evaluated with default inputs
- **THEN** it reads this owner's declaration and enables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Adopt existing records

- **WHEN** the owner reviews the matched existing provider imports through
  `//infra/forgejo/tf_setup:dns.plan` and `:dns.show`, then applies that saved
  plan through `:dns.apply`
- **THEN** the owner state binds each declared record to its existing provider identity
- **AND** the reviewed adoption plan proposes no record additions, changes,
  replacements, or deletions

#### Scenario: Reconcile unchanged declarations

- **WHEN** the adopted root plans unchanged declarations through
  `//infra/forgejo/tf_setup:dns.plan` and applies its reviewed saved plan
  through `:dns.apply`
- **THEN** it makes no DNS resource changes
- **AND** owned records retain their declared values and views while unrelated
  records remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrapper SHALL select `src_infra_dc1_forgejo1` through the repository AL
flow and keep secret values in injected variables. Dedicated DNS commands SHALL
select `dns=1`, retain the existing setup backend, and target `module.dns` without
starting unrelated service authentication. DNS apply SHALL consume only a
reviewed saved plan. Real DNS credentials and Vault policy grants SHALL be
prerequisites for operational Terraform calls. RouterOS DNS credentials SHALL be
isolated from unrelated RouterOS resources. The package SHALL expose a format
test that does not authenticate to Vault or contact DNS providers.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** the wrapper selects its Terraform stage labels
- **THEN** it selects DNS credential injection through the owner AppRole with
  `dns=1`, the existing setup backend, and the `module.dns` target
- **AND** disabled DNS resources do not imply offline provider configuration

## MODIFIED Requirements

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

- **WHEN** the owner adopts matched existing provider records through
  `//infra/forgejo_runner/tf_setup:dns.plan`, `:dns.show`, and `:dns.apply`
- **THEN** the owner state binds each declared record to its existing provider identity
- **AND** the reviewed adoption plan proposes no record additions, changes,
  replacements, or deletions

#### Scenario: Reconcile unchanged declarations

- **WHEN** the adopted root plans and applies unchanged declarations through its
  scoped `//infra/forgejo_runner/tf_setup:dns` wrappers
- **THEN** it makes no DNS resource changes
- **AND** owned records retain their declared values and views while unrelated
  records remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrapper SHALL select `src_infra_forgejo_runner` through the repository
AL flow and keep secret values in injected variables. Dedicated DNS plan, show,
and apply targets SHALL select `dns=1` alone and target `module.dns` while reusing
the owning setup root and backend. DNS apply SHALL consume only a reviewed saved
plan. Ordinary setup targets SHALL retain their existing `tf=setup` flow.
Real DNS credentials and Vault policy grants SHALL be prerequisites for
operational Terraform calls. RouterOS DNS credentials SHALL be isolated from
unrelated RouterOS resources. The package SHALL expose a format test that does
not authenticate to Vault or contact DNS providers.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** a dedicated DNS target selects its execution label
- **THEN** it selects `dns=1` and DNS credential injection through the owner AppRole
- **AND** it retains the setup backend while limiting reconciliation to `module.dns`
- **AND** disabled DNS resources do not imply offline provider configuration

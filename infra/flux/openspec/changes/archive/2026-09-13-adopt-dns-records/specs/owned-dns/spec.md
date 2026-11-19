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

- **WHEN** the owner applies a reviewed saved import plan through
  `//infra/flux/tf_setup:dns.apply`
- **THEN** the owner state binds each declared record to its existing provider identity
- **AND** the reviewed adoption plan proposes no record additions, changes,
  replacements, or deletions

#### Scenario: Reconcile unchanged declarations

- **WHEN** the adopted root plans and applies unchanged declarations through its
  `//infra/flux/tf_setup` wrappers
- **THEN** it makes no DNS resource changes
- **AND** owned records retain their declared values and views while unrelated
  records remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrappers SHALL select `src_infra_flux` through the repository AL flow and
keep secret values in injected variables. The `dns.plan`, `dns.show`, and
`dns.apply` targets SHALL select `dns=1` and retain the owner's setup backend;
plans SHALL target `module.dns` and apply SHALL require a reviewed saved plan.
Ordinary setup wrappers SHALL retain their service authentication behavior.
Real DNS credentials and Vault policy grants SHALL be prerequisites for operational
Terraform calls. RouterOS DNS credentials SHALL remain isolated from unrelated
RouterOS resources. The package SHALL expose a format test without live credentials.
DNS-targeted validation SHALL NOT establish unrelated VM or service health.

#### Scenario: Validate the Terraform source

- **WHEN** the package format test executes
- **THEN** it checks the packaged Terraform configuration without live credentials
- **AND** existing non-DNS authentication and backend paths remain unchanged

#### Scenario: Prepare an operational Terraform invocation

- **WHEN** an authorized operator selects the owner's DNS wrapper
- **THEN** it selects DNS credential injection through the owner AppRole and setup backend
- **AND** its scoped plan and saved-plan apply cover DNS and its dependencies
- **AND** disabled DNS resources do not imply offline provider configuration

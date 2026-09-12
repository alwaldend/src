## MODIFIED Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled before adoption and SHALL default to
enabled in the reviewed adoption revision and after adoption. Adoption SHALL
import the exact existing provider records and require a scoped no-change import plan before saved-plan
apply. Subsequent reconciliation with unchanged root inputs SHALL preserve
adopted records, unrelated DNS records, and non-DNS resources in the root.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before adoption
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect the adopted default

- **WHEN** the reviewed adoption revision or an adopted root is evaluated with default inputs
- **THEN** it reads this owner's declaration and enables managed DNS records
- **AND** enabled ownership remains the default after import

#### Scenario: Adopt existing records

- **WHEN** this owner applies a reviewed saved import plan through
  `//infra/harbor/tf_setup:dns.apply`
- **THEN** each existing provider record is imported at its exact owner-local address
- **AND** a reviewed adoption plan proposes no record additions, changes, replacements, or deletions before apply

#### Scenario: Reconcile unchanged declarations

- **WHEN** the adopted root reconciles unchanged canonical declarations and root inputs
- **THEN** the plan proposes no changes to the adopted DNS records
- **AND** unrelated DNS records and non-DNS resources in the root remain unchanged

### Requirement: Scoped DNS execution and offline checks

The DNS wrappers SHALL select `src_infra_harbor` through the repository AL flow and
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

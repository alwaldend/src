## ADDED Requirements

### Requirement: Owner-local DNS configuration

The `tf_setup` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled until the shared migration adopts them.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

### Requirement: Scoped DNS execution and offline checks

The DNS wrapper SHALL select `src_infra_flux` through the repository AL flow and keep
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

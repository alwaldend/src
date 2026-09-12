## MODIFIED Requirements

### Requirement: Owner-local DNS configuration

The `tf` root SHALL consume this owner's canonical `dnsconfig.json` through
the shared DNS Terraform module, preserving declared record identities and views.
DNS resources SHALL default to disabled until the shared migration adopts them
and SHALL retain enabled ownership after authorized adoption into the owner's
state. Reconciliation of adopted records with unchanged declarations SHALL
propose no record additions, changes, replacements, or deletions.

#### Scenario: Inspect the preparatory configuration

- **WHEN** the checked-in root is evaluated with default inputs before authorized
  adoption
- **THEN** it reads this owner's declaration and disables managed DNS records
- **AND** the package contains the module and declaration inputs

#### Scenario: Inspect adopted source defaults

- **WHEN** the adopted host root is evaluated with its checked-in source defaults
- **THEN** record ownership is enabled for the owner's declared records and views
- **AND** the package retains the canonical declaration and shared module inputs

#### Scenario: Reconcile existing DNS records

- **WHEN** the host root plans against its adopted state and unchanged declarations
- **THEN** it proposes no record additions, changes, replacements, or deletions
- **AND** records outside this owner's declarations remain unchanged

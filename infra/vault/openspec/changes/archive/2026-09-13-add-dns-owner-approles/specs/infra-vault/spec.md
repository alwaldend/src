## ADDED Requirements

### Requirement: Give each DNS owner a component identity

Every owner that manages DNS records through Terraform SHALL use its existing
component AppRole or a dedicated owner AppRole when one is missing. A DNS-only
identity SHALL retain the shared AppRole module's own-state and named shared
secret access without cloud provisioning, SSH, or PKI permissions. Existing
component identities and unrelated authentication flows SHALL remain stable.

#### Scenario: Introduce a landing site's DNS stage

- **WHEN** a project has DNS records but no existing component AppRole
- **THEN** Vault configuration declares an AppRole named for that owner
- **AND** the identity can manage its own Terraform state without access to
  another owner's state

#### Scenario: Add host DNS management

- **WHEN** host_bot uses its dedicated DNS AppRole
- **THEN** its existing Ansible authentication remains unchanged

### Requirement: Restrict provider secret reads to owned DNS views

DNS identities SHALL receive read access to the existing Cloudflare credential
only when they own global records and to the existing RouterOS DNS credential
only when they own dc1 records. DNS access SHALL NOT grant writes to either
provider secret or duplicate credential values.

#### Scenario: Configure a global-only landing site

- **WHEN** an owner declares only global records
- **THEN** its DNS policy grants Cloudflare credential reads
- **AND** its DNS policy grants no RouterOS credential access

#### Scenario: Configure records in both views

- **WHEN** an owner declares global and dc1 records
- **THEN** its DNS policy grants reads for both existing credential references
- **AND** neither read policy grants mutation of the provider credential

## Purpose

Provide a dedicated Vault identity for download infrastructure while keeping
state, credential references, and provider permissions scoped to its needs.

## ADDED Requirements

### Requirement: Dedicated component authentication

Vault configuration SHALL declare `src_infra_download` using the existing
AppRole backend and reusable identity modules. Its state permissions SHALL
cover the component's own environment and DNS state keys, not another
component's state namespace.

#### Scenario: Inspect component state access

- **WHEN** the download identity's effective state policy is evaluated
- **THEN** it grants access to download-owned state and locks
- **AND** it does not grant access to unrelated component state

### Requirement: Required provider and SSH assignments

The identity SHALL expose the membership and identity outputs needed for
Yandex folder provisioning and XCP-ng resource-set assignment. It SHALL
receive the required component SSH permissions and DNS provider-secret reads
for global and local records without administrative provider credentials.

#### Scenario: Prepare both environment assignments

- **WHEN** provider owners consume the declared download identity
- **THEN** they can assign its Yandex folder and XCP-ng resource set
- **AND** no unrelated component's identity is reused for deployment

### Requirement: Secret references and certificate automation

Checked-in declarations SHALL contain secret references only. Download
certificate automation SHALL receive role-scoped EAB access for internal Vault
HTTP-01 certificates. The role SHALL permit only `alwaldend.com`,
`download.alwaldend.com`, and `www.alwaldend.com`, without subdomain or client
certificate permissions. Public Let's Encrypt HTTP-01 SHALL require no DNS
challenge credential. Adding this
identity SHALL NOT broaden unrelated identities' policies.

#### Scenario: Review the source change

- **WHEN** the role and its policy declarations are inspected
- **THEN** they name Vault paths and permissions without embedding credential values
- **AND** unrelated identity permissions remain unchanged

### Requirement: Independent identity state

The download identity SHALL be deployed from its own Terraform root and
backend under `infra/vault/approles/src_infra_download`. Shared identity IDs
SHALL be resolved by name through data sources rather than parent-module
inputs or another root's state. Shared group membership lists SHALL retain
one owner in the core Vault stage.

#### Scenario: Bootstrap without the download identity

- **WHEN** shared Vault mounts and the administrative/operator identities already exist
- **THEN** the standalone root can create the download identity using administrative authentication
- **AND** a subsequent core Vault apply resolves the new identity/group by name and adds their memberships

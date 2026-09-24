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
certificate automation SHALL receive access only to the selected existing
credential references needed by its DNS challenge workflow. Adding this
identity SHALL NOT broaden unrelated identities' policies.

#### Scenario: Review the source change

- **WHEN** the role and its policy declarations are inspected
- **THEN** they name Vault paths and permissions without embedding credential values
- **AND** unrelated identity permissions remain unchanged

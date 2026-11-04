## Purpose

Define a dedicated Vault identity for Vast.ai infrastructure with scoped
component secrets and access for the existing operator to bootstrap it.

## ADDED Requirements

### Requirement: Scoped Vast.ai AppRole

Vault Terraform SHALL declare `src_infra_vastai` using the shared AppRole module
and existing secrets mount and AppRole backend. Its operator membership SHALL
include `simeonwarren`. It SHALL disable the optional Yandex Cloud policy and
SHALL NOT add the identity to the broad infrastructure AppRole group.

#### Scenario: Inspect checked-in identity configuration

- **WHEN** the Vault Terraform root is packaged
- **THEN** it includes the Vast.ai AppRole declaration and operator membership
- **AND** credentials and Terraform state use the component's own KV subtree
- **AND** validation does not apply the role or generate live credentials

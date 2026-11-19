# Infrastructure Yandex Cloud

## Purpose

Describe organization-level Yandex Cloud folder provisioning derived from Vault
AppRole identities. The direct `infra/yandex_cloud` owner currently contains one
concrete organization package, `org1`, covered here. This baseline does not
assert that the cloud organization or its folders were inspected live.

Baseline source revision: `550d7e79b1f5fdbc2b6017b75178471d6914082f`.
Observation date: 2026-09-08. Sources are linked in full; no excerpts are used.

Sources: [top-level documentation](../../../README.md),
[organization package](../../../org1/README.md),
[organization targets](../../../org1/BUILD.bazel),
[Terraform packaging](../../../org1/tf/BUILD.bazel),
[folder derivation](../../../org1/tf/folders.tf),
[provider configuration](../../../org1/tf/provider.tf), and
[credential injection](../../../org1/al.lua).

## Requirements

### Requirement: Keep organization configuration under its owner

The Yandex Cloud tree SHALL expose its documentation through the infrastructure
documentation boundary and keep `org1`'s AL configuration and Terraform package
under `infra/yandex_cloud/org1`.

#### Scenario: Discover the configured organization

- **WHEN** a contributor follows the Yandex Cloud package's child documentation
- **THEN** `org1` identifies the current organization configuration
- **AND** its Terraform wrapper uses `//infra/yandex_cloud/org1:al`

### Requirement: Derive folders from Vault AppRole membership

Terraform SHALL read the Vault identity group named `approles`, resolve its
member entities, and instantiate the shared `yc_folder` module once per entity
name. Folder names SHALL replace underscores with hyphens, and secret references
SHALL use the corresponding `yandex.cloud/org1/folders/` path.

#### Scenario: Map an AppRole identity to a cloud folder

- **WHEN** the `approles` group contains an entity named `example_service`
- **THEN** the folder module receives the name `example-service`
- **AND** its secret reference is
  `yandex.cloud/org1/folders/example-service`
- **AND** its cloud identifier comes from the configured `cloud_id` variable

### Requirement: Package Terraform with declared modules and injected environment

The Terraform Bazel package SHALL include its provider lock, shared folder and
backend modules, and Vault backend and injector tools. Its operational wrapper
SHALL select the Terraform plugin and the `default_default` Vault environment;
the package SHALL also expose the configured Terraform test map.

#### Scenario: Resolve the organization Terraform entry point

- **WHEN** the organization's Terraform wrapper is assembled
- **THEN** its runfiles contain the declared local modules and provider lock
- **AND** runtime authentication is supplied by the owning AL/Vault flow

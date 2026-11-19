# Reusable Terraform modules

## Purpose

Provide reusable infrastructure declarations for Vault, virtual machines,
storage, and related services. This baseline was observed at repository
revision `550d7e79` on 2026-09-08. The collection remains in progress according
to its [README](../../../README.md). Sources for the
representative module contracts below are the
[AppRole module](../../../vault_approle/main.tf),
[transit-key module](../../../vault_transit_key/main.tf),
[backup-bucket module](../../../backup_bucket/main.tf), and
[module packaging](../../../vault_approle/BUILD.bazel).
These are declarative source contracts; no live provisioning is asserted.

## Requirements

### Requirement: Package reusable module source

The AppRole, transit-key, backup-bucket, and Proxmox VM module packages SHALL
each expose their `main.tf` through a same-named Bazel filegroup available to
repository subpackages.

#### Scenario: A repository target consumes a module

- **WHEN** a target depends on `//projects/tf_modules/vault_approle:vault_approle`
- **THEN** its declared source input includes the AppRole module's `main.tf`

### Requirement: Associate AppRoles with Vault identities and policies

The AppRole module SHALL declare an identity entity, backend alias, internal
group, and named AppRole with configurable token and secret-ID limits. Its
role policies SHALL combine the shared and AppRole-secret policies with caller
policies and the Yandex-folder policy unless that policy is disabled for the role.

#### Scenario: Caller disables the Yandex-folder role policy

- **WHEN** `disable_yc_folder_policy` is true
- **THEN** the role's policy list excludes the module's Yandex-folder policy while preserving shared, AppRole-secret, and caller policies

### Requirement: Separate transit encryption and decryption membership

The transit-key module SHALL declare a named Vault transit key, an encryption
policy and group, and a decryption policy and group. Decryptor member groups
SHALL also be included in the encryption group.

#### Scenario: A group is listed as a decryptor

- **WHEN** a caller adds a group ID to `decryptors_member_group_ids`
- **THEN** the module includes that group in both transit permission groups

### Requirement: Declare encrypted versioned backup storage

The backup-bucket module SHALL declare a versioned Yandex storage bucket with
default KMS encryption, grant read and write permissions to the supplied
service accounts, and store bucket identifiers in the caller-selected Vault
KV v2 location.

#### Scenario: Backup storage is configured

- **WHEN** the caller supplies the bucket name components, folder, service accounts, and Vault destination
- **THEN** the declarations bind bucket encryption to the module's KMS key and publish the bucket ID, folder ID, and bucket name to that Vault destination

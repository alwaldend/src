# Ansible Collection Specification

## Purpose

Package the repository's `alwaldend.main` Ansible collection and its external
collection dependencies for repository automation. This baseline records
checked-in behavior at revision `550d7e79b1f5fdbc2b6017b75178471d6914082f`,
observed on 2026-09-08; it does not assert that any role has been deployed.

Sources: [project README](../../../README.md),
[collection packaging](../../../BUILD.bazel),
[Galaxy metadata](../../../galaxy.yml),
[role packaging](../../../roles/BUILD.bazel),
[playbook packaging](../../../playbooks/BUILD.bazel),
and [dependency declarations](../../../include.MODULE.bazel).

## Requirements

### Requirement: Collection namespace and source layout

The collection target SHALL package Galaxy metadata, role sources, and playbook
sources beneath `alwaldend/main`.

#### Scenario: Assemble the first-party collection

- **WHEN** a repository consumer requests the collection packaging target
- **THEN** its package mappings SHALL include `galaxy.yml`, the `roles` tree,
  and the `playbooks` tree under the `alwaldend/main` prefix.

### Requirement: Role and playbook aggregation

The roles package SHALL aggregate the `role` target of each role subpackage,
and the playbooks package SHALL include the playbook directory's YAML files.

#### Scenario: Add a packaged role

- **WHEN** a role subpackage exposes the required `role` packaging target
- **THEN** the aggregate roles package SHALL include that target under `roles`.

### Requirement: Locked collection dependencies

The combined collections target SHALL place the first-party collection beneath
`collections/ansible_collections` and include the external Galaxy collection
repository resolved through `ansible_lock.json`.

#### Scenario: Assemble collections for automation

- **WHEN** a repository consumer requests the combined collections target
- **THEN** its inputs SHALL include the prefixed first-party collection and
  the lock-backed `com_ansible_galaxy` repository.

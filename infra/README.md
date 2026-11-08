---
title: Infra
description: Infrastructure tree
weight: 2
cascade:
  - categories:
      - infra
---

This tree contains
[infrastructure as code](https://en.wikipedia.org/wiki/Infrastructure_as_code).
Tracked source follows the repository's public-source policy. Infrastructure
facts are not confidential merely because they are operational, generated, or
live. Reports may include them unless they contain credentials, other secrets,
or personal information. Inspect raw state, plans, inventories, and decrypted
configuration because those artifacts can contain prohibited content; do not
track the artifacts themselves.

Each infrastructure project owns its specifications and maintained changes
in `<project>/openspec/`. Use the [pinned OpenSpec workflow](../tools/openspec/README.md)
with that project selected. These specifications describe checked-in definitions;
they do not establish deployed state or authorize infrastructure operations.

## Requirements

- Bazel targets MUST use repository-internal visibility.
- Infrastructure definitions MUST NOT be published as production artifacts.
- Infrastructure targets MUST NOT be dependencies of production build
  targets.
- Public checked-in documentation MAY be included in the repository
  documentation site, including non-secret, non-personal operational facts.

## Transitive dependency enforcement

`//infra:transitive_consumers` allows internal infrastructure and developer
trees, website documentation, and the explicitly listed developer consumers.
The inherited Gazelle directive generates the package declarations. Ordinary
project packages, including `projects/ansible_collection/roles/traefik`, cannot
depend directly on infrastructure artifacts. Bazel 9.3.1 has an observed
[indirect enforcement limitation](../projects/bazel_gazelle_enforce_transitive_visibility/README.md#enforcement-limitation-in-bazel-931)
when an intermediary declares its own policy; do not treat this as a complete
transitive audit.
An exception grants access to every target in the consuming package; targets
sharing a package cannot have different transitive policies.

## Shared AL helpers

`infra_xo_login(name = ..., labels = ..., oidc_provider = ...)` from
[`al_config.bzl`](al_config.bzl) declares a configuration target for the
Xen Orchestra login plugin with the shared XO and Vault endpoints. The caller
must supply the Vault OIDC provider name; the helper constructs its discovery
URL. The helper declares the
`//infra/xcp_ng/cmd/xo_login` executable dependency. Consumers select its labels
in their command wrapper. Registration does not authenticate until a matching
command selects the plugin.

Public target macros in `al_config.bzl` use the `infra_` namespace, including
`infra_dns` and `infra_ansible_keys`. Backend calls use `vault_tf_backend` from
`//tools/vault/tf_backend:defs.bzl`, which owns their shared argument structure.

## New project

The commands below are state-changing operator examples. An agent must use
`bazel`, apply the repository Terraform and secret-handling procedures,
and receive explicit authority for the exact operation and environment before
running an equivalent command.

- Create an approle:
  [example](https://github.com/alwaldend/src/blob/master/infra/vault/tf/approle_src_infra_yandex_cloud.tf)
- Add it to approles:
  [example](https://github.com/alwaldend/src/blob/master/infra/vault/tf/group_approles.tf)
- Run apply:
  ```sh
  bazel run //infra/vault/tf:tf.apply
  ```
- Update Yandex Cloud folders:
  ```sh
  bazel run //infra/yandex_cloud/org1/tf:tf.apply
  ```
- Update Proxmox resource pools:
  ```sh
  bazel run //infra/pve/tf:tf.apply
  ```
- Update Xen Orchestra AppRole resource sets:
  ```sh
  bazel run //infra/xcp_ng/tf:tf.apply
  ```
- Set up al config:
  [example](https://github.com/alwaldend/src/blob/master/infra/pve/BUILD.bazel)
- Configure the component's Terraform state through the existing
  [Vault HTTP backend flow](../tools/agents/skills/repo-infra/references/flow.md).
  The `tf_backend` plugin creates state and lock KV entries on first use.

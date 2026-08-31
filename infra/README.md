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

## Requirements

- Bazel targets MUST use repository-internal visibility.
- Infrastructure definitions MUST NOT be published as production artifacts.
- Infrastructure targets MUST NOT be dependencies of production build
  targets.
- Public checked-in documentation MAY be included in the repository
  documentation site, including non-secret, non-personal operational facts.

## New project

The commands below are state-changing operator examples. An agent must use
`bazel_agent`, apply the repository Terraform and secret-handling procedures,
and receive explicit authority for the exact operation and environment before
running an equivalent command.

- Create an approle:
  [example](https://github.com/alwaldend/src/blob/master/infra/vault/tf/approle_src_infra_yandex_cloud.tf)
- Add it to approles:
  [example](https://github.com/alwaldend/src/blob/master/infra/vault/tf/group_approles.tf)
- Run apply:
  ```sh
  bazel_agent run //infra/vault/tf:tf.apply
  ```
- Update Yandex Cloud folders:
  ```sh
  bazel_agent run //infra/yandex_cloud/org1/tf:tf.apply
  ```
- Update Proxmox resource pools:
  ```sh
  bazel_agent run //infra/pve/tf:tf.apply
  ```
- Set up al config:
  [example](https://github.com/alwaldend/src/blob/master/infra/pve/al.lua)
- Set up a bucket for Terraform state:
  [example](https://github.com/alwaldend/src/blob/master/infra/vault/tf/cloud.tf)

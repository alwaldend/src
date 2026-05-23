---
title: Infra
description: Infrastructure tree
weight: 2
cascade:
  - categories:
      - infra
---

This tree contains [IaC](https://en.wikipedia.org/wiki/Infrastructure_as_code)

## Requirements

- MUST NOT be public
- MUST NOT be published
- MUST NOT be used in builds

## New project

- Create an approle: [example](./dc1/vault/tf/approle_src_infra_yandex_cloud.tf)
- Add it to approles: [example](./dc1/vault/tf/approles.tf)
- Run apply:
  ```sh
  bazel run //infra/dc1/vault/tf:tf.apply
  ```
- Update Yandex Cloud folders:
  ```sh
  bazel run //infra/yandex_cloud/org1:tf.apply
  ```
- Update Proxmox resource pools:
  ```sh
  bazel run //infra/dc1/pve1:tf.apply
  ```
- Set up al config: [example](./dc1/pve1/al.lua)
- Set up a bucket for terraform state: [example](./dc1/vault/tf/cloud.tf)

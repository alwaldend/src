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
- Create a Yandex Cloud folder: [example](./yandex_cloud/org1/tf/folders.tf)
- Set up secrets: [example](./dc1/pve1/al.lua)
- Set up a bucket for terraform: [example](./dc1/vault/tf/cloud.tf)

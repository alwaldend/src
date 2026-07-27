---
title: ingress
description: ingress
tags:
  - ingress
  - public
  - traefik
---

## Run ansible

```sh
bazel run //infra/ingress/ansible
```

## Apply terraform

```sh
bazel run //infra/ingress/tf
```

## Update signed image url

- Create a signed url: https://yandex.cloud/ru/docs/storage/operations/objects/link-for-download
- Update the secret: `alwaldend.com/vault1/approles/src_infra_ingress/image`

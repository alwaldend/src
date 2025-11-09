---
title: Infra
description: Infrastructure tree
weight: 2
cascade:
  - categories:
      - infra
---

## Setup DNS

```sh
bazel run //infra/dns:preview
bazel run //infra/dns:push
```

## Setup Rancher

```sh
bazel run //infra/ansible:rancher1
```

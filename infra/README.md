---
title: Infra
description: Infrastructure tree
weight: 2
cascade:
  - categories:
      - infra
---

This tree contains [IaC](https://en.wikipedia.org/wiki/Infrastructure_as_code)

## Features

- SHOULD NOT be public
- SHOULD NOT be published
- SHOULD NOT be used in builds

## Setup DNS

```sh
bazel run //infra/dns:preview
bazel run //infra/dns:push
```

## Setup Rancher

```sh
bazel run //infra/ansible:rancher1
```

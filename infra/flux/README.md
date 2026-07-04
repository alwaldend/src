---
title: flux
description: Fluxcd deployment
tags:
  - flux
  - ansible
  - pve
  - k3s
---

## Links

- Site: https://fluxcd.io/

## Run flux CLI

```sh
bazel run //infra/flux
```

## Run bootstrap

```sh
ssh -L 6443:127.0.0.1:6443 -N flux.alwaldend.com &
bazel run //infra/flux:flux.bootstrap
```

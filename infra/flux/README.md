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

## Set up port forwarding

```sh
ssh -L 6443:127.0.0.1:6443 -N flux.alwaldend.com
```

## Run flux CLI

```sh
bazel run //infra/flux/cl:flux
```

## Run bootstrap

```sh
bazel run //infra/flux/cl:flux.bootstrap
```

## Run flux operator

```sh
bazel run //infra/flux/cl:op
```

## Generate secret id for the cluster

```sh
bazel run //infra/flux/cl:secret_id
```

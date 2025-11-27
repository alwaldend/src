---
title: Cl1
description: Infrastructure for cl1.dc1.alwaldend.com
---

## Deployment

- Setup the host:
  ```sh
  bazel run //infra/alwaldend.com/dc1/cl1/ansible:deploy
  ```
- Deploy the chart:
  ```sh
  ssh -NR 127.0.0.1:6443:127.0.0.1:6443 pi1.cl1.dc1.alwaldend.com &
  bazel run //infra/alwaldend.com/dc1/cl1/helm:deploy
  ```

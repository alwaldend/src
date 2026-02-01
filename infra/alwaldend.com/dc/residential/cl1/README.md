---
title: Cl1
description: Cluster 1
---

## Deployment

- Setup the host:
  ```sh
  bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy
  ```
- Deploy the chart:
  ```sh
  ssh -NR 127.0.0.1:6443:127.0.0.1:6443 pi1.cl1.dc1.alwaldend.com &
  bazel run //infra/alwaldend.com/dc/residential/cl1/helm:deploy
  ```

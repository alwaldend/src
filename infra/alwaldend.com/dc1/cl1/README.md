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
  bazel run //infra/alwaldend.com/dc1/cl1/chart:deploy
  ```

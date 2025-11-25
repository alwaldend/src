---
title: Cluster1
description: Infrastructure for cluster1.dc1.alwaldend.com
---

## Deployment

- Setup the host:
  ```sh
  bazel run //infra/alwaldend.com/dc1/cluster1/ansible:deploy
  ```
- Deploy the chart:
  ```sh
  bazel run //infra/alwaldend.com/dc1/cluster1/chart:deploy
  ```

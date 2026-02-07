---
title: Cl1
description: Cluster 1
---

## Deployment

- Setup hosts:
  ```sh
  bazel run //infra/alwaldend.com/dc/residential/cl1/ansible:deploy
  ```
- Deploy the chart:
  ```sh
  ssh -NL 127.0.0.1:6443:127.0.0.1:6443 ansible@mini1.residential.dc.alwaldend.com &
  bazel run //infra/alwaldend.com/dc/residential/cl1/helm:deploy
  ```

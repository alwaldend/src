---
title: Forgejo runner
description: Forgejo Actions runner deployment
tags:
  - forgejo
  - ansible
  - pve
---

## Deploy VMs

```sh
bazel run //infra/forgejo_runner/tf_setup # Create VMs
bazel run //infra/forgejo_runner/ansible # Configure VMs
```

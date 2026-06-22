---
title: forgejo1
description: forgejo setup for forgejo1.dc1.alwaldend.com
tags:
  - ansible
  - pve
  - forgejo
---

## Links

- Docs: https://forgejo.org/docs/latest/admin/installation/binary/
- Config reference: https://forgejo.org/docs/latest/admin/config-cheat-sheet/

## Deployment

```sh
bazel run //infra/dc1/pve1/tf_setup:tf.apply # Create hosts
bazel run //infra/dc1/pve1/ansible # Configure hosts
bazel run //infra/dc1/pve1/tf:tf.apply # Configure forgejo
```

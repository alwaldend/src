---
title: forgejo
description: forgejo.alwaldend.com
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
bazel run //infra/forgejo/tf_setup:tf.apply # Create hosts
bazel run //infra/forgejo/ansible # Configure hosts
bazel run //infra/forgejo/tf:tf.apply # Configure forgejo
```

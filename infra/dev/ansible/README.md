---
title: Development host deployment
description: One deployment for both development hosts
---

```sh
bazel run //infra/dev/ansible
bazel run //infra/dev/ansible -- --limit secure1.dev.alwaldend.com
```

One invocation runs both host-specific plays using the `src_infra_dev` AppRole.
Workstation storage and services apply only to dev1; secure1 keeps the host-only
runner baseline.

Build `:ansible_bin` and run `:syntax` for offline validation.

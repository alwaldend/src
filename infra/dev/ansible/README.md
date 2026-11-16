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

The workstation's managed known_hosts trusts secure1 certificates through the
canonical repository server CA and pins dev1's public host key. The dev1 key
was verified from its local host filesystem and strict SSH handshake during
the deployment; host-key checking remains enabled through hostname cutover.

Build `:ansible_bin` and run `:syntax` for offline validation.

To install or reconcile only the pinned standalone fj binary on the current
dev1 host, run `bazel run //infra/dev/ansible:fj_install`. It uses the same
packaged playbook's `fj` tasks with a local connection and privilege escalation
disabled, independently of the Vault/SSH deployment wrapper. `-- --list-tasks`
and `-- --syntax-check` inspect this selection without installing it.

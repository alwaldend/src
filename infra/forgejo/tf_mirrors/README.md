---
title: Native push mirrors
description: Forgejo native HTTPS push mirror configuration
---

Terraform bootstraps native Forgejo HTTPS push mirrors for repositories
declared in the catalog. The pinned Forgejo provider exposes pull
mirrors. The packaged `push_mirror` API command handles push-mirror bootstrap
using the same creation-only pattern as runner registration. Vault supplies
the separate GitHub and GitLab destination credentials under the Forgejo
AppRole's `push_mirrors/<destination>` paths. Tokens stay out of Terraform state.

The command adopts a matching HTTPS mirror without writing to Forgejo.
Before creating one, it checks the advertised destination heads and tags with
Git. Each existing destination ref must have the same object ID on Forgejo.
Seed or reconcile a destination using normal non-force pushes first.
The new mirror uses Forgejo's native LFS transfer, sync on commit, and an
eight-hour interval. The bootstrap rejects any conflicting existing entry.

The bootstrap requires Git on the controller.
It runs on resource creation, so ordinary refresh cannot check remote drift.
Removing the bootstrap resource leaves the live mirror in place.

Add a `github` or `gitlab` entry under `forgejo.push_mirrors` in the
repository's catalog file. See [`src.json`](../../repos/orgs/alwaldend/repos/src.json)
for an example. Set `username` and `vault_key` for each destination and
include the corresponding destination forge in the same repository file.
Use `oauth2` for GitLab or the PAT owner's GitHub login.
Store the `token` field in the Forgejo AppRole's `push_mirrors/<vault_key>`
record.

```sh
bazel run //infra/forgejo/tf_mirrors:tf.plan
bazel run //infra/forgejo/tf_mirrors:tf.apply
```

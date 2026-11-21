---
title: Native push mirrors
description: Forgejo native HTTPS push mirror configuration
---

Terraform bootstraps native Forgejo HTTPS push mirrors for repositories
declared in the catalog. The packaged `push_mirror` API command reconciles
mirrors on Terraform resource creation and when the desired branch filter
changes. Vault supplies
the separate GitHub and GitLab destination credentials under the Forgejo
AppRole's `push_mirrors/<destination>` paths. Tokens stay out of Terraform state.

The branch filter combines the catalog's default branch with `releases/*`.
Forgejo also sends all tags and uses its native HTTPS LFS transfer. The
mirrors synchronize on commits and every eight hours.

The command adopts a matching HTTPS mirror without writing to Forgejo.
Forgejo's API requires deleting and recreating a mirror to change its filter.
Before either creation or replacement, the command checks that the source
advertises trunk and that every existing destination head and tag matches
Forgejo's object ID. It also obtains the destination credential before
deleting the old entry. A failed replacement can be retried through Terraform.

The first filtered sync prunes previously mirrored development branches from
the destination. Confirm that pruning is intended before applying the plan.
The preflight requires their tips to be recoverable from the source. The
command requires Git on the controller. Refresh cannot detect remote drift.
Removing the Terraform resource leaves the live mirror in place.

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

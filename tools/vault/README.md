---
title: Vault
description: Vault
---

## Links

- Repo: https://github.com/hashicorp/vault

The root `//:vault` and `//:vault.*` labels delegate to wrappers in
`//tools/vault/cmd/workspace`, using `//tools/al:config`. `//tools/vault:vault`
is the standalone CLI; the wrappers add the existing AL environment injection.

The pinned upstream binary is owned by
[`third_party/com_hashicorp_vault`](../../third_party/com_hashicorp_vault/README.md).

[`cmd/ssh_key`](cmd/ssh_key/README.md) creates or reuses an Ed25519 SSH pair
directly in Vault KV v2. Callers select its mount and path and own deployment;
the tool preserves the cached identity and refuses implicit rotation.

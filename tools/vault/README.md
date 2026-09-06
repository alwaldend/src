---
title: Vault
description: Vault
---

## Links

- Repo: https://github.com/hashicorp/vault

The root `//:vault` and `//:vault.*` labels delegate to wrappers in
`//tools/vault/cmd/workspace`, using `//tools/al:config`. `//tools/vault:vault`
is the standalone CLI; the wrappers add the existing AL environment injection.

---
title: Vast.ai
description: Vast.ai infrastructure configuration and Vault authentication
---

This project owns the infrastructure configuration for [Vast.ai](https://vast.ai/).
It provides a dedicated Vault authentication context, API-key injection for
the CLI, and Terraform HTTP backend wiring. It does not provision resources.
The pinned upstream CLI is owned by [tools/vastai](../../tools/vastai/README.md).

## Vault authentication

The [AL configuration](al.lua) uses the `src_infra_vastai` AppRole, declared in
[Vault Terraform](../vault/tf/approle_src_infra_vastai.tf). The existing operator
can bootstrap this role. Its component secret namespace is
`secrets/alwaldend.com/vault1/approles/src_infra_vastai/`.

An authorized apply of the owning Vault Terraform root must deploy the role
before authenticated commands can work. After deployment, inspect the role's
capabilities without reading a secret:

```sh
bazel run //infra/vastai:vault -- token capabilities \
  secrets/data/alwaldend.com/vault1/approles/src_infra_vastai/vastai
```

## Authenticated CLI

Store the API key as field `api_key` at
`alwaldend.com/vault1/approles/src_infra_vastai/vastai` in the `secrets` mount.
The component CLI selects the `vastai=1` injector, which reads that field and
sets `VAST_API_KEY` for the pinned upstream command:

```sh
bazel run //infra/vastai -- show user
```

The key is passed through the environment, without a command-line credential
or persistent CLI key file. Standard component Vault commands do not select
this injector, so credential setup does not require an existing API key.

## Terraform state

The `tf=main` label registers the existing Vault HTTP backend at
`alwaldend.com/vault1/approles/src_infra_vastai/tf_backend/tf` in the `secrets`
mount. A future Terraform command must package `//tools/vault/tf_backend` and
select that label. No Terraform stage is defined yet.

Building these targets does not authenticate or change live infrastructure.

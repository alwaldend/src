---
title: Vast.ai
description: Vast.ai infrastructure configuration and Vault authentication
---

This project owns the infrastructure configuration for [Vast.ai](https://vast.ai/).
It provides a dedicated Vault authentication context, API-key injection for
the CLI, and Terraform HTTP backend wiring. Its [bootstrap root](tf/README.md)
creates the provider's test team and generated project token without renting
compute or transferring credits. A catalog supports additional teams and keys
through the shared team module.
The pinned upstream CLI is owned by [tools/vastai](../../tools/vastai/README.md).

[GPU rental targets](rentals/README.md) use a dedicated `alwaldend` team key
from Vault with an aliased provider passed into a child module. They select
the cheapest on-demand offer and include an explicit interactive rent-and-release
test with separate state. They do not switch credential context or transfer credit.
The existing [team catalog](tf/README.md) declares the team's full-access key
through the shared team module and stores it in Vault. The stored issuer creates keys, but Vast.ai rejected selecting
alwaldend on the generated key with HTTP 400. Live rental acceptance remains
blocked by that team-binding operation.

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

Store the personal bootstrap API key as field `api_key` at
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
mount. The bootstrap root selects that label and packages the backend plugin.
The provider project's test root has its own AppRole, injector, generated token,
and state namespace; it does not reuse this personal bootstrap credential.
The team catalog uses each destination's Vault identity for generated tokens:
the provider project's AppRole for its test key and this component's AppRole
for the rental key. Both retain their existing namespace permissions.

Building these targets does not authenticate or change live infrastructure.

---
title: Vast.ai
description: Account, team, and rental infrastructure with Vault authentication
---

This project manages Vast.ai through three independent Terraform roots:

| Root                                  | Credential                      | Owns                                    | State suffix  |
| ------------------------------------- | ------------------------------- | --------------------------------------- | ------------- |
| [Account](tf/README.md)               | Personal account key            | Team identities and account settings    | `tf`          |
| [Alwaldend](alwaldend/README.md)      | Operator-provided alwaldend key | Full-access team keys and team settings | `alwaldend`   |
| [Rental test](rentals_test/README.md) | Generated alwaldend rental key  | One test GPU contract                   | `rental_test` |

[config.json](config.json) defines environment settings once. The account and
alwaldend roots both read that section because teams do not inherit account
variables. ENABLE_HTTPS, WEB_ENABLE_HTTPS and WEB_ENABLE_AUTH are all true.
The account root currently declares only alwaldend. Two former imported teams
remain tracked pending credential recovery for deletion; the managed test team
and its test objects have been deleted. See the [live evidence](openspec/changes/enable-team-gpu-rentals/validation.md).

## Authentication

The component uses the `src_infra_vastai` Vault AppRole declared in
[Vault Terraform](../vault/tf/approle_src_infra_vastai.tf). Secret references use
the `secrets` mount and its existing component namespace. Credentials are
injected into process environments; they are not command-line arguments or
persistent CLI key files.

| Vault path below `alwaldend.com/vault1/approles/src_infra_vastai/` | Owner                                  |
| ------------------------------------------------------------------ | -------------------------------------- |
| `vastai`                                                           | Operator-provided personal account key |
| `vastai/alwaldend`                                                 | Operator-provided root team key        |
| `vastai/alwaldend/rentals`                                         | Team root's generated rental key       |

The existing [alwaldend root](alwaldend/README.md) also declares a
dedicated token into the `src_infra_image_gen/vastai` component namespace.
The existing team Vault writer has an exact-path delivery permission for that
entry. The infrastructure rental reads it with its own AppRole and uses the
generated team key to provision the instance.

Each entry stores field `api_key`. The team root preserves its own seed key.
New API keys inherit the issuer's current account. The pinned CLI has no team
selector for `create api-key`; its team help directs operators to select an
existing team in the console and create its initial key there. Configure each
provider with the appropriate credential rather than switching API-key context.

The following one-time bootstrap command writes only the root team key and
keeps it out of shell history and process arguments:

```bash
(
  set +x
  unset vastai_team_key
  IFS= read -r -s -p 'Alwaldend team key: ' vastai_team_key || exit 1
  printf '\n'
  [ -n "$vastai_team_key" ] || exit 1
  printf '%s' "$vastai_team_key" |
    bazel run //infra/vastai:vault -- kv put -mount=secrets \
      alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend api_key=-
)
```

The user authorized this operator bootstrap; IaC contains only its reference.
The [pinned CLI](../../tools/vastai/README.md) can verify each context:

```sh
bazel run //infra/vastai -- show user
bazel run //infra/vastai/alwaldend -- show user
```

The HTTP backend stores each root's lock and state under its configured
`tf_backend/STATE_SUFFIX` entry. The image-generation token and its component Vault entry are owned by the existing
team backend without reusing the destroyed personal key.
Building or testing these targets does not authenticate or change live infrastructure.

---
title: Vast.ai team catalog
description: Team identities, generated project keys, and Vault delivery
---

This root owns the team and key catalog in `teams.json`, the personal-account
`ENABLE_HTTPS` setting, and Vault delivery of generated keys. Each entry with
`api_keys` calls `projects/tf_modules/vastai_team`, which owns both the team and
its keys. Entries without `api_keys` retain imported identity-only resources
with `prevent_destroy`. The root transfers no credit and declares no compute.
Generated credentials remain sensitive and are not root outputs.

Alwaldend and its full-access rental key belong to this catalog. No separate
credential root or credential backend is declared. The rental targets consume the
key from Vault and keep their compute lifecycle in separate states.

## Existing team addresses

The live alwaldend identity is currently at
`vastai_team.teams["alwaldend"]`. Its declared module address is now
`module.teams["alwaldend"].vastai_team.team`. Reconcile that address within
this same state before applying the rental key. This state migration is pending;
preserve the existing funded team and reject team creation, replacement, or
deletion. Do not import a duplicate identity or transfer it to another backend.
Remove temporary migration blocks after reconciliation and verify the final
owning plan without them. No temporary import or moved blocks remain in source.

The provider-test module, key, and `project_tokens` secret addresses are
unchanged. Other imported team identities retain their existing addresses.

## Vault identities

The personal Vast.ai issuer and HTTP state backend use `src_infra_vastai`.
AL injects both existing Vault AppRole identities for secret delivery:

| Catalog `vault_auth`         | Vault identity                           | Provider          |
| ---------------------------- | ---------------------------------------- | ----------------- |
| `provider_project` (default) | `src_projects_vastai_terraform_provider` | `vault`           |
| `default`                    | `src_infra_vastai`                       | `vault.component` |

The default provider retains the project's environment credential. The component
provider receives its sensitive token through `TF_VAR_component_vault_token`;
no resource receives a token input. Each AppRole writes only its existing
namespace. This requires no broader Vault policy or secret sharing between roles.

Secret destruction uses Vault soft deletion; the team module's key resource
revokes the API key. Purging Vault metadata is outside the roles' permissions.

## Add teams and keys

Add a distinct team object to `teams.json`. Its optional `name` defaults to the
catalog key. Include `api_keys`, even when empty, to retain its module address.
Each key declares a name, `secret_name`, permission object, and optional
`vault_auth`. A permission value of `null` explicitly requests default full
access; permission objects describe scoped keys. Imported-team owner binding
requires full access. Duplicate secret destinations and unknown Vault identities
fail planning.

```json
{
  "alwaldend": {
    "api_keys": {
      "rentals": {
        "name": "infra_vastai_alwaldend",
        "secret_name": "alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend",
        "vault_auth": "default",
        "permissions": null
      }
    }
  }
}
```

A new team receives zero credits and uses its generated management key for
same-apply project-key bootstrap. For a pre-existing team not yet owned in state,
declare its module and import its numeric account ID at that module's team
resource through the packaged import target before applying keys:

```sh
bazel run //infra/vastai/tf:tf.import -- RESOURCE_ADDRESS TEAM_ID
```

Already imported teams require address reconciliation rather than a second
import. Imports read identity without recovering old management secrets.
Imported-team key creation uses the configured issuer and the imported account
ID. Imported team destruction is refused by the provider.

## Review and apply

```sh
bazel run //infra/vastai/tf:tf.plan
bazel run //infra/vastai/tf:tf.apply
```

Review the plan displayed by the interactive apply. After address reconciliation,
the rental-key addition must create only its key and Vault entry, with no team
changes. The existing project key and secret must remain unchanged. Do not save
plans or use automatic approval.

The stored issuer can generate a full-access key, but the previous generated
key received HTTP 400 while selecting alwaldend. That tracked key was revoked
and no rental Vault entry was created. Binding and live rental acceptance remain
unresolved. This catalog refactor does not justify retrying that unchanged
operation. See the [rental targets](../rentals/README.md) and the
[reconciliation evidence](../openspec/changes/reconcile-live-account/validation.md).

`//infra/vastai/test/catalog:catalog_test` exercises this actual root with mocked
provider identities, checks both Vault routes, and rejects duplicate destinations
and unsupported identities. It uses no live credentials or infrastructure.

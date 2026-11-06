---
title: Vast.ai provider bootstrap
description: Terraform-created test team and generated project token
---

This root manages the personal-account `ENABLE_HTTPS` setting with value
`true` and the team and project-key catalog in `teams.json`.
The single `teams` map is loaded from `teams.json`, with one configuration
object per team. Entries without `api_keys` use `vastai_team.teams` and retain
their imported account IDs in state, with `prevent_destroy` enabled. Their
management credentials are unavailable, so import does not generate keys.
Import blocks were removed after adoption; a moved block preserves the previous
resource addresses.
Entries with `api_keys` retain the reusable `projects/tf_modules/vastai_team`
module and write each generated project token directly to its configured Vault
secret. The catalog contains three imported team identities and a provider-test
team with its project token. Existing module and secret addresses are preserved.
The root transfers no credit and declares no compute resources. Generated
credentials remain sensitive and are not root outputs.

The provider generates a separate management key to select the team context;
it never changes the personal bootstrap key's context. The project API key
grants user and team administration for the test fixtures, with no instance
permissions. The personal-key injector and Terraform backend use the
`src_infra_vastai` identity. The Vault provider authenticates separately as
`src_projects_vastai_terraform_provider` and uses its existing permissions on
its own secret subtree. No extra Vault policy or cross-role secret access is
required. Deploy that project AppRole before running bootstrap.

Secret destruction uses standard Vault soft deletion; the Vast.ai key resource
revokes the token. Purging Vault metadata is outside the role's permissions.

## Add another team

Add an object to `teams.json` with a distinct team key. The optional `name`
defaults to that key. Omit `api_keys` for a team identity without managed
project tokens, or include the map to select the credential-generating module.
Each API key declares a name, a JSON-compatible permission object, and a distinct
`secret_name` beneath the selected AppRole's namespace. For example:

```json
{
  "test_secondary": {
    "name": "test_secondary",
    "api_keys": {
      "test_reader": {
        "name": "test_secondary_reader",
        "secret_name": "alwaldend.com/vault1/approles/src_projects_vastai_terraform_provider/vastai/secondary",
        "permissions": { "api": { "user_read": {} } }
      }
    }
  }
}
```

Retain `api_keys`, even when empty, to preserve the module address. Adding or
removing that field changes the owning resource path and requires an explicit
state migration for an existing team. Duplicate Vault
destinations fail planning. Another project can reuse the module from its own
authenticated root; Vault identity and secret destinations remain caller-owned.

## Review and apply

```sh
bazel run //infra/vastai/tf:tf.plan -- -out=/absolute/private/path/bootstrap.tfplan
bazel run //infra/vastai/tf:tf.show -- /absolute/private/path/bootstrap.tfplan
bazel run //infra/vastai/tf:tf.apply -- /absolute/private/path/bootstrap.tfplan
```

Review the saved plan before applying it. Keep plans private because they can
contain generated credentials. After bootstrap, run the provider project's
non-compute test root using its own AppRole and generated-token injector.

## Reconcile existing resources

The authenticated import command uses this root's existing Vault-backed state
and provider credentials:

```sh
bazel run //infra/vastai/tf:tf.import -- RESOURCE_ADDRESS RESOURCE_ID
```

Declare a supported resource with its intended configuration before importing
it, then review a fresh plan. The provider supports import for instances,
SSH keys, environment variables, team roles, and owned team identities. Imported
teams have no recovered management credential and cannot be destroyed by the
provider. API keys cannot be imported. Adding an existing team to `teams.json`
requires importing it at its selected resource address before applying; a
catalog entry alone does not adopt a live resource.
Resources in the provider's test root retain that root's separate ownership.

The [reconciliation record](../openspec/changes/reconcile-live-account/validation.md)
records the observed inventory and unresolved adoption requirements. A
no-change plan covers declared resources; it does not establish that every
live resource is managed. API-masked variable values require their intended
values before they can be reproduced in configuration.

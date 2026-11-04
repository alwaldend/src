---
title: Vast.ai provider bootstrap
description: Terraform-created test team and generated project token
---

This root manages the team and project-key catalog in `teams.json`.
Each team uses the reusable `projects/tf_modules/vastai_team` module, and each
project token is written directly to its configured Vault secret. The default
catalog declares a provider-test team and its project token.
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

Add an entry to `teams.json` with a distinct team name and its `api_keys` map.
Each key declares a name, a JSON-compatible permission object, and a distinct
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

The same resource declarations manage every catalog entry. Duplicate Vault
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

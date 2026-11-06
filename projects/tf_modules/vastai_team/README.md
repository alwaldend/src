---
title: Vast.ai team module
description: Team and project API-key composition with caller-owned authentication
---

This module declares one Vast.ai team and a map of project API keys through the
in-tree provider. It declares no compute or Vault resources. The caller selects
the provider and owns state and secret storage. Generated key outputs are
sensitive and belong directly in secret resources.

A new team receives zero credits. Its project keys use the management key
generated during the same apply. Existing module resource addresses and outputs
remain unchanged; upgrading a managed team requires no state migration.

For an existing owned team, import its numeric account ID into the module's
`vastai_team.team` resource before applying keys. Import reads identity without
creating a team or recovering an old secret. Imported teams have no management
key, so project keys use the caller's configured provider and the imported team
ID. This binding requires `permissions = ""` for default full access. There is
no separate existing-team input or ownership mode.

```hcl
provider "vastai" {
  alias = "issuer" # VAST_API_KEY is injected from Vault.
}

module "alwaldend" {
  source    = "../../../projects/tf_modules/vastai_team"
  providers = { vastai = vastai.issuer }

  name = "alwaldend"
  api_keys = {
    rentals = { name = "infra_vastai_alwaldend", permissions = "" }
  }
}
```

Use the owning root's packaged import command with address
`module.alwaldend.vastai_team.team` and the verified account ID. If another
Terraform address or state already owns that team, migrate that ownership
instead of importing a duplicate. A plan for an imported team must show no team
creation or replacement. Imported team destruction is refused; cleanup can
revoke the declared project keys without deleting the existing team.

The provider verifies the generated key's team selection before returning its
sensitive value. Rejected selection fails apply and leaves a tracked key for
explicit cleanup. Instantiate this module with `for_each` for multiple teams.
Vault identity and destination paths belong to the caller; the module grants no
Vault permissions.

Terraform CLI fixtures verify managed bootstrap, imported-team key creation,
rejected activation and scoped input, unchanged plans and addresses, and key
cleanup while preserving imported teams. They require no live account. Live
alwaldend activation remains blocked by the recorded HTTP 400 response.

---
title: Alwaldend team configuration
description: Full-access team keys and shared team environment settings
---

This separate Terraform root uses the operator-provided alwaldend team key
from Vault as its default provider credential. It configures the shared
`environment_variables` map from `../config.json` and creates the full-access
keys listed in `teams.alwaldend.api_keys`. Each generated key is stored at
`secrets/alwaldend.com/vault1/approles/src_infra_vastai/vastai/alwaldend/KEY_NAME`.
Secret values are sensitive and are not outputs.

The dedicated image-generation key name comes from `teams.alwaldend.imagegen_api_key_name`.
The existing alwaldend Vast credential generates that key. The existing
`src_infra_vastai` Vault identity writes its value directly to
`alwaldend.com/vault1/approles/user_simeonwarren/imagegen/vastai`, using the
default Vault provider. Its [delivery policy](../../vault/tf/approle_src_infra_vastai.tf)
is limited to this exact entry. The personal rental reads the stored key with
its own user identity.

The existing key was transferred from its former credential state into this
root without regeneration. The [handoff record](../openspec/changes/relocate-imagegen-credentials/design.md)
records its completed reconciliation and no-change plan.

The root seed key at the parent `vastai/alwaldend` path is operator-owned.
This root neither replaces nor deletes it. The account root owns team creation;
this root owns no teams, rentals or credit transfers. Its Vault-backed backend
is `tf_backend/alwaldend`, independent of the account and rental states.

```sh
bazel run //infra/vastai/alwaldend -- show user
bazel run //infra/vastai/alwaldend:tf.plan
bazel run //infra/vastai/alwaldend:tf.apply
```

Verify that the seed selects alwaldend and review the interactive apply before
approving it. Do not save plans or use automatic approval. The component AppRole owns team credentials/state and the exact-path token
delivery permission. The owning Vault workflow manages that narrow permission. The
[acceptance record](../../vault/openspec/changes/authorize-vastai-key-delivery/verification.md)
records its deployed access and the team root's no-change plan.

Rent and release compute through the separate [rental test](../rentals_test/README.md),
which reads the generated `rentals` key. Running this root does not rent compute.

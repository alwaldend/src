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

The dedicated image-generation key name comes from `teams.alwaldend.image_gen_api_key_name`.
The existing alwaldend Vast credential generates that key. The existing
`src_infra_vastai` Vault identity writes its value directly to
`alwaldend.com/vault1/approles/src_infra_image_gen/vastai`, using the
default Vault provider. Its [delivery policy](../../vault/tf/approle_src_infra_vastai.tf)
is limited to this exact entry. The [infrastructure rental](../../image_gen/README.md) reads the stored key with
its own AppRole. The writer can delete metadata for this exact entry so its
managed key delivery is removed completely during teardown.

The former personal key and its delivery entry were destroyed on 2026-11-10.
The fresh component key has been created and delivered to its own namespace.
The [component receipt](../../image_gen/openspec/changes/add-imagegen-rental/design.md#credential-only-deployment)
records its reviewed creation without a rental. No state handoff or import is needed. The
[change record](../openspec/changes/relocate-imagegen-credentials/design.md)
retains the previous handoff history and current teardown evidence.

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
records the legacy grant; the component receipt above records the new deployment.

Rent and release compute through the separate [rental test](../rentals_test/README.md),
which reads the generated `rentals` key. Running this root does not rent compute.

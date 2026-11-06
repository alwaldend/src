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
approving it. Do not save plans or use automatic approval. The component
AppRole injects provider credentials and writes only its existing namespace.

Rent and release compute through the separate [rental test](../rentals_test/README.md),
which reads the generated `rentals` key. Running this root does not rent compute.

## Current desired state, 2026-11-10

The user requested destruction of the personal image-generation stack and fresh
IaC under infra/image_gen with its own AppRole. At 09:24 UTC, a reviewed targeted
apply destroyed key 36816419 and its current personal Vault entry: 0 additions,
0 changes, 2 destroys. The shared team issuer, unrelated keys and team remain.
Historical Vault/state versions remain protected; no state handoff is reused.

The new `vastai_api_key.image_gen` and `vault_kv_secret_v2.image_gen` declare a
fresh key named `infra_vastai_alwaldend_image_gen`, delivered directly to
`src_infra_image_gen/vastai`. The new entry deletes all versions on destroy;
Vault's exact metadata-delete grant is owned by infra/vault. These replacement
resources are not deployed. Earlier handoff receipts below describe the now
revoked personal key.

## Legacy ownership before teardown

The existing infra/vastai/alwaldend root owns the dedicated image-generation API key and its regular-user Vault entry. The team issuer remains unchanged; its existing infrastructure Vault identity delivers the generated key through an exact-path policy. The team root contains no recipient authentication or token variable. The personal rental only consumes that token. There is no additional credential root or AppRole. The delivery policy is owned by infra/vault.

## Completed state handoff

On 2026-11-09, the user explicitly approved the one-time exception for in-memory JSON-state processing and direct Vault state updates. The task-private procedure moved exactly API key 36816419 and its existing Vault entry from user_simeonwarren/tf_backend/imagegen_keys into src_infra_vastai/tf_backend/alwaldend. Both backends were locked and written with Vault KV compare-and-set; both state lineages and all pre-existing team resources were preserved. No raw state or credential value was printed or saved locally. Historical states remain protected in Vault.

The API key remains vastai_api_key.imagegen. The Vault resource became vault_kv_secret_v2.imagegen. The recipient provider alias used during the original handoff is retired by the subsequent team-writer correction. The key ID, name, value, and live secret path/value were unchanged; the old state's current resource list and outputs are empty.

The normal //infra/vastai/alwaldend:tf.plan then refreshed the existing key and user Vault entry and reported no changes. No key rotation, infrastructure apply, Vault policy change, or GPU rental occurred during the handoff.

## Verification

The actual account/team catalog fixture passes with mocked providers. Forty scoped checks and affected semantic lint pass. The public loopback handoff fixture covers success, wrong-key identity, destination conflicts, existing locks, and a partial-write failure that retains both locks for reconciliation. The one-time handoff is complete; no migration block or repeated migration command remains in desired configuration.

## Team writer correction

The team root now uses only its own Vast and Vault identities. Its exact-path
delivery policy is deployed, and the user project consumes the generated key
through its own identity. The existing team root reports no changes. A reviewed
refresh-only apply saved the default-provider association without changing any
remote resource. The [Vault acceptance record](../../../../vault/openspec/changes/authorize-vastai-key-delivery/verification.md)
owns deployment and permission evidence. No key was rotated.

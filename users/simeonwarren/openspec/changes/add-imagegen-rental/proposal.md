## Why

A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## What Changes

Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Use the newest published Vast Forge Neo/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Do not execute a live rental as implementation validation.

## Capabilities

### New Capabilities

- `imagegen-rental`: A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## Impact

Owner-local source, packaging, documentation, and offline validation. Reuses the existing external-data and other pinned Terraform providers; no new external library or GPU deployment.

## Dedicated team token

Create a dedicated image-generation API key in alwaldend using the existing
team issuer credential and write its generated value under the regular user's
`imagegen/vastai` Vault entry. The team writer has an exact-entry delivery policy; the existing user identity
reads that entry and provisions with the generated key. No extra role is added. The user explicitly requested token creation/storage, but not a GPU
rental. The existing infra/vastai/alwaldend root owns token generation and the regular user Vault write. The authorized state handoff preserved the existing key.

## Why

A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## What Changes

Add imagegen/tf_setup with random_password, Vault KV storage, and the reusable exposed-mode module. Reuse the existing user_simeonwarren AppRole, regular user secret namespace, and independent Vault-backed state. Declare a direct published Vast ComfyUI/CUDA 12.9 image tag without a template by default; provide optional declarative public model downloads; allow an optional reviewed environment-only template; retain the USD 1/hour cap. Offline validation SHALL NOT rent capacity; the user separately authorized one live rental after security review.

## Capabilities

The current follow-up switches this consumer to the separate ComfyUI module and
Chroma1-HD-Flash with the same encoder, VAE, credential ownership and rental
constraints. Publish the architecture and PR; leave the destroyed rental absent.

### New Capabilities

- `imagegen-rental`: A personal image generation rental needs generated credentials retained in Vault before paid instance creation.

## Impact

Owner-local source, packaging, documentation, and offline validation. Reuses the existing pinned Terraform providers; removes the external-data provider and adds no executable dependency. The approved Chroma bundle supplies model weights through the existing downloader. Live deployment is separately authorized by the user.

## Dedicated team token

Create a dedicated image-generation API key in alwaldend using the existing
team issuer credential and write its generated value under the regular user's
`imagegen/vastai` Vault entry. The team writer has an exact-entry delivery policy; the existing user identity
reads that entry and provisions with the generated key. No extra role is added. The user explicitly requested token creation/storage, but not a GPU
rental. The existing infra/vastai/alwaldend root owns token generation and the regular user Vault write. The authorized state handoff preserved the existing key.

## Why

The image-generation stack needs infrastructure ownership with an independent
Vault identity and state, rather than a personal user's credentials and backend.
The user requested full teardown before creating the replacement IaC.

## What Changes

Move the component to `infra/image_gen`. Declare the `src_infra_image_gen`
AppRole in the Vault owner, a fresh component KV/state namespace, and exact-path
team-key delivery from `infra/vastai/alwaldend`. Generate and persist the HTTPS
login before renting. Retain verified Secure Cloud, at least 32 GiB VRAM,
100 GB disk, a USD 1/hour ceiling and the approved Vast ComfyUI image.

Keep the requested native Qwen-Image-Edit-2511 bundle and neutral reference
workflows as defaults. The replacement remains undeployed. The latest teardown
request supersedes the earlier Qwen redeployment instruction.

## Capabilities

### New Capabilities

- `imagegen-rental`: Infrastructure-owned ComfyUI rental, credentials and reference workflows.

## Impact

Move source, documentation and the maintained change into the new owner. Update
Vault policy and team-key desired state; no new executable dependencies. Destroy
the legacy rental, generated login, dedicated team key and obsolete delivery
permission. Preserve the shared team, issuer, unrelated keys and personal identity.
No state handoff, credential reuse or new rental is part of this delivery.

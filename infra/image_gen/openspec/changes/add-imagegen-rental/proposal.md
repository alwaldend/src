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
workflows as defaults. After the teardown and credential-only deployment, the
user separately authorized one rental and live workflow acceptance with neutral
references. Retain the current identity, price and exposure boundaries.

## Capabilities

### New Capabilities

- `imagegen-rental`: Infrastructure-owned ComfyUI rental, credentials and reference workflows.

## Impact

Move source, documentation and the maintained change into the new owner. Update
Vault policy and team-key desired state; no new executable dependencies. Destroy
the legacy rental, generated login, dedicated team key and obsolete delivery
permission. Preserve the shared team, issuer, unrelated keys and personal identity.
No state handoff or old credential reuse is required. The authorized rental is
created through its owning interactive apply; offline source checks do not rent.

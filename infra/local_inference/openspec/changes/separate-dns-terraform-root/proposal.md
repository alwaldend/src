## Why

DNS currently uses a filtered view of the GPU provisioning root and shares its
state. It needs an independent lifecycle, provider selection, and state owner.

## What Changes

- Add `infra/local_inference/dns` with Cloudflare only and a dedicated backend.
- Reuse the canonical component DNS document, AppRole, and shared DNS module.
- Keep only the active DNS project, without migration scaffolding, import
  inputs or legacy backend wiring. Existing record ownership is already adopted.
- Remove DNS management from the provisioning PR #233.

## Capabilities

### New Capabilities

- `independent-dns`: Separate DNS root/state and reviewed ownership transfer.

## Impact

The user authorized this source extraction, state migration, and apply. No GPU,
Yandex folder, host deployment, DNS record recreation, or credential rotation
is included. A separate master-based PR owns the independent DNS project.

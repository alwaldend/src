## Why

Secure1 needs NAS-backed cache capacity presented as a virtual disk. Keeping
the storage connection in XCP-ng gives the guest an ordinary local filesystem.

## What Changes

- Bootstrap a dedicated SMB3 storage repository using the existing XO helper.
- Delegate that repository to the development resource set.

## Capabilities

### New Capabilities

- `runner-cache-storage`: NAS-backed virtual disk storage for secure1.

### Modified Capabilities

None.

## Impact

Uses existing XO, Go and SMB dependencies. NAS and guest desired state remain
with `infra/nas` and `infra/dev`; the exact credential grant is in `infra/vault`.
No deployment is authorized by this source change.

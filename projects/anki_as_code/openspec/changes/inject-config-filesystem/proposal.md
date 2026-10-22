## Why

Configuration operations currently access the host filesystem directly, preventing callers from supplying browser-compatible storage. The application needs a construction-time Afero filesystem that owns all configuration access.

## What Changes

- Add an optional Afero filesystem to application construction, retaining OS storage by default.
- Route all application filesystem operations through that filesystem: configuration discovery, reads, validation and writes, collection source checks, export/state output and archive staging/publication.
- Preserve server storage and hosting boundaries; no server migration is included. Backend libraries retain their existing filesystem contracts.
- Verify complete configuration workflows against memory-backed and restricted filesystems.

## Capabilities

### New Capabilities

- `configuration-filesystem`: Injected configuration storage, filesystem-relative paths and complete file-only operation without host storage.

### Modified Capabilities

None.

## Impact

Application composition, configuration/layout, collection source checks, export, rebuild, shared path handling, Go dependency exposure and documentation. Reuse the already pinned Afero v1.11.0 dependency; no new dependency or version upgrade is required. Server and backend-library filesystem implementations remain outside this change.

## Context

The local inference caller uses Ubuntu 24.04 and the current T3 release.
The npm launcher needs Node >= 20.11; the standalone bundle is already
published by the requested T3 Code project and needs no separate Node pin.

## Goals / Non-Goals

**Goals:** Add one opt-in installation mode to the existing shared owner.
**Non-Goals:** Change other deployments, add a new package manager, or deploy.

## Decisions

Preserve npm as the default and add `standalone` as an explicit choice.
The caller supplies the immutable release URL and SHA-256. Extract the full
bundle because it includes native modules and browser assets. Keep runtime
and parent directories administrator-owned; state, work, and home remain
owned by the ordinary T3 service account.

Use handlers for changed startup inputs and `state: started` for steady
state. Environment writes may contain credentials and need `no_log`.
Standalone mode uses native pairing/session authentication; preserve the
legacy token entry only in the existing npm mode.
An optional initial settings object seeds the userdata settings file with
`force: false`, preserving the application's subsequent UI changes.

## Risks / Trade-offs

- Upstream CLI flags and bundle layout are checked against the pinned release.
- Missing checksum, failed extraction, environment changes, and service-path
  changes must be covered by source checks before deployment acceptance.

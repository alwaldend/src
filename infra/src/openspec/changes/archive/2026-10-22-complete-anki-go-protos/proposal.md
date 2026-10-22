## Why

The aggregate Anki proto target includes all 25 pinned schemas, but Go generation selects only 17. A complete runtime CLI requires the missing schema bindings and access to every registered service.

## What Changes

- Add Go message/client targets for ankidroid, ankihub, ankiweb, frontend, github, i18n, image_occlusion and tags.
- Expose all 25 schema packages to intended consumers, retaining upstream import paths and restricted visibility.
- Reuse Anki 26.09.3 at the existing immutable commit; preserve native sources, pins and dependency locks.

## Capabilities

### New Capabilities

- `anki-go-protos`: Complete Go bindings for the pinned upstream schema set.

### Modified Capabilities

None.

## Impact

third_party/com_github_ankitects_anki BUILD definitions and documentation, with complete bindings exercised by the paired runtime CLI E2E. This is repository dependency integration and does not provision infrastructure or add external software.

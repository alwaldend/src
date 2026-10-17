## Why

The supplied October 17 collection supersedes the tracked October 7 export.
Grammar lessons need padded numeric prefixes for predictable browsing.

## What Changes

- Replace exported collection declarations from the supplied archive, removing stale declarations while preserving recognizable resource paths where possible.
- Prefix grammar lessons with `001. ` and subsequent numbers, alphabetically within each proficiency level.
- Split Guzhenren chapter vocabulary and sentences into aligned, numbered reading groups; number sentence files in chapter order.
- Save the native export's importable state sidecar and explicit moves from exported to organized paths.
- Use the current CLI's `version = "1"` schema, with native identities in the exported state and explicit moves connecting exported paths to organized paths.
- Document and verify the refreshed collection and ordering.

## Capabilities

No application behavior changes; this is a collection data refresh and adoption of the existing CLI format (`skip_specs: true`).

## Impact

Only the user-owned Anki collection text and its documentation change. The supplied archive remains local; media and review history are not committed. No dependency or live Anki profile changes are required.

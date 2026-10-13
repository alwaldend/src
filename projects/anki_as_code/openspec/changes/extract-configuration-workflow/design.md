## Context

The current CLI exposes native resource identities in editable TOML. The user replaces that workflow with paths and collection-backed state, similar to Terraform. Earlier completed GUID work is historical; the current contract below supersedes it.

## Decisions

Keep the fixed collection/deck hierarchy and format version 1. The root supports moved from/to pairs for individual configuration paths or directory prefixes. Parsed note/note-type wrappers retain each collection-relative source path; deck groups retain their marker path and directory-derived title. Resource configurations contain no top-level native IDs or note GUIDs. References use configuration paths; scheduling presets use readable names. Template/field order comes from arrays, and sort-field selection uses a field name. Reuse upstream presentation/configuration messages while omitting native identity, provenance, and derived metadata from export; reconciliation rejects attempts to set those native-only fields.

Load, Dump, and Export remain the pure internal/config APIs. Export writes supplied path/config wrappers directly; callers own valid fixed-layout paths, completeness, validation, staging, cleanup, and publication. Remove GUID generation, its staging helpers, and its CLI/app command. Config dump remains available for unfinished declarations without native access or relationship validation.

Keep one Go module. internal/config is the stable facade and dump-operation owner; config/layout owns fixed hierarchy discovery, file I/O, and direct object export; config/codec owns filesystem-independent TOML/JSON conversion through protobuf. The dependency direction is facade to layout/codec, then layout to codec. Keep layout helpers private and limit codec exports to the operations actually used across these packages. Preserve the current schema, supported version, datetime policy, error causes, unchecked writer behavior, and public entry points. Use existing CLI E2E coverage to verify this structural refactor without adding implementation-mirroring tests.

Collection state and native identity resolution remain in #128. This extraction supplies current schemas, source ownership, moved data, and file codecs without allocating identities, validating collection relationships, or opening Anki. Runtime models expose native identities solely to the reconciliation consumer.

## Verification

Write E2E failure cases before implementation. Verify current-schema loading/dumping/object export, source-path ownership, removed identity fields/command, optional presence and exact values, errors, and read-only sources. Native E2E verifies export initialization, persistent state after reopen/archive rebuild, import, moves, creation/deletion, untracked resource safety, structural migrations, filtered cards, scheduling/history/preset/media preservation, no-op convergence, and failure publication. Emit repeatable TOML/JSON/archive artifacts. Delivery receipts own exact-candidate validation and publication.

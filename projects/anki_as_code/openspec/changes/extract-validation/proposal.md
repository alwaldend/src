## Why

Validation should be reviewable independently of collection storage and reconciliation. Extract the existing pure validation code from PR #128 with focused corrections for invalid resource identities and conflicting deck-kind declarations.

## What Changes

- Add the existing validation package and its required Anki format helpers against master.
- Reject nonpositive and duplicate note identities; reject both deck kinds based on optional-field presence.
- Reuse pinned upstream regular/filtered deck and note-type/template configurations as fields in editable protobuf wrappers.
- Preserve package boundaries, documentation, and explicit Bazel dependencies.
- Keep collection commands as placeholders; no SQLite, archive, or reconciliation code is included.
- Keep PR #128 buildable against master until this prerequisite merges, then remove the shared additions by rebasing.

## Capabilities

### New Capabilities

- `desired-state-validation`: Reject invalid or duplicate note identities and conflicting deck-kind declarations.

### Modified Capabilities

None.

## Impact

Two internal Go packages, editable protobuf wrapper fields, upstream schema aliases, and project documentation. Reuse existing pinned dependencies; add no external dependency.

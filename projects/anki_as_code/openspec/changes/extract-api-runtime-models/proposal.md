## Why

PR #128 currently contains both the collection implementation and its shared data contracts. Extract the API schema and runtime models into an independently buildable prerequisite PR against master.

## What Changes

- Add the existing collection protobuf contract and generated-language Bazel targets.
- Add runtime note/deck wrappers and collection snapshots using generated messages directly.
- Preserve the command scaffold and leave serialization adapters and reconciliation behavior in PR #128.

## Capabilities

### New Capabilities

- `collection-models`: Shared generated serialization types and nonserialized runtime metadata.

### Modified Capabilities

None.

## Impact

Only `projects/anki_as_code/api/collection`, `internal/model`, and this OpenSpec change are extracted. Reuse existing protobuf toolchains and upstream Anki definitions; add no external dependencies. Both PRs target master. PR #128 retains these files until this prerequisite merges, then a rebase removes the shared diff.

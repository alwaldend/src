## Why

The Anki implementation PR combines command wiring and project integration with collection business logic. Extract the project scaffold so those boundaries can be reviewed and merged independently.

## What Changes

- Introduce the Go Cobra command tree with required flags and explicit operation placeholders.
- Add Bazel targets, project documentation, and the shared-site landing registration.
- Keep collection implementation and its CLI workflow test in PR #128.

## Capabilities

### New Capabilities

- `cli-scaffold`: Command discovery, argument validation, and explicit unavailable-operation errors.

### Modified Capabilities

None.

## Impact

Owns projects/anki_as_code command and documentation scaffolding plus projects registry integration. Uses the existing Cobra dependency. Both PRs target master; PR #128 will remove the extracted changes after this prerequisite merges.

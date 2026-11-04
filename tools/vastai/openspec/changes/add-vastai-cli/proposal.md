## Why

The repository lacks a pinned Vast.ai command-line tool. Add the requested upstream CLI through the existing Bazel Python dependency workflow.

## What Changes

- Expose `//tools/vastai` and document usage.
- Pin `vastai==1.10.0` in the component-local Python manifest and generate its native lock.

## Capabilities

### New Capabilities

- `vastai-cli`: Repository-internal access to the pinned upstream CLI.

### Modified Capabilities

None.

## Impact

Adds tools/vastai and registers a component-local pip repository. No host installation or cloud operation is included.

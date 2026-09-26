## Why

Independent placement leaves many isolated shapes. Users need a control that
encourages touching groups without increasing the number of occupied cells.

## What Changes

- Add `--clusterization`, a finite factor from 0 through 1, defaulting to 0.
- Rearrange occupied cells to increase horizontal and vertical contacts while
  preserving the occupied-cell count for the same seed and other settings.
- Preserve existing images at factor zero and deterministic replay at all factors.
- Document the control and retain reproducible comparison images in E2E outputs.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: configurable clustering and its interaction with density.

## Impact

The CLI, internal renderer, E2E acceptance suite, and project documentation change.
Rendering continues to use only Go's standard library.

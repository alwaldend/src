## Why

The user wants four configurable replication controls for directional inheritance,
forward growth, crowding avoidance, and occasional sideways branching, with
separate rendered examples to compare their effects.

## What Changes

- Add four finite 0..1 controls: `--replication-inheritance`,
  `--replication-forward-bias`, `--replication-crowding`, and `--replication-branching`.
- Default them to zero to preserve existing seeded images.
- Keep exactly two children or none, eight-neighbor adjacency, retained parents,
  initial-population density, and compact as the default algorithm.
- Document control interactions and provide repeatable comparison renders.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: configurable replication placement and direction memory.

## Impact

CLI, replication placement, E2E coverage, and project documentation change.
Rendering remains within the Go standard library. No additional dependency.

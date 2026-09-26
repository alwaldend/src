## Why

The user wants winding, branching strands with open gaps. Maximizing occupied
neighbor pairs packs cells into blobs and does not achieve that visual result.

## What Changes

- Make positive clusterization control directional strand growth and occasional branches.
- Prefer uncrowded eight-neighbor placements while retaining the sampled cell count.
- Replace the monotonic neighbor-count guarantee with a strand-oriented visual contract.
- Preserve default independent output, palette, geometry, density endpoints, and replay.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: organic clustering with winding strands and open gaps.

## Impact

Internal placement, CLI help, E2E morphology checks, and docs change. No new dependencies.

## Why

The user clarified that clustering should count diagonal neighbors too.

## What Changes

- Count all eight surrounding grid cells when measuring clustering.
- Preserve density and seeded replay, and describe circle corner geometry.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: include diagonal adjacency in clusterization.

## Impact

Update the internal neighbor count, command help, E2E measurements, and docs.
No dependencies or runtime programs are added.

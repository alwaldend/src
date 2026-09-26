## Why

The user wants to select the original compact clustering, the current organic
strands, or a new algorithm in which cells reproduce or stop reproducing.

## What Changes

- Add `--clusterization-algorithm compact|strands|replication`, defaulting to compact.
- Restore compact eight-neighbor clustering without changing existing strand images.
- Use density as the initial replication population and clusterization as the
  probability that a cell produces two distinct empty adjacent children.
- Process each cell once; death leaves the visible cell unchanged and adds no children.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: selectable clustering and replication population growth.

## Impact

CLI options, renderer placement, E2E coverage, docs, and the owning Go BUILD
declarations change. Only the standard library is used.

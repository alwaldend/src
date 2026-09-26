## Why

Regional replication improves clustering but still produces granular patches. The
requested winding strands need persistent direction and a main-stem/side-branch
hierarchy. The user approved a separate algorithm that may extend one child.

## What Changes

- Add `tips` beside compact, strands, and replication; compact remains default.
- Grow sparse initial cells into gradually curving paths, with occasional shorter
  side branches and stopping at collisions, boundaries, or an exhausted budget.
- Retain density as initial occupancy and use clusterization for growth length.
- Produce reproducible 1024-square comparisons with individual white pixels.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: selectable persistent tip growth and its lifecycle.

## Impact

CLI selector/help, generator implementation, command-level acceptance, project
documentation, and specifications. Go standard library only; no dependencies.
Existing algorithms and replication controls retain their behavior.

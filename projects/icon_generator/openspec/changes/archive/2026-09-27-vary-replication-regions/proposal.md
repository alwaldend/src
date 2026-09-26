## Why

The user prefers mixed-growth variant 10 but finds its texture too uniform.
Add the proposed smooth regional reproduction variation and uneven initial
population, retaining variant 10's directional controls for comparison renders.

## What Changes

- Add `--replication-field-scale` in logical cells, default 128.
- Add `--replication-growth-variation` as a maximum local probability offset,
  and `--replication-seeding-variation` as seeding contrast, both 0..1 and off
  by default.
- Use one deterministic smooth field for both effects, while retaining density
  as the average initial occupancy probability.
- Retain paired births, parents, eight-neighbor adjacency, and compact default.
- Generate separate reproducible comparison PNGs and provide a table of links.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `random-icon-generation`: regional initial occupancy and reproduction probabilities.

## Impact

Project CLI, sampling and replication internals, E2E tests, and docs change.
Use only the standard library. Family-trait variation is outside this change;
the selected experiment implements the first two proposed improvements.

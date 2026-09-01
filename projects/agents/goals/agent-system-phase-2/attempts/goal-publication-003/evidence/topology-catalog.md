# Topology catalog compilation evidence

## Binding

- Goal: `agent-system-phase-2`, resource version 9 (attempt `goal-publication-003`)
- Criteria revision: 2
- Candidate: deterministic topology compiler at
  `tools/agents/cmd/topology_check` and checked artifacts at
  `tools/agents/catalogs/topology.json` / `topology.md`.

## Live compilation

The compiler ran against the exact workspace:

```text
catalogID: agent-system.topology
completeness: complete
status: ok
bounds: eligible 43, emitted 43, unavailable 0
trees: 6, components: 28, workspaces: 9
limitations: []
conflicts: []
```

The JSON and Markdown digests match, and `--check` passes against the tracked
artifacts.

## Test coverage

- `//tools/agents/catalog/v1alpha1:v1alpha1_test` — envelope validation,
  deterministic canonical JSON, strict decoding, JSON/Markdown parity.
- `//tools/agents/cmd/topology_check:topology_check_test` — complete fixture,
  negative completeness on missing BUILD, byte-identical regeneration, stale
  checked-JSON rejection.
- `//tools/agents/cmd/topology_check:topology_check_check` — checked-artifact
  drift test (passes; fails on stale bytes).
- `//tools/agents/...` and `//:buildifier_test` and `//projects/goal/...` —
  green on the exact candidate; goal store validates 3/3.

## Honest verdict

This evidence supports `bounded-catalogs` progress for the topology slice
only: one complete, bounded, deterministic, checked catalog exists. It does
not support complete `bounded-catalogs` acceptance, which requires the
Policy, Action, Capability, WorkspaceCheck, Goal, and Index compilers.

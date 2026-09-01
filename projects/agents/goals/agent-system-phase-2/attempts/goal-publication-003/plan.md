# Phase 2B starter: shared catalog envelope and TopologyCatalog

## Bindings

- Goal: `agent-system-phase-2`
- Goal generation: 2
- Lifecycle generation: 4
- Criteria revision: 2
- Expected checkpoint resource version: 8
- Prior attempt: `goal-publication-002` (`refine`)
- Work type: `change`

The goal tool binds the portable goal-state and criteria digests when it
publishes this attempt.

## Target gap

Phase 2A (`goal-recovery`) is accepted; the remaining Phase 2 criteria
(`bounded-catalogs`, `system-index`, `context-capsule`, `runtime-isolation`,
`resource-baseline`, `legacy-migration`) are open. The catalog-inputs audit
prescribes a safe sequence: freeze the shared deterministic catalog envelope
first, then implement the static Topology compiler on the registered
`projects/*` universe before the richer catalogs and index.

## Decision and smallest slice

Implement the Phase 2B shared envelope plus the TopologyCatalog compiler.
Reuse the Phase 1 registry and `projects/*` owner-local facts as inputs;
produce one versioned deterministic JSON catalog, a checked Markdown render
from the same data, a `--check` mode that fails on completeness gaps, and
negative fixtures. This freezes the envelope (schema URI, kind, id,
derivationVersion, producerRef, sourceRevision, inputs, bounds, completeness,
limitations, conflicts, items, digest, canonical JSON rules) that the Policy,
Action, Capability, WorkspaceCheck, Goal, and Index compilers reuse.

Reject a database, daemon, network access, or stateful generation. The
compiler runs offline and reads bounded source authorities directly.

## Ready workstreams

1. Coordinator: implement the shared envelope types, canonical JSON, digest,
   and validation in `tools/agents/catalog/v1alpha1`.
2. Coordinator: implement the Topology compiler and `--check` semantics under
   `tools/agents/cmd/topology_check`.
3. Coordinator: add deterministic and negative completeness fixtures, Bazel
   wiring, checked JSON/Markdown outputs, registry generated-artifact entry,
   and README documentation.

Parallel worker boundaries are not needed for this slice: the modules are
small and tightly coupled, and the coordinator owns the canonical outputs.

## Acceptance for this attempt

- The envelope schema rejects unknown fields, duplicate identities, absolute
  or escaping paths, unknown enums, unsorted set-like arrays, and
  non-complete results without limitations.
- TopologyCatalog deterministically derives from the registry and
  `projects/*` README/BUILD facts with provenance (input digests, source
  revision) and no generation timestamp or absolute path.
- `--check` fails on an omitted or invalid eligible project and passes on a
  complete registered universe.
- Checked JSON and Markdown render the same model and state the JSON digest;
  fixtures prove both.
- `bounded-catalogs` remains `unverified` after this attempt: only the
  envelope plus Topology are implemented; Policy, Action, Capability,
  WorkspaceCheck, Goal, and Index stay open.

## Fixed regressions

- `//projects/goal/...` focused tests remain green.
- `//tools/agents/...` Phase 1 checker and new compiler tests pass.
- Repository Buildifier stays clean if Bazel metadata changes.
- `git diff --check` clean.

## Strategy reset

Reset if deterministic offline compilation turns out to require a database or
daemon, or if the registered `projects/*` universe cannot be represented
without inference. In that case narrow Topology to the exact verifiable
registered subset and retain the boundary evidence.

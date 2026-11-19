# Phase 2B starter: shared catalog envelope and TopologyCatalog result

## Outcome

This attempt closes with `refine`. It implements the Phase 2B shared
deterministic catalog envelope and the TopologyCatalog compiler, following the
safe sequence fixed by `goal-publication-001`. It claims `bounded-catalogs`
support for the topology slice only; the Policy, Action, Capability,
WorkspaceCheck, Goal, and Index compilers remain open Phase 2 criteria.

## Implemented

- `tools/agents/catalog/v1alpha1`: shared catalog envelope types, strict
  canonical JSON with self-digest semantics, deterministic validation
  (identity, path escaping, input digests, bounds, completeness and
  limitation coupling, sorted set-like arrays, duplicate and unknown-value
  rejection), and a human Markdown renderer from the same data.
- `tools/agents/cmd/topology_check`: offline, bounded compiler over the Phase 1
  registry, top-level boundary READMEs, registered `projects/*` README/BUILD
  facts, and tracked `MODULE.bazel` roots. It emits portable JSON plus
  checked Markdown, records every input with a content digest, derives a
  stable content-addressed `sourceRevision` over the input universe, and
  exposes `--check` that fails on completeness gaps or stale checked bytes.
- Checked artifacts: `tools/agents/catalogs/topology.json` and
  `tools/agents/catalogs/topology.md` regenerate via
  `//tools/agents/cmd/topology_check:topology_update` and are verified by
  `//tools/agents/cmd/topology_check:topology_check_check`.
- `tools/agents/declarations/registry.json`: registered the generated topology
  JSON artifact with its check target.

## Verification

- `//tools/agents/...`: envelope, parity, compiler, and Phase 1 suites pass;
  Phase 1 `phase1_check` accepts the new registered artifact.
- `//tools/agents/cmd/topology_check:topology_check_check` passes on the
  checked artifacts and fails when the checked JSON is stale (regression
  covered by `TestTopologyCheckFailsOnStaleJSON`).
- Negative completeness: `TestTopologyCompileFailsOnMissingProjectBuild`
  proves `--check` fails when an eligible project lacks its BUILD.
- Determinism: `TestCanonicalJSONTopologyDeterministicDigest` and
  `TestTopologyCompileNonDeterministicInputRejected` prove byte-identical
  regeneration.
- Live workspace catalog: `complete`, 6 trees, 28 components, 9 workspaces,
  no limitations or conflicts; JSON and Markdown digests match.
- `//:buildifier_test` and `//projects/goal/...` remain green; goal records
  validate (3/3).

## Not in scope

Policy, Action, Capability, WorkspaceCheck, Goal, and AgentSystemIndex
compilers remain open. `context-capsule`, `runtime-isolation`,
`resource-baseline`, and `legacy-migration` remain open. This attempt does not
claim complete `bounded-catalogs` acceptance.

## Dominant remaining failure

The remaining five static catalogs must each freeze their schema and
completeness gates against owner-local facts; the GoalCatalog must wait for a
pure read-only recovery-aware inspector before ingesting continuation state.

## Process audit

The envelope-first split kept the compiler, fixtures, and parity tests small
and reusable. The main feedback cost was iterating the frontmatter parsing
against real project READMEs (sites lists before statuses, missing H1s) and
the sourceRevision determinism decision; both are now covered by regressions.
No parallel worker was needed: the module is small and the coordinator owns
the canonical outputs.

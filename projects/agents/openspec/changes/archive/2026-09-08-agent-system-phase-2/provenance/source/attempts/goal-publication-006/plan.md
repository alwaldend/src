# goal-publication-006: Capability, WorkspaceCheck, Goal, and Index catalogs

## Goal

Complete the `bounded-catalogs` and `system-index` criteria with the
CapabilityCatalog, WorkspaceCheckCatalog, GoalCatalog, and AgentSystemIndex
compilers plus checked artifacts.

## Criterion

> Versioned bounded topology, policy, action, capability, workspace-check, and
> goal catalogs derive from owner-local facts with deterministic provenance,
> conflicts, bounds, and completeness checks.
>
> One bounded AgentSystemIndex contains only catalog identities, versions,
> input digests, conflicts, and query routes without duplicating catalog
> bodies.

## Plan

### Capability compiler (`tools/agents/cmd/capability_check`)

- Inputs: registry skills, runtimeTools, directBinaries, operation-file
  providers, `.agents/skills` discovery links, and each canonical
  `SKILL.md` plus its BUILD declaration.
- Emits `tools/agents/catalogs/capability.{json,md}` with skills (canonical
  path, discovery link, layer, activation, exclusions, capability refs,
  dependencies, conflicts, provider requirements, context cost, evaluation
  maturity) and providers (runtime tools, direct binaries, operation
  providers).
- Completeness: registry/discovery sets identical; every link relative
  non-escaping; frontmatter identity matches; dependency/conflict refs
  resolve; every registered runtime tool/direct binary appears.
- `--check` drift gate; Bazel update/check targets.

### WorkspaceCheck compiler (`tools/agents/cmd/workspace_check`)

- Inputs: tracked `MODULE.bazel` roots, `.bazelignore`, root local overrides,
  documentation aggregation, and full-check scripts.
- Emits `tools/agents/catalogs/workspace-check.{json,md}` with workspace
  records and projection membership (bazelIgnore, rootOverride,
  docsAggregation, fullCheck) plus phases.
- Completeness: every tracked MODULE root has a unique normalized path and
  module name; all projections have exactly the expected membership.

### Goal compiler (`tools/agents/cmd/goal_check`)

- Inputs: registered owner-local goals roots.
- Uses the pure read-only goal inspector (goal CLI `doctor` / validated
  store) — no mutating cleanup — to emit validated `identity` +
  `coarseStatus` only for complete records; invalid/interrupted records are
  `unavailable` with a bounded reason.
- Emits `tools/agents/catalogs/goal.{json,md}`.

### Index compiler (`tools/agents/cmd/index_check`)

- Reads the six catalog descriptor files (topology, policy, action,
  capability, workspace-check, goal) and emits `tools/agents/catalogs/index.json`
  + `index.md` with descriptors (id, kind, schema, derivation version,
  digest, input digests, completeness, query routes) and cross-catalog
  conflicts. No payloads embedded.
- Byte ceiling + rejection of embedded payloads via schema validation.

### Tests

- Compiler unit tests for complete and negative fixtures.
- Schema validation tests for capability/workspace/goal/index.
- Checked drift gates for all four.

## Out of scope

- `context-capsule`, `runtime-isolation`, `resource-baseline` run in
  subsequent attempts.

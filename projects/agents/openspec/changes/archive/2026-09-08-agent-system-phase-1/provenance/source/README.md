# Normalize agent semantics and declare safety contracts

> Generated bounded projection. Edit `goal.yaml`, `criteria.yaml`, and attempt records only through the goal tool.

- Goal ID: `agent-system-phase-1`
- API version: `goals.alwaldend.com/v1alpha1`
- Resource version: `14`
- Generation: `2`
- Lifecycle generation: `8`
- Scope: `project`
- Outcome: `achieved`
- Execution: `paused`
- Active attempt: `—`

## Relationships

- Parent: —
- Depends on: `repo-agent-system`
- Supersedes: —

## Acceptance criteria

- `shared-contracts` (r1, required): Shared v1alpha1 contracts round-trip deterministically and reject malformed identity, unknown effectful operations, widened authority, and incompatible information flows. — Evidence: Run focused schema round-trip and rejection tests.
- `registered-universe` (r1, required): A cheap report-only completeness check covers the closed registered agent-reachable universe and names missing or unclassified entries. — Evidence: Run the registered-universe completeness test.
- `information-policy` (r1, required): Public, secret, and personal-information classes remain independent from repository visibility, build-consumer, and publication policy axes. — Evidence: Run information-flow fixtures and inspect owner policy.
- `scratch-isolation` (r1, required): Task and run scratch is namespaced; concurrent tasks do not collide and stale workers fail expected-revision publication. — Evidence: Run concurrent namespace and stale-publication tests.
- `safe-operations` (r1, required): Every removed mutating alias has a safe explicit replacement and provider mutation ownership is unambiguous. — Evidence: Inspect operation declarations and run negative fixtures.
- `resource-baseline` (r1, required): A revision-bound Phase 1 baseline records numeric ceilings for correctness, unsafe actions, calls, context size, cold and warm time, universe coverage, and reused checks without estimating unavailable measurements. — Evidence: Validate the baseline fixture and its recorded ceilings.
- `exact-candidate-validation` (r1, required): All affected packages, schemas, declarations, generated checks, documentation, and goal records pass focused validation against the exact delivered candidate. — Evidence: Inspect the candidate-bound validation evidence.
- `encountered-bug-policy` (r1, required): Repository-wide policy requires primary agents to fix small bounded bugs encountered in repository tooling or the affected project instead of silently working around them, while subagents remain scoped and substantial redesigns or rewrites… — Evidence: Inspect the root agent policy and the goal CLI regression.

## Recent attempts

- [`phase1-completion-001`](attempts/phase1-completion-001/) — `integration`, `closed`, resource version `2`, criteria r3
- [`bazel-scratch-001`](attempts/bazel-scratch-001/) — `change`, `closed`, resource version `2`, criteria r3
- [`phase1-contracts-001`](attempts/phase1-contracts-001/) — `change`, `closed`, resource version `2`, criteria r3

## Record map

- [`goal.yaml`](goal.yaml): machine-authoritative goal state
- [`criteria.yaml`](criteria.yaml): versioned acceptance criteria
- [`attempts/`](attempts/): isolated attempt records and evidence

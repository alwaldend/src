# Generate a trustworthy system map and bounded context capsule

> Generated bounded projection. Edit `goal.yaml`, `criteria.yaml`, and attempt records only through the goal tool.

- Goal ID: `agent-system-phase-2`
- API version: `goals.alwaldend.com/v1alpha1`
- Resource version: `8`
- Generation: `2`
- Lifecycle generation: `4`
- Scope: `project`
- Outcome: `open`
- Execution: `active`
- Active attempt: `—`

## Relationships

- Parent: —
- Depends on: `agent-system-phase-1`
- Supersedes: —

## Acceptance criteria

- `goal-recovery` (r1, required): Goal publication exposes stable committed and incomplete states, and interruption at every publication boundary recovers to the prior valid record or idempotently completes the intended record. — Evidence: Run publication-boundary fault injection, doctor, recovery, and validation tests against exact record identities.
- `legacy-migration` (r1, required): Legacy goal migration creates a fresh valid identity, retains source path, digest, and raw bytes, maps fields explicitly, and preserves unmapped prose without modifying the source. — Evidence: Run migration fixtures and inspect retained provenance and unmapped source evidence.
- `bounded-catalogs` (r1, required): Versioned bounded topology, policy, action, capability, workspace-check, and goal catalogs derive from owner-local facts with deterministic provenance, conflicts, bounds, and completeness checks. — Evidence: Regenerate checked JSON and Markdown, run deterministic and negative completeness fixtures, and compare source digests.
- `system-index` (r1, required): One bounded AgentSystemIndex contains only catalog identities, versions, input digests, conflicts, and query routes without duplicating catalog bodies. — Evidence: Validate the index schema, byte ceiling, input bindings, and rejection of embedded catalog payloads.
- `context-capsule` (r1, required): A bounded offline context command returns the path-, label-, or task-relevant situation with provenance, structured unavailable states, safe next discovery actions, and a human rendering from the same data. — Evidence: Run clean-root task scenarios with optional providers unavailable and compare JSON and human projections.
- `runtime-isolation` (r1, required): The fixed runtime control kernel and cross-process task isolation remain available under optional package failure, timeout, stale revision, and untrusted-package fault fixtures. — Evidence: Run package timeout, fault containment, namespace, lease, and stale publication tests.
- `resource-baseline` (r1, required): A criteria-revision-bound baseline records numeric ceilings for correctness, unsafe actions, calls, context size, cold and warm time, registered coverage, and reused checks without estimating unavailable measurements. — Evidence: Validate the baseline fixture and run its standing scenarios against the exact candidate.
- `exact-candidate-validation` (r1, required): All affected packages, schemas, generated projections, documentation, goal records, and delivery state pass focused validation against the exact delivered candidate. — Evidence: Inspect candidate-bound build, test, formatting, goal-validation, and delivery evidence.

## Recent attempts

- [`goal-publication-002`](attempts/goal-publication-002/) — `change`, `closed`, resource version `1`, criteria r2
- [`goal-publication-001`](attempts/goal-publication-001/) — `change`, `closed`, resource version `2`, criteria r2

## Record map

- [`goal.yaml`](goal.yaml): machine-authoritative goal state
- [`criteria.yaml`](criteria.yaml): versioned acceptance criteria
- [`attempts/`](attempts/): isolated attempt records and evidence

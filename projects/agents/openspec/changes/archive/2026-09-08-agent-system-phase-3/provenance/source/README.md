# Add advisory planning, admission, and reusable evidence

> Generated bounded projection. Edit `goal.yaml`, `criteria.yaml`, and attempt records only through the goal tool.

- Goal ID: `agent-system-phase-3`
- API version: `goals.alwaldend.com/v1alpha1`
- Resource version: `3`
- Generation: `2`
- Lifecycle generation: `1`
- Scope: `project`
- Outcome: `open`
- Execution: `active`
- Active attempt: `attempt-9a199d4db6bd`

## Relationships

- Parent: `repo-agent-system`
- Depends on: `agent-system-phase-2`
- Supersedes: —

## Acceptance criteria

- `criterion-001` (r1, required): impact-planning: A deterministic ImpactPlan resolves the smallest sufficient plan for a given intent: selected capabilities, required and forbidden effects, changed and reverse-affected targets at narrowest scope, minimum checks plus cover… — Evidence: Inspect linked evidence against the criterion.
- `criterion-002` (r1, required): effect-admission: Provider gateways admit action contracts against exact authority, subject, environment, pre-state, and budgets; read and hermetic compute stay cheap; remote writes use prepare/validate/authorize/execute/verify; unknown ac… — Evidence: Inspect linked evidence against the criterion.
- `criterion-003` (r1, required): validation-sets: Execution emits candidate-bound immutable ValidationSets with exact candidate, profile, check identities, sanitized arguments, working scope, provider/config/toolchain/policy digests, clean pre/post state, results, coverag… — Evidence: Inspect linked evidence against the criterion.
- `criterion-004` (r1, required): evidence-assertions: Goal and task owners create EvidenceAssertions applying one or more validation sets to an exact criterion revision and semantic verdict; delivery consumes validation sets and goals consume both; the same immutable set… — Evidence: Inspect linked evidence against the criterion.
- `criterion-005` (r1, required): receipt-applicability: Receipt applicability uses relevant dependency-slice digests; tree-bound checks survive message-only rewrites; commit-bound checks do not; changes to base, tree, config, policy, contract, toolchain, generator, enviro… — Evidence: Inspect linked evidence against the criterion.
- `criterion-006` (r1, required): full-check-structured: full-repo-check emits incremental versioned JSON plus generated Markdown, supports workspace and phase selectors, binds exact inputs and profile, records the target-universe count, rejects zero or unexpected reductio… — Evidence: Inspect linked evidence against the criterion.
- `criterion-007` (r1, required): resource-bounds: Do not persist raw BEP or environment dumps; batch compatible Bazel checks; deduplicate identical concurrent checks; cap receipts; benchmark the persistent Bazel server's warm-query latency and keep scratch, output-base, a… — Evidence: Inspect linked evidence against the criterion.
- `criterion-008` (r1, required): exact-candidate-validation: All affected packages, schemas, generated projections, documentation, goal records, and delivery state pass focused validation against the exact delivered candidate. — Evidence: Inspect linked evidence against the criterion.

## Recent attempts

- [`attempt-9a199d4db6bd`](attempts/attempt-9a199d4db6bd/) — `change`, `open`, resource version `1`, criteria r1

## Record map

- [`goal.yaml`](goal.yaml): machine-authoritative goal state
- [`criteria.yaml`](criteria.yaml): versioned acceptance criteria
- [`attempts/`](attempts/): isolated attempt records and evidence

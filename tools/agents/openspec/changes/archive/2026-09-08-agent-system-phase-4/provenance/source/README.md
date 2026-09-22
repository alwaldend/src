# Join durable work, delivery, review, version, and release

> Generated bounded projection. Edit `goal.yaml`, `criteria.yaml`, and attempt records only through the goal tool.

- Goal ID: `agent-system-phase-4`
- API version: `goals.alwaldend.com/v1alpha1`
- Resource version: `5`
- Generation: `2`
- Lifecycle generation: `2`
- Scope: `project`
- Outcome: `achieved`
- Execution: `paused`
- Active attempt: `—`

## Relationships

- Parent: —
- Depends on: —
- Supersedes: `agent-system-phase-3`

## Acceptance criteria

- `criterion-001` (r1, required): goal-resume: A fresh agent can discover every maintained open goal and resume one without prior path knowledge or free-form archaeology. — Evidence: Inspect linked evidence against the criterion.
- `criterion-002` (r1, required): delivery-validation: prepare -> publish cannot succeed with only a caller-asserted head; the supplied validation set must match the exact candidate and policy. — Evidence: Inspect linked evidence against the criterion.
- `criterion-003` (r1, required): review-traversal: A review outcome is traversable to the delivered fix and regression without making delivery the goal or test owner. — Evidence: Inspect linked evidence against the criterion.
- `criterion-004` (r1, required): release-identity: Formal release identity, bundle head, version/channel, immutable artifact names, remote refs, and deployment manifest agree; unrelated tags do not truncate changelogs. — Evidence: Inspect linked evidence against the criterion.
- `criterion-005` (r1, required): release-plan: The same reviewed release-ref plan can publish and verify a nightly tag or a release branch/tag pair; an existing immutable release tag never moves. — Evidence: Inspect linked evidence against the criterion.
- `criterion-006` (r1, required): interrupted-goal-publish: Interrupted multi-file goal publication has a supported deterministic doctor or recovery result. — Evidence: Inspect linked evidence against the criterion.

## Recent attempts

- [`phase4-implementation`](attempts/phase4-implementation/) — `change`, `closed`, resource version `3`, criteria r1

## Record map

- [`goal.yaml`](goal.yaml): machine-authoritative goal state
- [`criteria.yaml`](criteria.yaml): versioned acceptance criteria
- [`attempts/`](attempts/): isolated attempt records and evidence

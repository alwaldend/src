# Goal graph organization

Use this reference when one objective may need related goals, parallel work,
or dependency-aware resumption. The graph is a view over ordinary local Goal
resources. The execution harness remains responsible for scheduling and
permissions.

## Decide whether to decompose

Keep ordinary ordered steps in one attempt plan. Create a separate goal only
when the subgoal benefits from at least one independent property:

- lifecycle or resumption;
- acceptance criteria and evidence;
- ownership or retention;
- substantial parallel execution; or
- reuse by more than one objective.

The coordination cost of another goal and edge must be lower than the value of
making that boundary explicit.

## Use typed relationships

`Goal.spec.relationships` supports three same-catalog references:

- `parentGoalRef` expresses hierarchy, not execution order;
- `dependsOnGoalRefs` names prerequisites; and
- `supersedesGoalRefs` records lineage.

A parent is not a Kubernetes owner reference and never implies garbage
collection. Supersession does not redirect dependencies or authorize mutation
of the older goal. Only an `achieved` prerequisite satisfies a dependency;
an open prerequisite waits, while an abandoned or superseded prerequisite
blocks until the relationship is explicitly revised.

Keep each relationship kind acyclic independently. Do not test the union of
different edge kinds as one directed acyclic graph. Missing targets make a
catalog view unknown, rather than making one otherwise well-formed Goal
resource invalid.

## Change structure deliberately

Treat inferred relationships as proposals, never facts. Bind a proposal to
its source evidence and the goal generation it inspected. A coordinator must
review it, analyze the complete catalog for cycles, and publish it through the
normal resource-version compare-and-swap path. A worker must not rewrite
canonical relationships.

Relationship publication completely replaces the dependency and supersession
lists. Omitting the parent preserves it; clearing the parent is an explicit
operation. Close an active attempt before changing relationships: each attempt
is bound to the goal generation and portable goal-state digest that existed
when work began. An accepted relationship request advances generation and
resource version even when normalization makes it a semantic no-op.

Per-goal locks do not make a catalog transaction. Cycle prevention checks a
snapshot assembled by reading each goal under its own lock. Concurrent writes
to different goals can both pass against older snapshots and jointly create a
cycle. The single coordinator must avoid those conflicting writes; always run
the deterministic graph projection again after parallel structural changes and
repair any reported cycle.

Dispatch only goals whose prerequisites are satisfied. Agent selection,
capability discovery, delegation, communication topology, and runtime
scheduling belong to the execution harness.

Use the complete deterministic graph analysis for inspection. A catalog with
unresolved references is `Unknown`; a cyclic catalog is `Invalid`. Graph
analysis does not truncate its input, so report a catalog read or validation
failure instead of inferring readiness from a partial view.

Attempt input bindings, criteria and goal-state digests, structured reviews,
and artifact digests form the provenance graph for evidence.

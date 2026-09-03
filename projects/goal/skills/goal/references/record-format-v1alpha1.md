# Goal resource format v1alpha1

Use this reference to create, resume, inspect, or validate a goal. The files
use Kubernetes-style resource envelopes for familiar versioning and evolution;
they are portable local records, not Kubernetes objects or CRDs.

## Layout

The same directory shape is used in both scopes:

```text
<goal-id>/
  goal.yaml
  criteria.yaml
  criteria-revisions/
    <revision>.yaml
  README.md
  attempts/
    <attempt-id>/
      attempt.yaml
      plan.md
      result.md
      evidence/
```

Workspace records normally live at `out/<task>/goals/<goal-id>/`. The
recommended direct-initialization layout for project records is
`<owner-root>/goals/<goal-id>/`. This is guidance rather than a CLI constraint;
direct initialization, promotion, and migration accept other safe in-workspace
roots. A goal ID is stable across promotion and must not be reused.

## Resource envelope

Every canonical YAML resource uses the Kubernetes-style envelope. For example,
a newly initialized local Goal has this shape:

```yaml
apiVersion: goals.alwaldend.com/v1alpha1
kind: Goal
metadata:
  name: example-goal
  resourceVersion: "1"
  generation: 1
  creationTimestamp: "2026-08-30T00:00:00Z"
  annotations:
    goals.alwaldend.com/local-owner-root: out/example
spec:
  title: Example goal
  scope: workspace
  retention:
    policy: ephemeral
  relationships:
    dependsOnGoalRefs: []
    supersedesGoalRefs: []
status:
  lifecycleGeneration: 1
  outcome: open
  execution: active
  criteriaRevision: 1
  observedAt: "2026-08-30T00:00:00Z"
```

The allowed kinds are `Goal`, `GoalCriteria`, `GoalAttempt`, and
`GoalSessionBinding`. `metadata.name` is a DNS-compatible stable resource
identity. Use only standard Kubernetes `ObjectMeta` field names. In the local
store, `metadata.resourceVersion` emulates an API-server optimistic-concurrency
token: carry it literally and do not perform arithmetic on it.
The local CLI advances `metadata.generation` for spec changes and for every
accepted relationship replacement, including a normalized no-op. Lifecycle
generation is separate because criteria replacements and outcome or execution
changes can invalidate in-flight work without changing desired relationships.

## Authority and projections

`goal.yaml`, `criteria.yaml`, immutable `criteria-revisions/<revision>.yaml`
snapshots, each `attempt.yaml`, and the exact `plan.md`, `result.md`, and
`evidence/*.md` bytes bound by its artifact manifest are canonical. Validation
checks those Markdown sidecars against their SHA-256 digests. `README.md` is
generated, bounded, and replaceable. It is never canonical input; regenerate
it from the canonical record instead of editing it by hand.

Use the repository tool for canonical operations when it is available:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- <command> --help
```

The command owns IDs, resource validation, per-goal locking,
resource-version checks, atomic file replacement, path normalization, and
rendering. Do not emulate those mechanics with direct YAML edits. If the tool
is unavailable, keep exactly one writer and state that revision and
concurrency safety were not verified.

The local adapter keeps a deterministic, persistent per-goal lock file under
`$XDG_RUNTIME_DIR/alwaldend/goal/locks/`. It is outside the workspace and
cannot be committed. While held, it contains the current process PID for
diagnostics; clean release clears it, while a crash can leave a stale PID.
`flock`, never the PID text, is authoritative, and the tool never unlinks a
lock based on that text. The lock is not a resource or part of the portable
record format. A mutation writes a sibling temporary file, closes it, and
renames it over the destination while holding that goal's lock. Each file
replacement is atomic; a command that changes several files does not claim a
cross-file transaction.

## Goal

`Goal.spec` owns the requested configuration: title, `workspace` or `project`
scope, retention policy, and stable `parentGoalRef`, `dependsOnGoalRefs`, and
`supersedesGoalRefs` relationships. The workspace-relative owner root is
local adapter metadata in the `goals.alwaldend.com/local-owner-root`
annotation. The filesystem backend requires it, but portable goal-state
digests exclude it.
Each local goal reference is an object containing `name`, rather than a path or
bare ID. Relations do not grant permission to change related goals or imply a
scheduler in v1alpha1.

`Goal.status` owns observed state: outcome, execution, active and accepted
attempt IDs, accepted result digest, current criteria revision, lifecycle
generation, promotion or migration provenance, and update time. Criteria
replacement and accepted outcome or execution changes advance lifecycle
generation. Starting, updating, or closing an attempt without such a transition
changes goal state but does not advance lifecycle generation.

## Criteria

`GoalCriteria.spec` identifies its goal by a `{name: ...}` reference and a
revision. Each item has a stable criterion ID, its own revision, a required
flag, statement, and evidence method. Every accepted complete criteria
replacement advances the criteria resource version, generation, and spec
revision. An item's revision advances only when its meaning changes. Evidence
proves exactly one item revision; an earlier pass becomes historical rather
than current.

Keep a fixed regression set in the criteria resource. Run it after every
attempt that can affect it and against the frozen final result.

## Attempts

`GoalAttempt.spec` binds the work to a `{name: ...}` goal reference, goal and
lifecycle generations, criteria revision, portable criteria and goal-state
digests, optional `planID`, and work type. It does not store a Goal resource version: that value
is the checkpoint caller's mutation-time compare-and-swap token, not a durable
attempt input. Attempt status records whether it is open or closed, relevant
timestamps, the SHA-256 artifact manifest, and the structured close review.
Review evidence references may name only the manifest's plan, result, or
evidence paths. The digest-bound Markdown bytes live beside the resource.

An attempt may carry optional structured resume fields so a fresh agent can
resume an open goal without free-form archaeology:

- `stableDefect`: the reproducible problem under investigation.
- `hypothesis`: the candidate explanation being tested.
- `subject`: the exact system, artifact, or reference under test.
- `affectedCriteria`: criterion IDs this attempt exercises; unique and sorted.
- `regressionRefs`: reviewed regression set or fixtures; unique and sorted.
- `priorAttemptID`: an earlier attempt this one resumes or corrects.
- `dominantFailure`: the single most useful failure signal observed.
- `measurableDelta`: the measured difference this attempt produces.
- `nextAction`: the deterministic next step for a resuming agent.
- `blocker`: an external condition preventing progress, if any.
- `resumeCondition`: what must hold before a resuming agent resumes this
  attempt.

## Plans

`Goal.status.plans` stores durable plan summaries. Each entry has a portable
`planID`, a bounded `strategy`, and a state of `active`, `accepted`,
`rejected`, or `superseded`. A rejected plan carries a bounded
`rejectionReason`. At most one plan is active. Create or transition a plan
with `goal checkpoint --plan-id ... --plan-strategy ...` or
`--plan-state ... --plan-only`; a new plan supersedes the previous active plan.
An attempt may bind to one of these summaries with `spec.planID`; the plan is
input context and the attempt review remains the acceptance evidence.

These fields are advisory input, not acceptance evidence: a closed attempt may
omit them entirely. Prose fields must be trimmed, bounded, and free of NUL;
identifiers must be portable record IDs; lists must be unique, sorted, and
bounded. The generated README projection surfaces `stableDefect`,
`dominantFailure`, `nextAction`, `blocker`, and `resumeCondition` for open
attempts so an agent can resume directly from the goal record.

For the registered repository goals root, use the goal command's catalog-backed
resume view instead of scanning goal directories:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- resume \
  --goals-root projects/agents/goals
```

The output is a bounded `GoalResumePacket` decoded from the checked,
digest-verified goal catalog. It contains only open goals with a resumable
open attempt and their exact candidate paths and continuation fields. Override
`--catalog` only when a task-specific goals root has its own generated goal
catalog; otherwise the command deliberately refuses stale or invalid catalogs.

`Goal.status.acceptedResultDigest` is the accepted attempt's `result.md`
SHA-256, not automatically the identity of an external deliverable. Record an
external subject's exact identity in the result or evidence when acceptance
depends on it.

Allowed work types are `investigation`, `candidate`, `change`, `integration`,
`validation`, and `decision`. Not every work unit has a file
candidate or hypothesis. Stateful work should identify source revision, target
and environment/config revision, plus an operation receipt.

An open attempt may receive isolated evidence. Closing it publishes its final
resource, plan, result, and evidence set. After close, treat its directory as
immutable; a correction or additional evidence is a new attempt that refers to
the old one.

## Bounds and paths

List, goal-detail, and README views have explicit item or byte limits and
report truncation. Graph analysis instead consumes the complete accepted
catalog or fails. Prefer useful recent entries in bounded views over failing
because history is large. Canonical and migrated Markdown files must be
bounded regular files containing valid UTF-8 and no NUL bytes. Store links
relative to the workspace or owner root, or as stable external URLs. Never
place local absolute paths, credentials, or private environment values in a
project goal.

## Kubernetes compatibility boundary

The resources deliberately follow Kubernetes API conventions, but the local
files are not directly managed by a Kubernetes API server. The filesystem
adapter supplies `resourceVersion`, `generation`, and local annotations. A
future Kubernetes adapter must let the API server populate its read-only
metadata, install structural OpenAPI schemas, enable the status subresource,
and reconcile or strip local-only annotations during import and export.

Keep portable identity and desired state in the resource. Store a session's
workspace-relative goal path only in the namespaced annotation
`goals.alwaldend.com/local-goal-ref`; `GoalSessionBinding.spec.goalRef` remains
a normal `{name: ...}` object. Digest-bound plan, result, and evidence sidecars
are a local storage representation; a cluster adapter must map their bytes to
cluster-accessible objects or artifact references.

This boundary lets a future CRD/controller reuse the API types. Do not claim
that a local record can be passed to `kubectl apply` without that adapter and
the CRDs.

When implementing that adapter, verify the current Kubernetes
[custom-resource](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
and [ObjectMeta](https://kubernetes.io/docs/reference/kubernetes-api/definitions/object-meta-v1-meta/)
contracts rather than copying the local store's metadata behavior blindly.

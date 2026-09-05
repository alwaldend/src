---
title: Goal command
description: CLI for versioned local goal records
languages:
  - go
tags:
  - agent
  - workflow
---

`goal` is the command surface for the repository goal skill. Portable resource
types and domain validation live in `api/v1alpha1`; deterministic local
persistence lives in `internal/fsstore`. The command keeps authoritative
Kubernetes-inspired YAML envelopes and digest-bound attempt Markdown separate
from the bounded generated `README.md` projection, uses the same record format
for workspace and project goals, and stores session focus separately under
ignored workspace scratch.

The record layout is:

```text
goals/<goal-id>/
  goal.yaml
  criteria.yaml
  criteria-revisions/<revision>.yaml
  README.md
  attempts/<attempt-id>/
    attempt.yaml
    plan.md
    result.md
    evidence/
```

`<owner-root>/goals/<goal-id>/` is the preferred project placement, not a CLI
constraint. Direct initialization, promotion, and migration also accept other
safe in-workspace goals roots while recording the explicit owner root.

The YAML objects use `apiVersion: goals.alwaldend.com/v1alpha1` and the kinds
`Goal`, `GoalCriteria`, `GoalAttempt`, and `GoalSessionBinding`. Their
envelopes validate Kubernetes-compatible qualified metadata keys and label
values and use object-shaped local goal references. The shared API treats
`resourceVersion` as opaque. The filesystem backend allocates canonical local
numeric versions and requires complete persisted metadata. The files are not
raw `kubectl` input, and this project supplies neither CRDs nor a controller.

Run it through Bazel:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- init \
  --goals-root out/example/goals \
  --goal-id verify-the-release \
  --title "Verify the release" \
  --criterion "All affected tests pass"
```

Use a task-specific binding directory when changing session focus:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- attach \
  --session-root out/example/goal-sessions \
  --session-id current \
  --goal-dir out/example/goals/verify-the-release
```

`GoalSessionBinding.spec.goalRef` contains the stable object name. The
namespaced annotation `goals.alwaldend.com/local-goal-ref` contains the
workspace-relative storage path. Session bindings never create a repository-
global catalog.

`Goal.metadata.annotations[goals.alwaldend.com/local-owner-root]` contains the
workspace-relative owner used only by the local adapter. It is deliberately
absent from portable `Goal.spec` and portable digests.

Commands are `init`, `list`, `resume`, `show`, `attach`, `checkpoint`,
`learning-proposal`, `graph`, `set-relationships`, `validate`, `promote`,
`render`, and `migrate`.
`set-relationships` replaces the complete dependency and supersession lists;
the parent is preserved unless explicitly set or cleared. It requires no
active attempt, advances Goal generation and resource version on every accepted
request, and refreshes the bounded README projection. Its cycle check is a
per-goal-locked catalog snapshot rather than a catalog-wide transaction, so a
`graph` call after concurrent writes settle is authoritative.

`resume` requires `--catalog`, resolved relative to `--goals-root`, and prints
a bounded `GoalResumePacket` containing only open goals with a resumable
attempt. The packet carries candidate paths and structured continuation
fields; it never mutates records or opens goal plans, results, or evidence.
The strict decoder validates catalog structure and its stored digest, not
freshness against current goal records. Use `show --goal-dir` for current
local state.

`learning-proposal` validates one Phase 5 proposal. Repeated friction references
are required, but the command never promotes the proposal or edits source; the
owning project adopts it through ordinary review and delivery.

`checkpoint` starts or publishes an attempt, changes execution/outcome state,
or applies a complete desired criteria-items file with `--criteria-file`.

For ordinary work, initialize the objective and acceptance once with `init
--title ... --criterion ...`, then use an inline checkpoint:

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- checkpoint \
  --goal-dir out/example/goals/fix-parser \
  --expected-resource-version 1 \
  --subject "candidate digest or source revision" \
  --summary "Parser repaired; focused checks remain unverified." \
  --next-action "Run the parser acceptance check."
bazel_agent bazel run //projects/goal/cmd/goal -- show \
  --goal-dir out/example/goals/fix-parser
```

`--summary` is nonblank UTF-8 Markdown, limited to 8192 bytes, and requires
`--subject` and `--next-action`. It cannot combine with `--plan-file` or
`--result-file`. The store generates the initial plan from the objective,
criteria, and next action, then writes progress into the existing canonical
`result.md`. Later checkpoints reuse the open attempt, preserve its initial
plan and imported evidence, and update the declared subject and next action.
Changing an existing subject requires a replacement summary or result;
retained evidence applies only to the candidate identified in that evidence.
A next-action-only checkpoint is also supported.

`show` includes `activeAttempt` with current continuation fields, an 8192-byte
UTF-8 result preview, complete result byte count/digest, evidence digests,
observation time, and source paths relative to the goal directory. Its
`resultTruncated` field indicates when the canonical result needs a separate
read. Subject and progress are caller declarations, not live Git observations
or criterion verdicts. Initialization, summary checkpoints, and bounded views
use the same store and recovery protocol as detailed attempts; there is no
separate compact record format. Closing still requires an explicit
`--close-attempt --review-file` with criterion verdicts and evidence; a
summary does not imply acceptance.

Detailed attempts can provide separate plan, result, and review payloads.
Plans are durable summaries in `Goal.status.plans`. Create one with
`--plan-id ... --plan-strategy ... --plan-only`; transition the active plan
with `--plan-id ... --plan-state accepted|rejected|superseded --plan-only`,
adding `--plan-rejection-reason` for a rejection. A newly created active plan
supersedes the prior active plan. Attempts may bind to a plan with
`--plan-id` when starting or by using the active plan of an existing attempt.
Criteria updates require paused execution and use ordinary atomic file
replacement. Immutable criteria snapshots retain the exact canonical criteria
revision for historical attempts; portable, domain-separated criteria and
goal-state digests avoid binding durable records to local resource-version
tokens. Closing an attempt requires `--review-file` with an `accept`, `refine`,
or `reset` decision and per-criterion verdicts linked to frozen plan, result,
or evidence artifacts.
The review file is plain YAML with exactly two keys. Criteria entries must be
sorted by `criterionID`, and `evidenceRefs` must be unique, sorted, and name
only frozen artifacts (`plan.md`, `result.md`, or files under `evidence/`).
Verdicts are `pass`, `fail`, or `unverified`; a non-`unverified` verdict
requires at least one evidence reference. Do not add `apiVersion` or `kind`
headers:

```yaml
decision: accept
criteria:
  - criterionID: friction-baseline
    criterionRevision: 1
    verdict: pass
    evidenceRefs:
      - evidence/friction-baseline.md
  - criterionID: optimization-pipeline
    criterionRevision: 1
    verdict: pass
    evidenceRefs:
      - evidence/optimization-pipeline.md
```

An achieved outcome must close an accepting attempt whose exact passes cover
every current required criterion. Structured verdicts are kept in
`attempt.yaml`; richer narrative stays in `result.md`.

For attempt and lifecycle checkpoints, a new attempt is fully staged before
`goal.yaml` advances its resource version. That goal write is the optimistic-
concurrency commit point; canonical attempt content follows, and `README.md`
is last. An immediately closed new attempt uses an intermediate active pointer
until its directory is published and the Goal is finalized at the same
resource version, so interruption at either gap fails validation. The store
returns the committed Goal reference after any post-commit failure, and the CLI
error identifies that resource version.

Promotion also requires a paused workspace goal and preserves the goal name,
portable input bindings, immutable criteria history, and content-digest
provenance. It rejects known absolute workspace/file links in promoted attempt
artifacts; callers remain responsible for semantic privacy and credential
review.

Migration is a non-destructive import, not an in-place conversion. The legacy
source remains unchanged, while a complete validated record is published at
`<destination-goals-root>/<goal-id>` with one final directory rename. Repeating
the same import is idempotent only while source provenance and import options
match the existing target.

```sh
bazel_agent bazel run //projects/goal/cmd/goal -- migrate \
  --source-goal-dir out/example/legacy/verify-the-release \
  --destination-goals-root out/example/imported/goals
```

Mutations to an existing Goal's canonical state take an exact expected local
resource version. Cooperating processes serialize them with an advisory lock
keyed to the canonical goal path under
`$XDG_RUNTIME_DIR/alwaldend/goal/locks/`, reread the manifest while holding the
lock, and replace individual files by renaming sibling temporary files. Lock
files are outside the workspace and cannot enter version control. Different
goals do not share a lock. Local owner and session-link annotations are
workspace-relative normalized paths; absolute host paths are never written to
manifests or command output.

Each individual file replacement is atomic; commands that update several files
do not claim cross-file transaction semantics. `README.md` is a derived
projection and is written last. Plan, result, and evidence Markdown is instead
canonical through the SHA-256 artifact manifest in `attempt.yaml`. Direct file
writers are outside this cooperative trust boundary, and digest validation
makes later edits detectable. Execution and outcome transitions advance
`status.lifecycleGeneration`, invalidating in-flight attempt input bindings.

This experimental `v1alpha1` format is a local file protocol. One coordinator
should own a goal record. Delegated workers write only to isolated scratch;
the coordinator imports selected artifacts through `checkpoint`.

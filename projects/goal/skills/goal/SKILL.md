---
name: goal
description: >-
  Persist and pursue a concrete multi-step research or implementation goal
  through versioned records, evidence-backed work units, honest acceptance
  checks, and strategy-changing retries. Use when iterative work benefits from
  state that survives a turn; do not use for a simple one-step task or a
  question answerable in one response.
---

# Goal

Keep the objective, decisions, evidence, and failed approaches in an
inspectable record. The skill owns the execution protocol; the goal tool owns
record identity, revisions, validation, projections, and safe file updates.
A session may focus on one goal while the store contains many.

Within the request's authority, safety constraints, and acceptance
obligations, the requested result is the priority. Records, process changes,
tools, and delegation exist only to help produce and verify it. Do not let
process optimization or record maintenance displace the next
highest-leverage work on the result.

## Start or resume the right record

A substantial research task can be a goal even when it only reads source
material: the distinction is iterative work and persisted evidence, not
whether product files change. Follow the request's authority boundaries. A
question alone does not authorize a durable project record or unrelated
mutations.

Classify the record before work:

- Use a **workspace goal** for current-task coordination, investigations, and
  experiments. Store it under the repository's ignored task output, normally
  `out/<task>/goals/<goal-id>/`.
- Use a **project goal** only when the user explicitly requests, or the project
  already establishes, a maintained record. Prefer the narrowest owning root
  and the layout `<owner-root>/goals/<goal-id>/`. This placement is guidance,
  not a CLI constraint; direct initialization, promotion, and migration accept
  other safe in-workspace roots.

Both scopes use the same format. Never make a record durable merely because a
task became long. Never key persisted state only by a thread or session ID.

## Preserve the goal across interruptions

Treat an active goal's lifecycle separately from the priority of the current
turn. A question, status request, clarification, or additional task may need
to run first, but it does not implicitly pause, abandon, supersede, complete,
or otherwise terminate the goal. Answer or perform the interruption within
its own authority, preserve any in-flight state safely, and then resume the
highest-leverage ready goal work in the same turn when practical or on the
next available turn.

A question grants no new mutation authority, and an additional task does not
broaden the goal's scope. Resume only work already authorized by the goal and
reconcile explicit changes to its objective or acceptance criteria through the
normal versioned workflow. Stop pursuing the goal only when it is achieved,
the user explicitly pauses, cancels, abandons, or supersedes it, or its
execution is honestly waiting or blocked under the lifecycle rules. Ending a
response, answering a question, servicing a higher-priority task, or receiving
no repeated “continue” instruction is not such a stop signal.

Before working, create or explicitly attach the goal, then read its manifest,
current criteria revision, active attempt, and bounded projection. If the
record is unversioned, preserve it and follow the migration procedure; do not
silently reinterpret prose as structured fact.

Read [record-format-v1alpha1.md](references/record-format-v1alpha1.md)
whenever creating, resuming, inspecting, or validating a record.

## Run an evidence-backed work unit

Read [lifecycle-and-evidence.md](references/lifecycle-and-evidence.md) before
starting an attempt.

1. Define observable acceptance criteria and how each will be checked. Mark
   inferred criteria as assumptions. Never weaken the requested outcome to
   make it pass.
2. Start one bounded work unit against exact goal generation, lifecycle,
   criteria, and portable digest bindings. Keep the current Goal resource
   version separately as the next checkpoint's compare-and-swap token. Record
   the work type, target uncertainty or defect, plan, inputs, and intended
   evidence before implementation.
3. Preflight cheap deterministic assumptions before expensive work. A rejected
   preflight can revise an unfrozen plan; after work or publication starts, a
   material change is a new attempt.
4. Prefer a small, well-defined, high-leverage module over a broad uncertain
   rewrite. Include the minimum whole-result context needed to judge its fit,
   interfaces, and regressions; local polish that harms the whole does not
   count as progress.
5. Produce the smallest artifact that can reliably test the plan, then inspect
   the actual result. Close the work unit early when decisive evidence
   falsifies its module, interface, or approach. Successful commands prove
   execution, not acceptance.
6. Evaluate every affected criterion as `pass`, `fail`, or `unverified`, run
   the fixed regression set, and bind evidence to the exact criteria revision
   and candidate or operation identity.
7. Close the attempt with its result and decision: accept, refine, or reset.
   Closed attempt directories are immutable.
8. Checkpoint canonical state using the resource version from the most recent
   canonical read. On a stale resource version or lifecycle generation, reread
   and reconcile; never overwrite newer state.

If a defect survives two attempts, change the diagnosis, hypothesis, workflow,
or validation method. Change the representation only when bounded local
evidence establishes that its owner or interface cannot express the required
result without regression. Use a stable defect name so recurrence cannot be
hidden by relabeling. Preserve failed attempts and rejected evidence.

### Prefer evidenced deltas before replacement

Start from the last accepted or frozen artifact when one exists and remains a
viable baseline. Otherwise use the best viable current artifact and record why
no accepted baseline is available. Before replacing a whole module,
representation, or architecture, identify one acceptance-visible defect and
the evidence that exposes it. Find the smallest connected owner of that
defect, make one reversible change there, and compare the before and after
states with the same acceptance-relevant views, tests, or observations. Keep
the change only on clear evidence of improvement without a material
regression; undo or decline to promote ambiguous changes, no-ops, and
regressions.

Before spending a work unit, run a lightweight causal-reach check. Trace how
the proposed owner and operation can change the named acceptance evidence,
compare its largest safe effect with the measured gap, and check that it can
produce the required category of result rather than merely nearby activity.
Cancel or redirect a task whose bounded influence is immaterial or whose
interface cannot express the required cue. Keep the check proportional and
prefer existing evidence or one cheap probe over a new planning detour. When
new evidence disproves reachability during implementation, return early at
that evidence boundary instead of finishing the obsolete plan.

Numerical or mechanical movement is not result progress by itself. More
topology, code, parameters, completed commands, attempts, or measured activity
counts only when acceptance-relevant evidence shows that the requested result
improved. Permit new topology, architecture, representation, or a whole-module
replacement only after several bounded local edits establish with evidence
that the present owner or interface has a structural limit. Record that limit
and the discriminating test for the replacement before starting it.

### Stage feedback before promotion

When a cheap noncanonical artifact can expose the dominant uncertainty, do not
make promotion-grade machinery a prerequisite for seeing that artifact. First
preserve the immutable source, authorization boundary, and safety invariants;
then produce and inspect a disposable probe in the minimum whole-result context
needed to falsify the hypothesis. Give it an exact identity and label it
clearly as diagnostic and non-promotable. A preview, dry run, integration
coupon, benchmark sample, or partial render is evidence about the hypothesis,
not evidence that the requested result is acceptable.

Escalate validation only as the artifact earns it:

1. Run the cheap checks necessary to create the diagnostic safely.
2. Inspect the observable result against the dominant failure and stop
   immediately on a categorical miss.
3. Freeze interfaces and invest in exhaustive correctness, regression,
   reproducibility, and promotion checks only for a plausible survivor.
4. Rerun the complete acceptance plan against the exact promotable result.

Do not defer a safety, authorization, destructive-action, or source-preservation
gate merely to obtain faster feedback. Conversely, do not treat speculative
schemas, validators, builders, dashboards, or reports as progress when a safe
observable probe would settle the current decision sooner. Contracts should
describe a representation that survived the relevant early evidence, not force
an untested hypothesis to satisfy prematurely frozen mechanics.

For long or repeatedly stalled loops, set and record a task-local target for
time to the next decision-bearing artifact. Base it on measured prior latency
and the cost of the work rather than a universal duration. If the target is
missed, stop expanding the work unit, account for the delay, and either emit the
smallest safe probe or reset the approach. Keep independent preparation in
parallel only when it can still affect that next decision.

### Put correction on the critical path

Do not rely on the user to notice stagnation, and do not postpone all process
judgment until an attempt closes. At every attempt preflight and when its first
decision-bearing artifact arrives, evaluate the correction gate. Open one
review episode when any of these conditions holds:

- the same stable defect survived two closed attempts;
- the latest acceptance-visible artifact is flat or regressed against the
  frozen immediate comparison or the best relevant historical comparison;
- the recorded artifact-time target was missed for a controllable reason;
- an unresolved acceptance-material assumption can change the next decision
  and lacks a cheap deterministic discriminator; or
- the user identifies a material defect or stagnation that the running process
  should have detected.

Coalesce overlapping triggers into one episode bound to the current decision;
resuming a session does not open a duplicate. Read
[correction-gate.md](references/correction-gate.md) whenever an episode opens
or remains unresolved. Resolve its reviewer verdict before expanding beyond
the first decision-bearing artifact, taking a costly or irreversible step,
closing or promoting the attempt, or starting the next attempt.

Before implementation, still record the strongest recurrence risk and the
first artifact that distinguishes success from recurrence. Review that exact
artifact as soon as it exists; do not wait for an archival packet or user
complaint. Record attempt start, worker start, first-artifact time, and
decision time when latency matters. The reviewer supplies disconfirming
evidence and a stop signal, not mutation or publication authority; the
coordinator owns canonical state and the final evidence-backed disposition.

Separate invariant execution scaffolding from the changing hypothesis. Reuse
a proven runner, fixture, renderer, comparator, or test harness across bounded
variants and change only the payload under test. If one authoring modality
repeatedly produces the same qualitative defect despite parameter changes,
switch to a more directly controllable modality rather than writing a larger
generator for the same abstraction. Examples include moving from batch
generation to an interactive edit loop, from prose-only design to an
executable coupon, or from end-state inspection to continuous fixed-reference
comparison. Preserve reproducibility around the surviving result; do not make
reproducibility machinery a substitute for producing it.

## Leave inspectable progress

Every user turn that advances an active goal must create or materially update
an inspectable artifact. Examples include a report, test result, comparison,
candidate, deterministic script, or attempt result. A progress message or
cosmetic record edit is not enough.

Link the artifact in commentary or the final response with a short label and
honest verdict. Record it through the current attempt so the bounded README
projection can surface it. Put disposable outputs under the repository's
ignored task directory. Acceptance-critical evidence for a project goal must
instead be retained by a repository-relative link, stable public reference, or
reproducible regeneration instructions.

## Preserve durable progress remotely

Before yielding a turn that advances an active goal and leaves task-owned
durable repository changes, use the repository's delivery workflow to run
proportional checks, stage only task-owned paths or hunks, commit or amend
according to branch policy, and push the resulting commit. This applies to
accepted work, rejected attempts, and useful early returns. Verify that the
remote feature ref contains the exact local commit; an unpushed local commit
is not a completed remote checkpoint.

Commit project goal records, deterministic sources, reusable tools, and
required deliverables needed to resume or reproduce the work. Keep workspace
goal records and disposable outputs ignored, do not stage unrelated or
user-owned changes, and do not create empty or cosmetic commits merely to
satisfy this rule. Preserve decision-relevant findings from ignored artifacts
in durable records or reproducible inputs when the task requires retention.

If validation, authentication, a changed remote ref, conflicts, branch policy,
or unavailable delivery transport prevents a safe commit or push, do not claim
that progress is remote. Preserve what can safely be preserved, record the
exact blocker and recovery action in the active attempt, and surface it to the
user.

## Keep lifecycle axes separate

- `outcome`: `open`, `achieved`, `abandoned`, or `superseded`.
- `execution`: `active`, `paused`, `waiting`, or `blocked` while the outcome is
  open.
- `scope`: `workspace` or `project`.
- `retention`: how raw and acceptance-critical evidence survives.

Set `achieved` only after freezing one result and rerunning the full evidence
plan against that exact content or operation identity. `unverified` is not a
pass. Set `blocked` only when progress needs unavailable authority, input, or
external state; ordinary difficulty and incomplete work remain active.

Before each next attempt, audit what measurably improved, regressed, or stayed
unchanged; the highest-leverage unresolved issue; whether the approach should
continue, change, or be discarded; and where avoidable feedback time went.
Turn that audit into a live control for the next attempt when recurrence risk
is material; a retrospective alone is insufficient. Optimize for cheap
falsification without weakening final acceptance.

At least every three closed attempts, and immediately when progress stalls,
review the full attempt history and the whole delivery process more deeply
than the per-attempt audit. Reconsider decomposition, interfaces, tools,
evidence quality, feedback latency, and how this protocol is being applied.
Change only what evidence indicates will move the requested result faster or
more reliably. Record the review and the last attempt it covers in attempt
evidence or the attempt result so a resumed session can determine when another
review is due. Keep process changes task-local and reversible unless broader
change is separately authorized.

## Coordinate safely

Use one coordinator as the canonical writer for a goal. Workers may produce
isolated attempt evidence against immutable input versions, but they must not
publish canonical state.

At the start or resumption of every long-running goal, and before each new
attempt, enumerate the ready independent workstreams. When two or more can
produce useful, independently reviewable outputs concurrently, actively
delegate them unless a concrete dependency, shared-state risk, or immediate
review-and-integration bottleneck makes sequential work faster. Record that
reason in the attempt plan or result when remaining sequential. This is the
default execution check, not an optional optimization.

Good worker boundaries include bounded candidate variants, reference or
source analysis, feasibility research, implementation of disjoint modules,
implementation-blind review, verification, and artifact preparation. Treat
agents, context, and coordination as real costs: do not fragment trivial
operations merely to occupy slots, and do not recursively delegate without a
separate benefit. Prefer a few meaningful workstreams over one worker per file
or command.

Schedule already-authorized task-local agents, compute, and reviews around the
critical path. Limit work in progress to what can be evaluated and integrated
promptly; unused capacity is preferable to speculative work that creates a
review queue. Stop or redirect safely interruptible work when it can no longer
affect the next decision or required acceptance and delivery evidence.

Repeat the parallelism check whenever an attempt closes, stalls, changes
strategy, or exposes a new independent workstream.

When a shared artifact has no stable independently mergeable interfaces, do
not split the canonical artifact merely to enable parallel edits. Test bounded
changes on isolated copies of the exact frozen baseline, approved when the
task requires approval. Review each copy in the minimum whole-result context
needed to detect integration regressions, and promote only candidate-bound
evidence permitted by the task's acceptance and approval policy. Promotion
remains subject to existing authority to mutate the canonical artifact; this
workflow grants none. Send only inputs permitted for each worker or tool
boundary, and minimize ancillary context.

Read [sessions-and-concurrency.md](references/sessions-and-concurrency.md)
before resuming in a new session, switching among goals, using parallel
workers, or recovering from a stale update.

Read [graph-organization.md](references/graph-organization.md) before splitting
an objective into related goals, changing goal relationships, or dispatching
dependency-aware work. Keep ordinary linear steps inside one attempt.

Read [promotion-and-migration.md](references/promotion-and-migration.md) before
promoting a workspace goal, importing an unversioned record, or retaining
evidence in a public project goal.

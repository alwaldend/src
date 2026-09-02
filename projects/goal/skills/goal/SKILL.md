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

Before authoring a candidate, review the goal from its purpose outward. Ask
what object or outcome is intended, how it is made or used in reality, which
physical or operational model best represents it, which acceptance failures
that model implies, and whether alternative representations would test the
goal more directly. Do not rush into local implementation merely because an
existing artifact invites editing. If repeated attempts fail, reopen this
model review before choosing another localized variation.

1. Define observable acceptance criteria and how each will be checked. Mark
   inferred criteria as assumptions. Never weaken the requested outcome to
   make it pass.
2. Represent the approach as a durable plan. Create or select a plan with a
   portable `planID` and bounded `strategy` before starting an attempt. At
   most one plan is active; creating a replacement supersedes the previous
   active plan, while a rejected plan records why it cannot continue.
3. Start one bounded work unit against exact goal generation, lifecycle,
   criteria, portable digest bindings, and the plan ID when one exists. Keep
   the current Goal resource version separately as the next checkpoint's
   compare-and-swap token. Record the work type, target uncertainty or defect,
   inputs, and intended evidence before implementation.
4. Preflight cheap deterministic assumptions before expensive work. A rejected
   preflight can revise an unfrozen plan; after work or publication starts, a
   material change is a new attempt, and a rejected approach is a rejected
   plan rather than a hidden retry.
5. Prefer a small, well-defined, high-leverage module over a broad uncertain
   rewrite. Include the minimum whole-result context needed to judge its fit,
   interfaces, and regressions; local polish that harms the whole does not
   count as progress.
6. Produce the smallest artifact that can reliably test the plan, then inspect
   the actual result. Close the work unit early when decisive evidence
   falsifies its module, interface, or approach. Successful commands prove
   execution, not acceptance.
7. Evaluate every affected criterion as `pass`, `fail`, or `unverified`, run
   the fixed regression set, and bind evidence to the exact criteria revision
   and candidate or operation identity.
8. Close the attempt with its result and decision: accept, refine, or reset.
   Close the plan as accepted, rejected, or superseded when its approach is
   settled. Closed attempt directories are immutable.
9. Checkpoint canonical state using the resource version from the most recent
   canonical read. On a stale resource version or lifecycle generation, reread
   and reconcile; never overwrite newer state.

If a defect survives two attempts, change the hypothesis, representation,
workflow, or validation method. Use a stable defect name so recurrence cannot
be hidden by relabeling. Preserve failed attempts and rejected evidence.

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
Optimize for cheap falsification without weakening final acceptance.

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

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
4. Produce the smallest artifact that can reliably test the plan, then inspect
   the actual result. Successful commands prove execution, not acceptance.
5. Evaluate every affected criterion as `pass`, `fail`, or `unverified`, run
   the fixed regression set, and bind evidence to the exact criteria revision
   and candidate or operation identity.
6. Close the attempt with its result and decision: accept, refine, or reset.
   Closed attempt directories are immutable.
7. Checkpoint canonical state using the resource version from the most recent
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

## Coordinate safely

Use one coordinator as the canonical writer for a goal. Workers may produce
isolated attempt evidence against immutable input versions, but they must not
publish canonical state. Treat agents, context, and coordination as real
costs; delegate only substantial independent work whose quality or
critical-path benefit justifies them.

Read [sessions-and-concurrency.md](references/sessions-and-concurrency.md)
before resuming in a new session, switching among goals, using parallel
workers, or recovering from a stale update.

Read [graph-organization.md](references/graph-organization.md) before splitting
an objective into related goals, changing goal relationships, or dispatching
dependency-aware work. Keep ordinary linear steps inside one attempt.

Read [promotion-and-migration.md](references/promotion-and-migration.md) before
promoting a workspace goal, importing an unversioned record, or retaining
evidence in a public project goal.

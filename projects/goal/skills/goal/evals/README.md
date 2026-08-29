---
title: Goal evaluations
---

# Goal evaluations

This suite records the semantic contract for persistent, evidence-backed goal
pursuit. The required offline Bazel target validates the Promptfoo
configuration, referenced cases, and staged skill without making a model call.
The cases cover research routing, workspace-versus-project ownership, explicit
session focus, stale-update reconciliation, honest acceptance,
result-prioritized modular work, isolated candidate promotion, and
critical-path delegation. Concurrency cases require active delegation when a
long-running goal exposes multiple independently reviewable workstreams,
require a recorded reason for sequential execution, and reject fanout whose
only purpose is occupying available slots. The interruption case distinguishes
turn priority from lifecycle state: questions and additional tasks do not
silently stop an already-authorized active goal or expand its authority.
The staged-feedback case requires a safely isolated, explicitly non-promotable
observable probe before promotion machinery when that probe can cheaply
falsify the dominant uncertainty, while preserving all safety and final
acceptance gates. The standing-correction case requires an adversarial reviewer
and a first-artifact veto on the live critical path, reuse of invariant
execution scaffolding, measured feedback latency, and a change of authoring
modality when repeated batch generation preserves the same qualitative defect.
The correction-episode cases require a Medium coordinator to detect recurrence
and regression without a user diagnosis, route exactly one context-isolated
Ultra reviewer, gate further authoring on its compact verdict, suppress
duplicate reviewers across overlapping triggers and resume, and honor a
`STOP` result unless cited contrary evidence supports an override. They also
require a fresh same-tier reviewer when the coordinator already uses the
strongest tier.
The delta-first cases require one reversible change to the smallest connected
owner of an observable defect, fixed before/after acceptance evidence, and
rejection of activity metrics as result progress. They permit a broader
replacement only after several local edits establish an evidenced structural
limit and a discriminating test for the new approach. The causal-reach case
requires cancelling even an in-flight bounded edit when its maximum safe
effect is immaterial to the measured gap or cannot express the required output
category.

A manual `//projects/goal/skills/goal:eval` target exercises the staged skill
with a Medium subject and an isolated stronger judge. Run it only when live,
credentialed semantic evaluation is intended:

```sh
bazel_agent test //projects/goal/skills/goal:eval \
  --test_env=CODEX_HOME=/absolute/path/to/.codex \
  --test_env=CODEX_PATH_OVERRIDE=/absolute/path/to/codex
```

That target judges the response contract; it does not prove that T3's
collaboration tool spawned a reviewer, preserve a correction episode across
turns, expose the spawn's model, effort, or fork arguments, or prevent a later
mutation. The isolated Promptfoo provider does not currently expose that
longitudinal orchestration surface. A complete tool-trace eval still needs a
writable multi-turn fixture that captures collaboration calls and injects a
deterministic reviewer result. Promptfoo validation proves only that these
assets load. Deterministic store behavior is covered by
`//projects/goal/internal/fsstore:go_test`; a future
longitudinal eval remains a separate harness capability, not an inferred pass.

The static storage case covers the ordinary one-goal boundary: one path-keyed
goal lock plus atomic per-file rename. Promotion and migration use a separate
two-path protocol that acquires distinct source and destination locks in
canonical-path order.

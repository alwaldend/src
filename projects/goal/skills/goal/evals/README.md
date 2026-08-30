---
title: Goal evaluations
---

# Goal evaluations

This suite records the semantic contract for persistent, evidence-backed goal
pursuit. The required offline Bazel target validates the Promptfoo
configuration, referenced cases, and staged skill without making a model call.
The cases cover research routing, workspace-versus-project ownership, explicit
session focus, stale-update reconciliation, honest acceptance, and bounded
delegation.

A live target is omitted because representative behavior spans multiple turns
and requires filesystem tools, safe writable fixtures, and fresh-session
resume. A tool-free single response cannot verify those longitudinal
postconditions. Promptfoo validation therefore proves only that these assets
load. Deterministic store behavior is covered by
`//projects/goal/internal/fsstore:go_test`; a future
model eval still needs an isolated multi-turn workspace fixture.

The static storage case covers the ordinary one-goal boundary: one path-keyed
goal lock plus atomic per-file rename. Promotion and migration use a separate
two-path protocol that acquires distinct source and destination locks in
canonical-path order.

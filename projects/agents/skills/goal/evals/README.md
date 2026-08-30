---
title: Goal evaluations
---

# Goal evaluations

This suite records the behavioral contract for persistent, evidence-backed
goal pursuit. The required offline Bazel target validates the Promptfoo
configuration, referenced cases, and staged skill without making a model call.
The cases also record that ordinary task-coordination goals remain in the
repository's ignored temporary-output area rather than becoming source-tree
churn, and that agent delegation must justify its coordination and compute
cost rather than merely produce a small speedup.

A live target is omitted because representative behavior spans multiple turns
and requires creating and reviewing real artifacts, maintaining goal records,
and changing strategy after measured failures. A tool-free single response
cannot verify those filesystem and longitudinal postconditions. Promptfoo
validation therefore proves only that these evaluation assets load; behavior
needs an isolated multi-turn workspace fixture before a live target would be
meaningful.

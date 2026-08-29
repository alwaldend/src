---
title: Goal evaluations
---

# Goal evaluations

This suite records the behavioral contract for persistent, evidence-backed
goal pursuit. The required offline Bazel target validates the Promptfoo
configuration, referenced case, and staged skill without making a model call.

A live target is omitted because representative behavior spans multiple turns
and requires creating and reviewing real artifacts, maintaining goal records,
and changing strategy after measured failures. A tool-free single response
cannot verify those filesystem and longitudinal postconditions. Promptfoo
validation therefore proves only that these evaluation assets load; behavior
needs an isolated multi-turn workspace fixture before a live target would be
meaningful.

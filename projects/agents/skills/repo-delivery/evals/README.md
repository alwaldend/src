---
title: Repository delivery evaluations
---

# Repository delivery evaluations

This suite describes the behavioral contract for safely finalizing a feature
branch and its pull request. Its required offline Bazel target validates the
Promptfoo configuration, referenced case, and staged skill without making a
model call.

A live target is omitted because representative delivery requires inspecting
and changing Git history, pushing a remote ref, interacting with a forge, and
responding to live review state. A read-only, tool-free model evaluation cannot
reproduce those operations safely or verify their postconditions. Promptfoo
validation therefore proves only that the evaluation assets load; delivery
behavior still needs isolated integration fixtures or review of a real,
authorized delivery.

---
title: Task workspace evaluations
---

# Task workspace evaluations

This suite records the behavioral contract for isolating a task workspace
before the first mutation or task-scratch write. The required offline Bazel
target validates the Promptfoo configuration, referenced case, and staged skill
without making a model call.

A live target is omitted because representative behavior requires a mutable
multi-worktree Git checkout, existing and concurrent task branches, shared
system temporary directories, and a filesystem that can be inspected for
residue. A read-only response cannot prove those placements or the absence of
leftover files. Promptfoo validation therefore proves only that these
evaluation assets load; behavior needs isolated worktree and temporary-directory
fixtures before a live target would be meaningful.

---
title: Remote-rebase evaluations
---

# Remote-rebase evaluations

This suite records the behavioral contract for safely rebasing a feature
branch and updating its remote ref. The required offline Bazel target validates
the Promptfoo configuration, referenced case, and staged skill without making
a model call.

A live target is omitted because representative behavior requires a mutable
Git repository and remote, concurrent-ref-change fixtures, history rewriting,
and a push whose postconditions can be inspected. A read-only response cannot
prove those operations. Promptfoo validation therefore proves only that these
evaluation assets load; behavior needs isolated local remotes and controlled
race fixtures before a live target would be meaningful.

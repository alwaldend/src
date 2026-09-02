---
title: Ast-grep evaluations
---

# Ast-grep evaluations

This suite describes the structure-aware search and rewrite contract for the
repository-pinned ast-grep toolchain. Its required offline Bazel target
validates the Promptfoo configuration, referenced case, and staged skill
without making a model call.

A live target is omitted because representative behavior requires invoking
the pinned ast-grep binary against the repository, inspecting matched AST
nodes, and optionally generating and reviewing a rewrite diff. Those tool and
side-effect interactions cannot be reproduced by the read-only Promptfoo
subject. Configuration validation is not evidence that a pattern was applied
or rewritten correctly.

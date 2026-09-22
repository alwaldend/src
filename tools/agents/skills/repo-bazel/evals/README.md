---
title: Repository Bazel evaluations
---

# Repository Bazel evaluations

This suite describes the required behavior for a focused Bazel maintenance
change. Its offline Bazel target validates the Promptfoo configuration,
referenced case, and skill staging without making a model call. The
configuration names `openai:codex-sdk` in read-only mode, but this suite is
intended only for offline validation and does not invoke it.

A live target is omitted because representative behavior requires inspecting
and editing BUILD files, querying the build graph, and invoking Bazel and
Buildifier in a real workspace. Configuration validation does not establish
that the generated build graph is valid or that the requested tests pass.

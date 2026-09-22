---
title: Nested Bazel module evaluations
---

# Nested Bazel module evaluations

This suite describes the required behavior for creating and integrating a
standalone Bzlmod workspace. Its offline Bazel target validates the Promptfoo
configuration, referenced case, and skill staging without making a model call.
The configuration names `openai:codex-sdk` in read-only mode, but this suite
is intended only for offline validation and does not invoke it.

A live target is omitted because representative behavior requires creating
workspace files, updating module locks, running Gazelle, and invoking Bazel in
both the nested and root workspaces. Configuration validation does not prove
that those integrations or generated files are correct.

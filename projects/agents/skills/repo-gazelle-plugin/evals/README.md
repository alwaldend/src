---
title: Repository Gazelle plugin evaluations
---

# Repository Gazelle plugin evaluations

This suite describes the required design and integration behavior for a
repository Gazelle plugin. Its offline target validates the Promptfoo
configuration, referenced case, and staged skill without making a model call.

A live target is omitted because representative behavior requires creating a
nested module, updating locks and BUILD files, invoking Gazelle and Bazel, and
reviewing generated diffs. Promptfoo validation is not evidence that an agent
performs those tool calls correctly; the plugin's Go and integration tests
provide deterministic coverage of the generated behavior.

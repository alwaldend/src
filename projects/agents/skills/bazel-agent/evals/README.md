---
title: Bazel agent evaluations
---

# Bazel agent evaluations

This suite describes the required behavior for invoking and troubleshooting
the repository's Bazel runner. Its offline Bazel target validates the
Promptfoo configuration, referenced cases, and skill staging without making a
model call. The configuration names `openai:codex-sdk` in read-only mode, but
this suite is intended only for offline validation and does not invoke it.

Cached-tool cases distinguish confirmed support from a known older runner,
requiring reuse of the observation and the supported Bazel target fallback
without an unauthorized host update.

A live target is omitted because representative behavior requires checking the
host runner and invoking Bazel in a real workspace. Promptfoo configuration
validation is not evidence that an agent uses the runner correctly or
diagnoses a Bazel failure accurately.

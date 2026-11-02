---
title: Bazel agent evaluations
---

# Bazel agent evaluations

This suite describes the required behavior for invoking and troubleshooting
the repository's Bazel workflow. Its offline Bazel target validates the
Promptfoo configuration, referenced cases, and skill staging without making a
model call. The configuration names `openai:codex-sdk` in read-only mode, but
this suite is intended only for offline validation and does not invoke it.

Installation cases require explicit authority and a source-built AL executable.
They reject private executable caches and unauthorized host updates.

Affected-test cases cover an explicit Git base and bounded target patterns,
nested-workspace selection, conservative full-scope fallback, and graph-query
failures. They keep ordinary Bazel commands valid and distinguish a successful
empty selection from an unsuccessful query. The selector does not certify
delivery readiness or claim implicit coverage of other workspaces.

A live target is omitted because representative behavior requires checking the
host configuration and invoking Bazel in a real workspace. Promptfoo configuration
validation is not evidence that an agent invokes Bazel correctly or
diagnoses a Bazel failure accurately.

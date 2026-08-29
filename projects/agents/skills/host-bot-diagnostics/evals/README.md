---
title: Host bot diagnostics evaluations
---

# Host bot diagnostics evaluations

This suite describes the required behavior for a safe host-bot audit. Its
offline Bazel target validates the Promptfoo configuration, referenced case,
and skill staging without making a model call. The configuration names
`openai:codex-sdk` in read-only mode, but this suite is intended only for
offline validation and does not invoke it.

A live target is omitted because representative behavior requires read-only
access to host services, listeners, firewall state, Codex diagnostics, and
deployed configuration. Configuration validation does not show that an agent
interprets that state correctly or avoids exposing local secrets.

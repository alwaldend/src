---
title: Repository Ansible evaluations
---

# Repository Ansible evaluations

This suite describes the required behavior for safely changing packaged
Ansible automation. Its offline Bazel target validates the Promptfoo
configuration, referenced case, and skill staging without making a model call.
The configuration names `openai:codex-sdk` in read-only mode, but this suite
is intended only for offline validation and does not invoke it.

A live target is omitted because representative behavior requires editing a
fixture deployment and invoking its Bazel packaging and syntax checks.
Configuration validation is not evidence that the automation is idempotent,
that injected credentials remain protected, or that no live inventory was
contacted.

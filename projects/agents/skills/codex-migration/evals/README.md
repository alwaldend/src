---
title: Codex migration evaluations
---

# Codex migration evaluations

This suite describes safe migration behavior for provider, model,
authentication, and shared-configuration changes. Its offline Bazel target
validates the Promptfoo configuration, referenced case, and skill staging
without making a model call.

A live target is omitted because representative behavior requires access to a
host secret manager and two provider requests, including a Codex tool round
trip. Configuration validation is not evidence that authentication works,
that provider-bound encrypted history is portable, or that a migration is safe
for concurrent sessions.

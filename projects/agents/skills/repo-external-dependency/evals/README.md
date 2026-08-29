---
title: External dependency evaluations
---

# External dependency evaluations

This suite describes the reproducibility and ownership contract for updating a
shared external dependency. Its required offline Bazel target validates the
Promptfoo configuration, referenced case, and staged skill without making a
model call.

A live target is omitted because representative behavior requires inspecting
repository consumers, querying an upstream registry or release service,
verifying downloaded artifacts, updating locks, and invoking Bazel against the
changed graph. Those tool and network interactions cannot be reproduced by the
read-only Promptfoo subject. Configuration validation is not evidence that a
dependency was selected, pinned, or verified correctly.

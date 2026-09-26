---
title: External dependency evaluations
---

# External dependency evaluations

This suite describes reproducibility, ownership, and publisher selection for
external dependencies. Cases distinguish reputable upstream or distribution
sources from obscure repackagers without treating star counts as a trust
threshold. Its required offline Bazel target validates the Promptfoo
configuration, referenced cases, and staged skill without making a model call.

A live target is omitted because representative behavior requires inspecting
repository consumers, querying an upstream registry or release service,
verifying downloaded artifacts, updating locks, and invoking Bazel against the
changed graph. Those tool and network interactions cannot be reproduced by the
read-only Promptfoo subject. Configuration validation is not evidence that a
dependency was selected, pinned, or verified correctly.

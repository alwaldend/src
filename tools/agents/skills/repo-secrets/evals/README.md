---
title: Repository secrets evaluations
---

# Repository secrets evaluations

This suite describes the trust-boundary and injection contract for a secret
needed by repository-managed infrastructure. Its required offline Bazel target
validates the Promptfoo configuration, referenced case, and staged skill
without making a model call.

A live target is omitted because representative behavior requires access to
repository-specific `al.lua` wiring, Vault authentication and policy state,
and possibly sensitive Terraform or deployment tooling. Supplying those tools
or credentials to a model evaluation would be unsafe and non-reproducible.
Offline validation does not demonstrate correct secret handling; it only proves
that the public, placeholder-only evaluation assets load.

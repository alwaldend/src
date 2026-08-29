---
title: Repository Terraform evaluations
---

# Repository Terraform evaluations

This suite describes the safe implementation and validation contract for a
Terraform resource rename. Its required offline Bazel target validates the
Promptfoo configuration, referenced case, and staged skill without making a
model call.

A live target is omitted because representative behavior requires inspecting a
real component, invoking Bazel and Terraform wrappers, obtaining Vault and
provider access for a useful plan, and reasoning about existing remote state.
The read-only Promptfoo subject has none of that tool or state surface. Offline
validation does not show that a plan is safe or that a resource will avoid
replacement; those claims require an isolated state fixture or authorized plan
review.

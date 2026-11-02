---
title: Project layout evaluations
---

# Project layout evaluations

This suite describes role-based placement across repository boundaries,
`internal` naming, the treatment of legacy layouts, and landing-page onboarding
for new first-party projects including repository-only helpers. Its offline Bazel
target validates the Promptfoo configuration, referenced cases, and skill
staging without making a model call.

A live target is omitted because the representative requests require an agent
to inspect the surrounding repository before choosing an owner and path.
Promptfoo configuration validation is not evidence that an agent makes those
decisions correctly.

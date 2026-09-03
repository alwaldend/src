---
title: Repository agent goals
description: Durable goals, attempts, and evidence for the repository agent system
tags:
  - agent
  - goals
---

# Repository agent goals

Each child directory is an owner-local durable goal record. Canonical resource
files and attempt artifacts are created and validated through the repository
`goal` workflow; do not edit generated records or projections directly.

Goal evidence may contain source and operational analysis, including
non-secret live or generated facts. It must never contain credentials, other
secrets, or personal information. Inspect raw artifacts before promoting their
content.
Structured task artifacts such as FrictionRecord and LearningProposal JSON
files are task scratch: keep them under `out/<task>/` and cite them from goal
evidence Markdown rather than committing them into `evidence/`, which accepts
only `.md` files.

The current open implementation goal is under
`agent-system-phase-six-ergonomics/`. It measures agent friction across
representative task archetypes and adopts optimizations through the Phase 5
learning-proposal contract. The maintained, link-safe public-system synthesis
is the [current-state document](../docs/current-state.md).
Context-bound raw audits remain in canonical goal records rather than the
documentation projection.

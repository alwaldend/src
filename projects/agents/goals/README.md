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

Discover current open work through the goal catalog and its bounded resume
view; status is owned by each goal record. The
`agent-system-phase-six-ergonomics/` goal records a completed friction-reduction
attempt, with the original evidence and its limitations preserved.

The [current-state document](../docs/current-state.md) describes supported
interfaces and current evidence boundaries. Closed attempts remain historical
records; later implementation reviews may narrow claims they made without
rewriting those immutable records.

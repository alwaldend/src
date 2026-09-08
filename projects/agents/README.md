---
title: Agents
description: Repository-wide agent-system architecture and skills
statuses:
  - active
languages:
  - markdown
  - bzl
tags:
  - agent
  - skills
---

# Repository agent system

This project owns the repository-wide agent-system contract and reusable
cross-repository skills. Maintained improvement work lives in this project's
[OpenSpec workspace](openspec/). It does not centralize
component facts or runtime state: each fact remains canonical at its natural
owner, and system-wide views are derived projections.

## Start here

| Document                                                                   | Purpose                                                  |
| -------------------------------------------------------------------------- | -------------------------------------------------------- |
| [Current state](docs/current-state.md)                                     | Evidence-backed baseline and material seams              |
| [Architecture](docs/architecture.md)                                       | Canonical abstraction tower, authorities, and invariants |
| [Roadmap](docs/roadmap.md)                                                 | Dependency-ordered future work and acceptance signals    |
| [OpenSpec changes](openspec/changes/)                                      | Maintained work, acceptance, and preserved goal history  |
| [Root agent guide](https://github.com/alwaldend/src/blob/master/AGENTS.md) | Current repository-wide operating policy                 |

The current-state document describes supported entry points and their evidence
boundaries. Dated audits remain in [migrated history](../../infra/src/openspec/migration.md).
The architecture defines
the intended composition contract; the roadmap does not claim that proposed
interfaces already exist.

## Ownership boundaries

- `docs/` owns the cross-layer system model and plan, not duplicated component
  configuration.
- `openspec/` owns agent-system requirements, changes, and evidence;
  [the migration map](../../infra/src/openspec/migration.md) locates historical records.
- Other components own their specifications and changes under their own
  `openspec/` directory. `infra/src/openspec/` describes repository evolution.
- `skills/` owns reusable repository-wide agent procedures and their
  development-time evaluations.
- Product-specific skills remain with their project at
  `projects/<project>/skills/<name>`.
- Repository-internal executors and build integrations may live under
  `tools/`; their behavior remains owned and documented there.

## Skill packaging and discovery

Every canonical skill is packaged as a `skill_library` in its owning
directory.

The repository discovery directory `.agents/skills/` contains one relative
symlink per skill. Each link points directly to its canonical project-owned
directory. Bazel ignores the discovery directory and builds only canonical
targets, preventing duplicate packages while allowing skills from more than
one owning project.

Skill evaluation data is not part of the runtime `skill_library` unless a
skill explicitly declares otherwise. Every new or updated skill includes an
offline Promptfoo configuration-validation target. That target checks the eval
harness, not answer correctness. Live behavioral evaluations are manual,
networked tests and must not be included in ordinary wildcard test runs; they
may be omitted when representative coverage requires tool calls or external
state that cannot be provided safely and reproducibly, with the gap documented
beside the eval configuration.

---
title: Agents
description: Repository-wide agent-system architecture, goals, and skills
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

This project owns the repository-wide agent-system contract, its durable
improvement goals, and reusable cross-repository skills. It does not centralize
component facts or runtime state: each fact remains canonical at its natural
owner, and system-wide views are derived projections.

## Start here

| Document                                                                   | Purpose                                                  |
| -------------------------------------------------------------------------- | -------------------------------------------------------- |
| [Current state](docs/current-state.md)                                     | Evidence-backed baseline and material seams              |
| [Architecture](docs/architecture.md)                                       | Canonical abstraction tower, authorities, and invariants |
| [Roadmap](docs/roadmap.md)                                                 | Dependency-ordered future work and acceptance signals    |
| [Durable goals](goals/)                                                    | Versioned attempts, evidence, and acceptance state       |
| [Root agent guide](https://github.com/alwaldend/src/blob/master/AGENTS.md) | Current repository-wide operating policy                 |

The current-state document is a dated snapshot. The architecture defines the
intended composition contract. The roadmap is not a claim that proposed
interfaces already exist.

## Ownership boundaries

- `docs/` owns the cross-layer system model and plan, not duplicated component
  configuration.
- `goals/` owns durable repository-agent work records and evidence.
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

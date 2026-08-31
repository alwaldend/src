---
title: Src
description: Monorepo source and repository control system
---

# Source repository

This monorepo contains first-party projects, repository tooling,
infrastructure definitions, data, and pinned external sources. Facts remain
with their owning component; the repository agent system connects those facts
without replacing their authorities.

## Start here

- [Agent policy](AGENTS.md) is the repository-wide operating contract. A
  nearer `AGENTS.md` takes precedence within its subtree.
- [Agent-system current state](projects/agents/docs/current-state.md) records
  the evidence-backed baseline.
- [Agent-system architecture](projects/agents/docs/architecture.md) defines
  how intent, policy, capabilities, work, execution, evidence, delivery, and
  learning compose.
- [Agent-system roadmap](projects/agents/docs/roadmap.md) distinguishes
  proposed work from current guarantees.
- [Durable agent-system goals](projects/agents/goals/README.md) retain
  versioned attempts, acceptance state, and evidence provenance.
- [Repository setup](https://alwaldend.com/docs/misc/repo/) explains the human
  development environment.

The intended control flow is:

```text
request -> owner and policy -> capability -> work -> execution
        -> evidence -> delivery -> reviewed learning
```

## Repository shape

- [`projects/`](projects/README.md): products and reusable project code.
- [`tools/`](tools/README.md): repository-wide rules, toolchains, and
  repository-internal automation.
- [`infra/`](infra/README.md): infrastructure definitions and operations.
- [`data/`](data/README.md): repository-owned data and documentation assets.
- [`third_party/`](third_party/README.md): vendored and externally sourced
  inputs.
- [`users/`](users/README.md): user-owned code and infrastructure.
- [`projects/agents/`](projects/agents/README.md): repository-wide agent
  architecture, goals, skills, and evaluation assets.

Each linked tree README is authoritative for that tree. Bazel `BUILD` and
`MODULE.bazel` files remain authoritative for executable and dependency
structure; generated maps and catalogs are projections, not new sources of
truth.

## External links

- Homepage: https://alwaldend.com/
- Docs: https://alwaldend.com/docs/

## License

[AGPL-3.0](https://spdx.org/licenses/AGPL-3.0-or-later.html), see
[LICENSE.txt](data/license/LICENSE.txt).

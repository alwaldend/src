---
title: Tools
description: Tools tree
weight: 4
cascade:
  - categories:
      - tool
---

This tree contains tools that are used inside the repo

## Requirements

- MUST NOT be public (except toolchain types)
- MUST NOT be published
- MUST NOT be used as dependencies of production build targets
- MUST be available to the whole repo (`visibility = ["//:__subpackages__"]`)
- MAY be used in tests and explicit source-generation/update targets

An explicit source-generation/update target is a developer workflow, not part
of the production build graph. It may depend on a tool from this tree and use
`write_source_file` to update a checked-in source artifact. Normal production
and documentation targets must consume that checked-in artifact directly;
they must not depend on the generator, its tool, or an action-generated copy.

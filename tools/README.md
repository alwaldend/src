---
title: Tools
description: Tools tree
weight: 4
cascade:
  - categories:
      - tool
---

This tree contains tools used inside the repository. All tracked content
follows the repository's public-source policy.

## Requirements

- Bazel targets MUST use repository-internal visibility except for toolchain
  types whose owner explicitly exposes a public build API.
- Tool artifacts MUST NOT be published as first-party product artifacts.
- Tool targets MUST NOT be used as dependencies of production build targets.
- Tool targets intended for repository-wide use MUST use
  `visibility = ["//:__subpackages__"]`.
- Tool targets MAY be used in tests and explicit source-generation/update
  targets.

An explicit source-generation/update target is a developer workflow, not part
of the production build graph. It may depend on a tool from this tree and use
`write_source_file` to update a checked-in source artifact. Normal production
and documentation targets must consume that checked-in artifact directly;
they must not depend on the generator, its tool, or an action-generated copy.

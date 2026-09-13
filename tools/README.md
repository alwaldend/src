---
title: Tools
description: Tools tree
weight: 4
cascade:
  - categories:
      - tool
---

This tree contains developer tools and reusable Bazel rule modules. All tracked
content follows the repository's public-source policy. Standalone `rules_*`
modules retain their own Bzlmod names and public build APIs, with documentation
on the main repository site. They do not own separate landing sites.

## Requirements

- Bazel targets MUST use repository-internal visibility except for standalone
  Bazel rule modules and toolchain types whose owner explicitly exposes a public
  build API.
- Tool artifacts MUST NOT be published as first-party product artifacts.
- Standalone Bazel rule modules MAY publish their reusable build APIs through
  their explicit module release workflow.
- Production build targets MAY consume the public build APIs of standalone
  Bazel rule modules and toolchains. Other tool targets MUST NOT be dependencies
  of production build targets.
- Tool targets intended for repository-wide use MUST use
  `visibility = ["//:__subpackages__"]`.
- Tool targets MAY be used in tests and explicit source-generation/update
  targets.

An explicit source-generation/update target is a developer workflow, not part
of the production build graph. It may depend on a tool from this tree and use
`write_source_file` to update a checked-in source artifact. Normal production
and documentation targets must consume that checked-in artifact directly;
they must not depend on the generator, its tool, or an action-generated copy.

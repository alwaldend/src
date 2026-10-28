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
- Production build targets MAY load Starlark rules and use execution tools and
  toolchains, including the public build APIs of standalone Bazel rule modules.
  Other tool targets MUST NOT be ordinary source, library, or runtime dependencies
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

## Transitive dependency enforcement

`//tools:transitive_consumers` records the packages permitted to consume this
tree. Gazelle selects package groups using the path mappings in
[`transitive_visibility/transitive_visibility.toml`](transitive_visibility/transitive_visibility.toml).
The listed project exceptions preserve documentation, tests, source updates,
and explicit developer/release commands. An exception permits the entire
consumer package, including other targets sharing it; it does not prove that
every target in that package follows the production rule above. New exceptions
require review of the consuming workflow. Standalone modules are separate
repositories and retain their public build APIs.

Public toolchain-type packages and shared build settings explicitly select
`//tools:build_api_consumers`. This preserves external toolchain registration
and repository build configuration without granting other tool packages the
same exception.

The intended artifact boundary excludes execution-tool and Starlark-load
edges. Bazel 9.3.1's native transitive check cannot express that distinction
for every attribute; see the [generator's enforcement limitations](../projects/bazel_gazelle_enforce_transitive_visibility/README.md#enforcement-limitation-in-bazel-931).
Use `--notool_deps` for queries of ordinary target dependencies. A package
exception also permits ordinary dependencies from that package, so it is not
equivalent to an exemption for execution tools alone.

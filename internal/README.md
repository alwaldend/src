---
title: Internal libraries
description: Shared first-party runtime libraries for this repository
---

This tree contains shared first-party runtime libraries used within this
repository. All tracked content follows the repository's public-source policy.

Each library owns its documentation, specifications and maintained changes
under its component directory. A component used by only one project belongs
in that project's own internal directory.

## Requirements

- Bazel targets MUST use repository-internal visibility.
- Internal artifacts MUST NOT be independently published as product artifacts.
- Internal targets MAY be used by production build targets and MAY be runtime
  dependencies of repository applications.

These libraries do not own separate product landing pages or release workflows.
See the [tools](../tools/README.md) and [third-party](../third_party/README.md)
boundaries for developer tooling and externally maintained dependencies.

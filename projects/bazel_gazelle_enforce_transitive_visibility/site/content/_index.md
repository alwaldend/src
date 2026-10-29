---
title: Bazel Gazelle transitive visibility
linkTitle: Gazelle transitive visibility
description: Maintain Bazel dependency boundaries with TOML path mappings
layout: landing
statuses:
  - active
languages:
  - go
tags:
  - bazel
---

This Gazelle extension generates Bazel's package-level transitive visibility
declarations from a central TOML configuration. Directory paths select package
groups, and the most specific matching path controls descendant packages.

Register the extension after other Gazelle languages, declare the policy, and
enable Bazel's experimental transitive visibility enforcement. The restriction
uses package-granularity exceptions. Direct dependencies are checked; the
documentation records an indirect enforcement limitation in Bazel 9.3.1.

See the [documentation](/docs/projects/bazel_gazelle_enforce_transitive_visibility/)
for setup, examples, and limitations.

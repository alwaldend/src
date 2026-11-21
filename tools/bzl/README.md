---
title: Bzl
description: Bazel rules related to Bazel itself
languages:
  - bzl
tags:
  - bzl_rules
  - bzl_registry
---

## First-party ruleset archives

Develop rulesets from the monorepo root using their package paths, for example
`bazel test //tools/rules_docs/...`. Each project declares its standalone name
in its own `release.MODULE.bazel`.

To package the source archives, run:

```sh
bazel run //tools/bzl/cmd/release_archives -- --output-directory "$PWD/out/releases"
```

The command packages the committed source at `HEAD`, translates package labels
for standalone consumers, removes the enclosing monorepo visibility directive
from OpenSpec packages, and installs each `release.MODULE.bazel` as
`MODULE.bazel` inside its archive. Use those
manifests and archive checksums to prepare registry entries. DNSControl,
ISO, and systemd provide their own `release_archive` build targets.

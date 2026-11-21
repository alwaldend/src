---
title: OpenSpec Gazelle
description: Gazelle extension for owner-local OpenSpec workspaces
---

The OpenSpec extension generates `openspec_workspace(name = "source")` and
`openspec_validate_test(name = "validate_test")` in directories named
`openspec` that contain `config.yaml`. The workspace's source glob stops at
Bazel package boundaries. The generator preserves existing attributes and
leaves packages with conflicting target names alone.

The repository's Gazelle binary includes the extension. A custom binary can
include the public language library in `//tools/rules_openspec_gazelle/gazelle`,
or use the plugin-only `:gazelle_openspec` binary. The extension loads
`@rules_openspec//rules:defs.bzl` using the consumer's Bzlmod repository mapping.
The root invocation supplies the equivalent in-tree label.

---
title: Bazel dependency analyzer
linkTitle: Bazel dependency analyzer
description: Turn Bazel dependency graphs into Mermaid diagrams
layout: landing
statuses:
  - active
languages:
  - go
  - bzl
tags:
  - bazel
  - bzl_rules
---

Generate dependency diagrams directly from a Bazel target. The analyzer reads
Bazel's JSON target graph and produces a Mermaid flowchart for tracing
dependencies or a mindmap for browsing package ownership. Both views retain
target identities and rule kinds, with optional grouping by package.

Limit dependency depth, include files or external targets, and render the
result with Mermaid. The same graph produces consistent output regardless of
the order of its input records.

See the [rule and command documentation](/docs/projects/bazel_deps_analyzer/)
for usage and a generated example.

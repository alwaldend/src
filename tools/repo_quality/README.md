---
title: Repository quality
description: Whole-index formatting, linting and validation coordinator
languages:
  - go
  - sh
tags:
  - formatter
  - linter
  - validation
---

This package uses upstream `aspect_rules_lint` formatting and linting rules to
discover Git-tracked files and inspect declared source targets with established
tools acquired by Bazel. It does not implement a formatter or linter and does
not download tools at runtime.

Format all safe, hand-maintained files:

```sh
bazel_agent run //:format
```

Run the non-mutating repository check:

```sh
bazel_agent test //:repo_quality_test
```

Run semantic source linters over declared targets:

```sh
bazel_agent build --config=lint //...
```

The integration test intentionally uses the upstream rule's `no_sandbox` mode
so it can inspect files that are tracked by Git but not declared in BUILD
files. This is a local-checkout guarantee, like the root Buildifier test; tool
acquisition remains hermetic and pinned through Bazel.

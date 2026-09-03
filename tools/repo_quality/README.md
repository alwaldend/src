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
not download tools at runtime. Thin workspace adapters extend the upstream
language list with the official Go SDK for module and workspace manifests,
StyLua and Selene for Lua, and the existing Prettier XML plugin for Qt `.ui`
and `.qrc` files.

Format all safe, hand-maintained files:

```sh
bazel_agent bazel run //:format
```

Run the non-mutating repository check:

```sh
bazel_agent bazel test //:repo_quality_test
```

Run semantic source linters over declared targets:

```sh
bazel_agent bazel build --config=lint //...
```

The CSS formatter bucket includes CSS, Less and SCSS. The JavaScript bucket
includes JavaScript, JSON, JSON-with-comments, JSON5, TypeScript, TSX and Vue.
Go source uses gofumpt; `go.mod` and `go.work` use the official Go SDK.
Lua checks cover every tracked file not excluded by Git attributes, using
StyLua for canonical formatting and Selene for correctness diagnostics.

Explicit templates (`.j2`, `.tmpl`, `.tpl`, and the Hugo layout tree), generic
application configuration, and non-Terraform HCL are deliberately excluded
from general formatting. Ordinary Markdown, YAML, HTML and other files remain
covered even when their contents include template expressions. Generated,
vendored, lock, binary, diagram and exact-content files likewise retain their
owning regeneration, validation or integrity workflow.

The integration test intentionally uses the upstream rule's `no_sandbox` mode
so it can inspect files that are tracked by Git but not declared in BUILD
files. This is a local-checkout guarantee, like the root Buildifier test; tool
acquisition remains hermetic and pinned through Bazel.

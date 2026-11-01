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
al bazel run //:format
```

Run the non-mutating repository check:

```sh
al bazel test //:repo_quality_test
```

The suite also parses root `BUILD.bazel` and rejects `load()` statements to keep
unrelated tool dependencies out of root package loading. Run this check alone
with `al bazel test //tools/repo_quality/test/root_build:root_build_test`.

Run semantic source linters over declared targets:

```sh
al bazel build --config=lint //...
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

The quality suite also runs pinned ast-grep against every tracked first-party
Python file not marked `linguist-vendored` or `linguist-vendored=true`, rejecting symbol imports in
favor of module imports and qualified member access. Formatting exclusions do
not exempt first-party Python from this check. Explicit public re-exports
listed in `__all__` may use direct symbol imports in `__init__.py`; wildcard
imports remain forbidden.
Its syntax fixtures distinguish real imports from
comments and strings. Module from-imports are verified against source paths
and pinned Python module specifications without executing imported packages.
Like formatting, this complete-checkout check is local
and uncached; downloaded external repositories are outside its coverage.

## Owner inventories

`al bazel run //tools/repo_quality:inventory_update` derives the central
OpenSpec validation map from declared owner source exports, resolving standalone
module identities from their MODULE files. It also derives the project index
from project README metadata. Draft-only specification directories without a
source export are not registered. Landing-page membership remains owned by
`projects/projects.bzl`, independently of the documentation index.

The inventory freshness test in the quality suite inspects the actual checkout
and rejects stale projections, including newly declared exports. Like the other
whole-checkout checks it is local and uncached. Its isolated command fixture
checks regeneration and preserves the resulting inventory in test outputs.

---
title: Rules docs
description: Bazel documentation packaging rules and Gazelle extension
languages:
  - bzl
  - go
tags:
  - bzl_rules
---

`rules_docs` packages Markdown documentation under a common archive prefix.
Its Gazelle extension adds a `docs_filegroup` to existing Bazel packages that
contain a `README.md`.

## Documentation rule

Add the module dependency:

```starlark
bazel_dep(name = "rules_docs", version = "<VERSION>")
```

Then declare the package documentation. `srcs` defaults to
`glob(["*.md"])`; `visibility` is optional.

```starlark
load("@rules_docs//docs:defs.bzl", "docs_filegroup")

docs_filegroup(
    name = "docs",
    deps = ["child"],
)
```

Relative dependency names without a colon are normalized to the child
package's `docs` target. `prefix` defaults to the current package beneath
`content/docs/`.

## Gazelle extension

Add the public language target to the repository's custom Gazelle binary:

```starlark
load("@gazelle//:def.bzl", "gazelle_binary")

gazelle_binary(
    name = "gazelle_binary",
    languages = [
        "@rules_docs//gazelle",
    ],
)
```

The extension only operates in directories where both a BUILD file and a
`README.md` already exist. Newly generated rules use `glob(["*.md"])` and are
visible only to their nearest ancestor Bazel package. Existing `srcs`, `deps`,
`visibility`, and `prefix` attributes are preserved.

The public `@rules_docs//gazelle:gazelle_docs` binary contains only this
extension and can update a nested workspace without initializing unrelated
language plugins from its parent repository.

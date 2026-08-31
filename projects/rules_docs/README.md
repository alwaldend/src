---
title: Rules docs
description: Bazel documentation packaging rules
statuses:
  - active
languages:
  - bzl
tags:
  - bzl_rules
---

`rules_docs` packages Markdown documentation under a common archive prefix.

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

Generation support is packaged separately so consumers of the documentation
rule do not inherit Gazelle and Go dependencies. Add it as a development-only
module dependency:

```starlark
bazel_dep(
    name = "rules_docs_gazelle",
    version = "<VERSION>",
    dev_dependency = True,
)
```

Add the public language target to the repository's custom Gazelle binary:

```starlark
load("@gazelle//:def.bzl", "gazelle_binary")

gazelle_binary(
    name = "gazelle_binary",
    languages = [
        "@rules_docs_gazelle//gazelle",
    ],
)
```

The extension only operates in directories where both a BUILD file and a
`README.md` already exist. Newly generated rules use `glob(["*.md"])` and are
visible only to their nearest ancestor Bazel package. Existing `srcs`, `deps`,
`visibility`, and `prefix` attributes are preserved.

The public `@rules_docs_gazelle//gazelle:gazelle_docs` binary contains only
this extension and can update a nested workspace without initializing
unrelated language plugins from its parent repository.

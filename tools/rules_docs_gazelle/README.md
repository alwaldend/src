---
title: Rules docs Gazelle
description: Gazelle extension for Bazel documentation packaging rules
statuses:
  - active
languages:
  - go
tags:
  - bzl_rules
---

`rules_docs_gazelle` adds a `docs_filegroup` to existing Bazel packages that
contain a `README.md`. The generated rule loads its macro from `rules_docs`.

Add `rules_docs` as a normal dependency and this generator as a development
dependency:

```starlark
bazel_dep(name = "rules_docs", version = "<VERSION>")
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

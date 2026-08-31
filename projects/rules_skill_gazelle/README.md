---
title: Rules skill Gazelle
description: Gazelle extension for Bazel skill libraries
statuses:
  - active
languages:
  - go
tags:
  - bzl_rules
---

`rules_skill_gazelle` adds a `skill_library(name = "skill")` to each named
subpackage that contains a `SKILL.md`. The generated rule loads its macro from
the apparent `rules_skill` repository, so Bzlmod repository mappings are
respected.

Add `rules_skill` as a normal dependency and this generator as a development
dependency:

```starlark
bazel_dep(name = "rules_skill", version = "<VERSION>")
bazel_dep(
    name = "rules_skill_gazelle",
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
        "@rules_skill_gazelle//gazelle",
    ],
)
```

The generated source bundle is:

```starlark
load("@rules_skill//skill:defs.bzl", "skill_library")

skill_library(
    name = "skill",
    srcs = glob(
        ["**"],
        exclude = [
            "BUILD.bazel",
            "BUILD",
            "evals/**",
        ],
    ),
)
```

The generator always excludes `BUILD.bazel` and `BUILD`, plus every custom
BUILD filename configured in Gazelle. Duplicate configured names are removed.

`SKILL.md` is sufficient for Gazelle to create a new BUILD file below the
repository root. A root-level `SKILL.md` is ignored because `rules_skill`
requires a named package from which it can derive the skill name. Existing
attributes on a manually maintained `skill_library(name = "skill")` are
preserved, and a missing `SKILL.md` does not delete a manual rule.

The public `@rules_skill_gazelle//gazelle:gazelle_skills` binary contains only
this extension and can update a nested workspace without initializing
unrelated language plugins from its parent repository.

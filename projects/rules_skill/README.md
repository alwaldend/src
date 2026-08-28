---
title: Rules skill
description: Bazel rules and validation for Codex skills
languages:
  - bzl
  - go
tags:
  - bzl_rules
---

`rules_skill` packages all files belonging to a Codex skill and validates its
instructions and optional OpenAI metadata with a hermetic Bazel aspect.

## Getting started

Add the module and register the validation aspect:

```starlark
bazel_dep(name = "rules_skill", version = "<VERSION>")
```

```text
build --aspects @rules_skill//skill:defs.bzl%skill_validation_aspect
build --output_groups=+skill_validation
```

Declare one library in the skill's package:

```starlark
load("@rules_skill//skill:defs.bzl", "skill_library")

skill_library(
    name = "skill",
    srcs = glob(["**"], exclude = ["BUILD.bazel"]),
)
```

Building the library materializes the `skill_validation` output group. The
aspect checks `SKILL.md`, verifies that the frontmatter name matches its
package directory, and validates `agents/openai.yaml` when it is present.

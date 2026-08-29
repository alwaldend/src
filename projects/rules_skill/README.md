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

Declare one library in the skill's named, non-root package:

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

The package's final path segment is the skill name, so a `skill_library`
cannot be declared in a repository root package. Move a root-level skill into
a named subpackage before declaring it.

Building the library materializes the `skill_validation` output group. The
aspect checks `SKILL.md`, verifies that the frontmatter name matches its
package directory, and validates `agents/openai.yaml` when it is present.

## Consuming a skill

Rules that install or evaluate skills can require the public `SkillInfo`
provider:

```starlark
load("@rules_skill//skill:defs.bzl", "SkillInfo")

attrs = {
  "skill": attr.label(providers = [SkillInfo]),
}
```

`SkillInfo` exposes these fields:

- `name` is the logical skill name derived from the final segment of `root`.
  The validation aspect verifies that the `SKILL.md` frontmatter uses the same
  name.
- `root` is the owning Bazel package path within the skill's repository, such
  as `projects/agents/skills/answer-question`. It has no repository,
  execution-path, or runfiles prefix. It is always non-empty because
  repository root packages are unsupported.
- `files_by_path` maps slash-separated paths relative to `root` to Bazel
  `File` values. It includes `SKILL.md`, preserves nested paths such as
  `agents/openai.yaml`, and has the same shape for source and generated files.
- `files` remains the `depset` of all skill files. `skill` is the distinguished
  `SKILL.md` file, and `openai_yaml` is the optional distinguished
  `agents/openai.yaml` file.

Consumers should use `files_by_path` when staging a bundle instead of parsing
`File.path` or `File.short_path`. Every source must belong to the skill's Bazel
package, and duplicate logical paths are rejected during analysis.

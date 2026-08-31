---
title: Rules skill
description: Bazel rules and validation for Codex skills
statuses:
  - active
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

## Maintaining source-tree discovery links

Use `skill_discovery_links` when a repository exposes canonical skill
directories through a local discovery index such as `.agents/skills`. Pass
the `skill_library` labels directly; the BUILD declaration is the complete
source of truth, so there is no separate link manifest or configuration file.

```starlark
load(
    "@rules_skill//skill:defs.bzl",
    "skill_discovery_links",
)

skill_discovery_links(
    name = "write_skill_links",
    skills = [
        "//projects/agents/skills/answer-question:skill",
        "//projects/example/skills/example:skill",
    ],
)
```

Reconcile the checkout with the runnable target:

```sh
bazel run //:write_skill_links
```

The macro also declares `//:write_skill_links_test`. That local, uncached test
reads the source checkout and requires the discovery directory to contain
exactly the skill names passed in BUILD, each as a direct relative symlink to
its canonical `SkillInfo.root`. It rejects missing, stale, extra, absolute, or
indirect links. Run the updater before the test after changing the label list:

```sh
bazel test //:write_skill_links_test
```

Both targets intentionally accept only source `SKILL.md` files from the same
repository. The exact-state test is a pinned-local checkout check rather than
a hermetic or remotely executable test: it resolves a source runfile back to
the workspace and then inspects the local symlinks. The generated programs
require Bash and POSIX symlink support and work on ordinary Linux and macOS
checkouts; native Windows checkouts are not supported.

The updater holds a sibling `<discovery_dir>.lock` directory while it changes
links, so concurrent declarations cannot produce a mixed set. A process killed
without running its exit trap can leave that lock behind; after confirming no
updater is active, remove the empty lock directory and rerun the updater. The
exact-state test detects any interrupted partial reconciliation.

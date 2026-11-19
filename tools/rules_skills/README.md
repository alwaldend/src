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

`rules_skills` packages all files belonging to a Codex skill and validates its
instructions and optional OpenAI metadata with a hermetic Bazel aspect.

## Getting started

Add the module and register the validation aspect:

```starlark
bazel_dep(name = "rules_skills", version = "<VERSION>")
```

```text
build --aspects @rules_skills//skill:defs.bzl%skill_validation_aspect
build --output_groups=+skill_validation
```

Declare one library in the skill's named, non-root package:

```starlark
load("@rules_skills//skill:defs.bzl", "skill_library")

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
load("@rules_skills//skill:defs.bzl", "SkillInfo")

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

## Packaging skills from an external archive

Use `skill_archives` when skills are maintained upstream and consumed as a
pinned archive. The external repository's `BUILD` file derives one target per
skill directory from a wildcard, so a new upstream skill needs no edit here:

```starlark
load("@rules_skills//skill:defs.bzl", "skill_archives")

skill_archives(
    name = "skills",
    roots = [
        manifest.rsplit("/", 1)[0]
        for manifest in glob(["skills/*/SKILL.md"])
    ],
)
```

`skill_archives` is a macro over the `skill_archive` rule. Each generated
target's name is its skill directory's final segment, `root` is the skill
directory inside the package, and its packaged logical paths are relative to
`root`. Analysis fails when a root has no `SKILL.md`. Declare one
`skill_archive` directly when a single skill needs a hand-written target.

## Writing the source-tree discovery directory

Use `skills_write` to reconcile `.agents/skills` from declared skills. Each
`symlinks` label is installed as a direct relative symlink to its canonical
root; each `archives` label is copied in as regular files, which is how archive
skills are materialized into the consuming repository:

```starlark
load("@rules_skills//skill:defs.bzl", "skills_write")

skills_write(
    name = "write_skills",
    archives = ["@org_fissionai_openspec//:openspec-propose"],
    discovery_dir = ".agents/skills",
    symlinks = ["//projects/agents/skills/answer-question:skill"],
    workspace_marker = "//:AGENTS.md",
)
```

`workspace_marker` is a source file in the consuming repository used to
resolve the workspace root from test runfiles. Reconcile and verify:

```sh
bazel run //.agents:write_skills
bazel test //.agents:write_skills_test
```

The generated check requires exactly the declared names: missing, extra, or
stale entries fail, and a written entry must match its declared source
byte-for-byte. The updater holds a sibling `<discovery_dir>.lock` directory
while it changes entries; a process killed without running its exit trap can
leave that lock behind. After confirming no updater is active, remove the empty
lock directory and rerun the updater.

Both rules require POSIX symlinks and Bash and are limited to Linux and macOS
checkouts; native Windows checkouts are not supported.

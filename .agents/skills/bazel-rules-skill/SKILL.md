---
name: bazel-rules-skill
description: >-
  Add or update repository skills packaged by the rules_skill Bazel rule and
  checked by its validation aspect. Use for .agents/skills content and
  BUILD.bazel declarations in this monorepo.
---

# Add a repository skill

## Create the skill

1. Read the root `AGENTS.md`, use the `skill-creator` skill, and inspect the
   closest existing repository skill.
2. Put the skill in `.agents/skills/lowercase-hyphen-name`. Keep the directory
   name and the `name` in `SKILL.md` identical.
3. Add only resources the skill needs. Put product metadata in
   `agents/openai.yaml`; its `default_prompt` must mention
   `$lowercase-hyphen-name`.
4. Do not create nested Bazel packages inside the skill. A recursive Bazel glob
   stops at package boundaries and would omit those files from the library.

## Package every skill file

Add `.agents/skills/lowercase-hyphen-name/BUILD.bazel`:

```starlark
load("@rules_skill//skill:defs.bzl", "skill_library")

skill_library(
    name = "skill",
    srcs = glob(["**"], exclude = ["BUILD.bazel"]),
)
```

The library exposes the skill files through `DefaultInfo` and `SkillInfo`.
Excluding `BUILD.bazel` keeps repository build metadata out of the runtime
skill bundle. Ensure every referenced script, icon, asset, and metadata file
remains inside the glob.

## Validate through the aspect

The root Bazel configuration registers `skill_validation_aspect` and requests
its `skill_validation` output group. Build the library to execute validation:

```sh
bazel build --config=agent //.agents/skills/lowercase-hyphen-name:skill
bazel test --config=agent //:buildifier_test
```

The aspect validates the `SKILL.md` frontmatter and body, checks that the name
matches the package directory, and validates optional `agents/openai.yaml`
metadata. Fix the underlying skill instead of suppressing the aspect.

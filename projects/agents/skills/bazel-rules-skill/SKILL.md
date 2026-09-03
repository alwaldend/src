---
name: bazel-rules-skill
description: >-
  Add or update repository skills packaged by the rules_skill Bazel rule and
  checked by its validation aspect. Use for canonical project-owned skill
  content, .agents discovery links, and skill BUILD.bazel declarations in this
  monorepo.
---

# Add a repository skill

## Create the skill

1. Read the root `AGENTS.md`, use the `skill-creator` skill, and inspect the
   closest existing repository skill.
2. Identify the narrowest project that owns the skill. Put a product-specific
   skill at `<owner-project>/skills/lowercase-hyphen-name`. Use
   `projects/agents/skills/lowercase-hyphen-name` only for a repository-wide
   agent workflow with no narrower owner. Keep the directory and frontmatter
   names identical.
3. Add only resources the skill needs. Put product metadata in
   `agents/openai.yaml`; its `default_prompt` must mention
   `$lowercase-hyphen-name`.
4. Do not create nested Bazel packages inside the skill. A recursive Bazel glob
   stops at package boundaries and would omit those files from the library.

## Package every skill file

Add `<owner-project>/skills/lowercase-hyphen-name/BUILD.bazel`:

```starlark
load("@rules_skill//skill:defs.bzl", "skill_library")

skill_library(
    name = "skill",
    srcs = glob(
        ["**"],
        exclude = [
            "BUILD",
            "BUILD.bazel",
            "evals/**",
        ],
    ),
    visibility = ["//:__pkg__"],
)
```

The library exposes the skill files through `DefaultInfo` and `SkillInfo`.
The root-package visibility lets the repository discovery-link target consume
the provider while keeping the canonical skill unavailable to other packages.
Excluding both conventional BUILD filenames keeps repository build metadata
out of the runtime skill bundle. Keep Promptfoo configurations, cases, and
eval documentation under `evals/` and exclude them as well; they test the skill
but are not part of the instructions delivered to an agent. Ensure every
referenced script, icon, asset, and metadata file remains inside the glob.

## Register skill discovery

Add the canonical `:skill` label to the complete `skills` list in the root
`skill_discovery_links` declaration:

```starlark
load(
    "@rules_skill//skill:defs.bzl",
    "skill_discovery_links",
)

skill_discovery_links(
    name = "write_skill_links",
    skills = [
        # Existing canonical skill targets...
        "//projects/storage/skills/database-backup:skill",
    ],
)
```

The Bazel declaration is the only discovery configuration. The updater derives
each link name and relative target from `SkillInfo`; do not add a separate
manifest or create `.agents/skills` links by hand. Run the updater, then its
generated exact-state test:

```sh
bazel_agent bazel run //:write_skill_links
bazel_agent bazel test //:write_skill_links_test
```

The test rejects missing, extra, indirect, or incorrectly targeted local
links. Package, validate, and edit only the canonical skill directory; never
edit through a discovery link.

## Add Promptfoo coverage

Every new or updated repository skill must include at least one checked-in
Promptfoo configuration and an offline `promptfoo_validate_test` for every
configuration. Keep the configuration and cases under `evals/`, declare all
referenced files in `data`, and stage the skill through `skills = [":skill"]`
when the config tests skill discovery. Offline validation proves that Bazel can
stage the skill and that Promptfoo can load the complete configuration without
credentials or model calls; do not describe it as evidence that the skill's
answers are correct.

```starlark
load(
    "@rules_promptfoo//promptfoo:defs.bzl",
    "promptfoo_test",
    "promptfoo_validate_test",
)

filegroup(
    name = "eval_cases",
    testonly = True,
    srcs = ["evals/cases.yaml"],
)

promptfoo_validate_test(
    name = "eval_config_test",
    config = "evals/promptfooconfig.yaml",
    data = [":eval_cases"],
    skills = [":skill"],
)
```

Add a manual `promptfoo_test` for behavioral coverage when the suite can run
safely and meaningfully with the available provider. An online test is
optional when representative behavior requires tool calls or external state
that the eval cannot safely and reproducibly provide. When omitting it,
document the missing tool surface and the resulting coverage gap in
`evals/README.md`; the offline validation target remains required. Never make a
credentialed or billable live eval part of ordinary `//...` execution.

For live coverage, include a no-skill or previous-version control when it can
support the claimed improvement, keep quality and skill-routing assertions
separate, and treat model-judge results as regression evidence rather than
proof. Use the `rules_promptfoo` README for credential isolation, local-only
execution, result privacy, and target-tag requirements.

## Validate through the aspect

The root Bazel configuration registers `skill_validation_aspect` and requests
its `skill_validation` output group. Build the library to execute validation:

```sh
bazel_agent bazel build //<owner-project>/skills/lowercase-hyphen-name:skill
bazel_agent bazel test \
  //<owner-project>/skills/lowercase-hyphen-name:eval_config_test
bazel_agent bazel test //:write_skill_links_test
bazel_agent bazel test //:buildifier_test
```

The aspect validates the `SKILL.md` frontmatter and body, checks that the name
matches the package directory, and validates optional `agents/openai.yaml`
metadata. Fix the underlying skill instead of suppressing the aspect.

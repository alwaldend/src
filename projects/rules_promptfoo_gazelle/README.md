---
title: Rules Promptfoo Gazelle
description: Gazelle extension for offline Promptfoo validation tests
languages:
  - go
tags:
  - bzl_rules
  - llm_evals
---

`rules_promptfoo_gazelle` generates required, offline
`promptfoo_validate_test` targets for Promptfoo configurations below an
`evals/` directory. It never creates an `evals` subpackage and never generates
live `promptfoo_test` targets.

Add `rules_promptfoo` as a normal dependency and this generator as a
development dependency:

```starlark
bazel_dep(name = "rules_promptfoo", version = "<VERSION>")
bazel_dep(
    name = "rules_promptfoo_gazelle",
    version = "<VERSION>",
    dev_dependency = True,
)
```

For one-pass skill package creation, also add `rules_skill` as a normal
dependency and `rules_skill_gazelle` as a development dependency:

```starlark
bazel_dep(name = "rules_skill", version = "<VERSION>")
bazel_dep(
    name = "rules_skill_gazelle",
    version = "<VERSION>",
    dev_dependency = True,
)
```

Add both public language targets to the repository's custom Gazelle binary:

```starlark
load("@gazelle//:def.bzl", "gazelle_binary")

gazelle_binary(
    name = "gazelle_binary",
    languages = [
        "@rules_promptfoo_gazelle//gazelle",
        "@rules_skill_gazelle//gazelle",
    ],
)
```

With both extensions installed, a fresh directory containing `SKILL.md` and a
conventional Promptfoo configuration receives both `skill_library` and
`promptfoo_validate_test` targets in one Gazelle run; it does not need an
existing BUILD file.

The extension recognizes these conventional configuration paths:

```text
evals/promptfooconfig.yaml
evals/promptfooconfig.<variant>.yaml
evals/promptfooconfig.yml
evals/promptfooconfig.<variant>.json
```

The default file produces `eval_config_test`; a variant produces
`eval_<sanitized_variant>_config_test`. Variants are lowercased, and runs of
characters other than ASCII letters and digits become underscores. Colliding
sanitized names receive a stable path-derived `__<8 lowercase hex digits>`
suffix. The double underscore cannot be produced by the sanitizer, so a
variant cannot claim a collision target's namespace.

When a non-root containing package has a `SKILL.md`, each ordinary validation
target stages that package's conventional `:skill` target. Repository-root
skills are unsupported by `rules_skill`, so a root-package `SKILL.md` does not
infer that label. The exact `no_skill` variant omits `skills` so it remains a
control. A generic Promptfoo package without `SKILL.md` receives no inferred
skill label. Every non-configuration regular file under `evals/`, including
`README.md` and test cases, is added to `data` in sorted order. Directories
containing `BUILD`, `BUILD.bazel`, or an additional filename configured as
valid in Gazelle are package boundaries and are not traversed.

The extension reconciles `config`, `data`, and `skills` when a discovered
configuration matches a conventionally named validation target. Adding or
editing eval files therefore updates those attributes on the next run.
Removing or renaming a configuration does not delete the old target: a name
shape is not reliable proof that Gazelle created a rule, so stale targets must
be removed explicitly. Gazelle's `# keep` semantics remain available at the
rule, attribute, and individual list-item levels for intentional overrides.
Other attributes, including `args`, `env`, `tags`, and `visibility`, remain
manual and are preserved while the generated target exists. Existing manual
or live `promptfoo_test` rules are also left intact.
The extension does not follow a symlink used as the package's `evals`
directory, because it could escape the repository. Such a symlink, or another
unexpected filesystem scan error including an error while checking a candidate
package boundary, is a non-destructive no-op rather than a signal to remove
targets. Before scanning, it also resolves the repository root and package
directory and refuses a package whose resolved path escapes the workspace,
including a package reached through Gazelle's symlink-following support. The
extension manages the shared load for both Promptfoo rule symbols and uses
Bzlmod's apparent repository name.

The public `@rules_promptfoo_gazelle//gazelle:gazelle_promptfoo` binary
contains only this extension and can update a nested workspace without
initializing unrelated language plugins from its parent repository.

---
name: repo-bazel
description: Build, test, query, and maintain targets in this Bazel monorepo. Use for BUILD.bazel, .bzl, MODULE.bazel mechanics, Gazelle, buildifier, toolchains, target selection, or repository-wide validation work; use repo-external-dependency to add or upgrade external software.
---

# Work with Bazel

## Inspect before editing

1. Follow the `bazel-agent` skill for every repository Bazel invocation. Use
   `bazel_agent <command> ...`; do not repeat its flags or call `bazel`
   directly.
2. Read the nearest `AGENTS.md` and the applicable `BUILD.bazel` files.
3. Inspect `.bazelrc` and its imports before changing flags. Preserve the final
   `try-import %workspace%/user.bazelrc` so local overrides remain last.
4. Find existing rules with `rg`; copy the closest repository pattern instead of
   inventing a parallel macro.
5. Use `bazel_agent query` to confirm labels and dependency boundaries when
   uncertain.

## Edit targets

- Keep Starlark formatted with buildifier conventions.
- Prefer repository macros under `tools/` and `projects/al/rules/` over raw
  upstream rules. Infrastructure commonly uses generated target maps such as
  `terraform_binary_map`, `terraform_test_map`, and `vault_binary_map`.
- Keep runtime files in `data`; declaring a file in `srcs` does not necessarily
  make it available to an executed tool.
- Preserve least-privilege `visibility` and use package-relative labels when
  that is the surrounding convention. When the same consumer target or package
  must be granted access from several packages, define a named `package_group`
  and grant visibility to that group instead of repeating the same package
  label in every dependent package. Name the group for the capability it
  authorizes, not for a list of unrelated packages.
- When a checked-in generated file must not be rewritten by repo-wide
  formatting (for example, generator-owned JSON with a fixed layout), exclude
  it in `.gitattributes` with `rules-lint-ignored=true`. Verify the exclusion
  with the relevant `//tools/repo_quality` format test before delivery.
- Update dependency lockfiles only through the owning repository workflow. Do
  not mix unrelated generated-file churn into a change.

## Validate narrowly, then broadly

Run the cheapest relevant commands first:

```sh
bazel_agent query //path/to/package:all
bazel_agent test //path/to/package:all
bazel_agent build //path/to/package:all
bazel_agent test //:buildifier_test
```

Use `bazel_agent test //...` or `bazel_agent build //...` only when the scope
and available time justify a repository-wide check. Report remote-cache,
credential, network, or platform failures separately from failures caused by
the patch.

For generated targets, inspect available labels before guessing:

```sh
bazel_agent query '//path/to/package:*'
```

Never run deploy, apply, or other mutating `bazel_agent run` targets merely as a
validation step.

---
name: repo-bazel
description: Build, test, query, and maintain targets in this Bazel monorepo. Use for BUILD.bazel, .bzl, MODULE.bazel, dependency, Gazelle, buildifier, toolchain, target-selection, or repository-wide validation work.
---

# Work with Bazel

## Inspect before editing

1. Read the nearest `AGENTS.md` and the applicable `BUILD.bazel` files.
2. Inspect `.bazelrc` and its imports before changing flags. Preserve the final
   `try-import %workspace%/user.bazelrc` so local overrides remain last.
3. Find existing rules with `rg`; copy the closest repository pattern instead of
   inventing a parallel macro.
4. Use `bazel query` to confirm labels and dependency boundaries when uncertain.

## Edit targets

- Keep Starlark formatted with buildifier conventions.
- Prefer repository macros under `tools/` and `projects/al/rules/` over raw
  upstream rules. Infrastructure commonly uses generated target maps such as
  `terraform_binary_map`, `terraform_test_map`, and `vault_binary_map`.
- Keep runtime files in `data`; declaring a file in `srcs` does not necessarily
  make it available to an executed tool.
- Preserve least-privilege `visibility` and use package-relative labels when
  that is the surrounding convention.
- Update dependency lockfiles only through the owning repository workflow. Do
  not mix unrelated generated-file churn into a change.

## Validate narrowly, then broadly

Run the cheapest relevant commands first:

```sh
bazel query --config=agent //path/to/package:all
bazel test --config=agent //path/to/package:all
bazel build --config=agent //path/to/package:all
bazel test --config=agent //:buildifier_test
```

Use `bazel test --config=agent //...` or
`bazel build --config=agent //...` only when the scope and available time
justify a repository-wide check. Report remote-cache, credential, network, or
platform failures separately from failures caused by the patch.

For generated targets, inspect available labels before guessing:

```sh
bazel query --config=agent '//path/to/package:*'
```

Never run deploy, apply, or other mutating `bazel run` targets merely as a
validation step.

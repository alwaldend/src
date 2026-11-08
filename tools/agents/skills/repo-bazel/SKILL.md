---
name: repo-bazel
description: Build, test, query, and maintain targets in this Bazel monorepo. Use for BUILD.bazel, .bzl, MODULE.bazel mechanics, Gazelle, buildifier, toolchains, target selection, or repository-wide validation work; use repo-external-dependency to add or upgrade external software.
---

# Work with Bazel

## Inspect before editing

1. Follow the `al-bazel` skill for invocation and setup. Run
   `bazel <command> ...` directly; account-specific defaults come from the
   user's Bazel configuration.
2. Read the nearest `AGENTS.md` and the applicable `BUILD.bazel` files.
3. Inspect `.bazelrc` and its imports before changing flags. Preserve the final
   `try-import %workspace%/user.bazelrc` so local overrides remain last.
4. Find existing rules with `rg`; copy the closest repository pattern instead of
   inventing a parallel macro.
5. Use `bazel query` to confirm labels and dependency boundaries when
   uncertain.

## Edit targets

- Keep Starlark formatted with buildifier conventions.
- Keep short, single-use configuration in BUILD.bazel, including small
  dictionaries and direct calls to existing macros. Extract a local defs.bzl
  only when it provides significant reuse or makes complex construction logic,
  resource graphs, long templates or bulky settings substantially easier to
  read. Do not add a single-use wrapper merely to bind a few fixed attributes
  to an existing rule or macro. Put genuinely common patterns in reusable
  namespaced target macros with their owner. Compose small targets instead of
  hiding a whole large config dictionary behind a single macro.
- For AL configuration, combine small targets through configs. Declare intended
  overrides with per-target extends relationships, including component auth
  extending //tools/al:auth. Aggregate list order does not select overrides;
  let the formatter sort references. Ordered resource dependency lists retain
  their runtime meaning. Use plugin_labels dictionaries for call selection;
  wrappers own CLI encoding and run_args is reserved for other AL options.
  Prefer al_injector_call with named `res_<type>` maps;
  keep structured settings and explicit deps/Vault metadata in each entry.
  Generic calls default call and plugin names to the target name; override
  config_name/plugin only when needed. Injector calls select injector directly.
- Namespace public Starlark rules, macros, and configuration constructors with
  their owning component prefix: `infra_dns`, `al_terraform_binary_map`, and
  `vault_tf_backend`, for example. Providers use the corresponding CamelCase
  namespace, such as `AlConfigInfo`. Apply this convention when creating or
  renaming an API and update its callers.
- Prefer repository macros under `tools/` and `@rules_al//:defs.bzl` over raw
  upstream rules. Infrastructure commonly uses generated target maps such as
  `al_terraform_binary_map`, `al_terraform_test_map`, and `vault_binary_map`.
- For Terraform, use the public rules in
  [`rules_terraform`](../../../../tools/rules_terraform/README.md). Repository
  consumers use `//tools/al:al_terraform.bzl` for AL integration, passing
  `configs` and `plugin_labels` directly. Plain offline checks can use the standalone
  Terraform rules. Reuse `vault_tf_backend` from
  `//tools/vault/tf_backend:defs.bzl` for backend call targets.
  Provider pins belong to [third_party/terraform](../../../../third_party/terraform/README.md), with selected
  provider labels in rule attributes; do not add source Terraform lockfiles.
- Keep runtime files in `data`; declaring a file in `srcs` does not necessarily
  make it available to an executed tool.
- Do not use `genrule`. Write a proper Bazel rule, a Go binary, or a separate
  checked-in script instead. When no suitable rule exists, consider adding one
  to the owning tools package; do not encode build logic in shell snippets.
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

## Sandbox and network policy

Run every build action inside a Bazel sandbox and without network access.
If an action genuinely needs network access or must disable sandboxing,
document the reason in the owning project's README before enabling it and
prefer the narrowest possible exemption.

## Cached outputs and no-op edits

A successful no-op patch is settled evidence that something other than the
target bytes is in control. Do not repeat it. Diagnose the action cache,
generated-file ownership, macro name expansion, or output path, then make a
single state-changing correction. When a build reports up-to-date after a
behavioral source change, inspect the exact generated inputs and action
arguments before rebuilding; if a cache is stale, invalidate the target
once with a real input change rather than alternating no-op edits with
builds.

## Validate narrowly, then broadly

Run the cheapest relevant commands first:

```sh
bazel query //path/to/package:all
bazel test //path/to/package:all
bazel build //path/to/package:all
bazel test //:buildifier_test
```

Use `bazel test //...` or `bazel build //...` only when the scope
and available time justify a repository-wide check. Report remote-cache,
credential, network, or platform failures separately from failures caused by
the patch.

For generated targets, inspect available labels before guessing:

```sh
bazel query '//path/to/package:*'
```

Never run deploy, apply, or other mutating `bazel run` targets merely as a
validation step.

## Automate with Go and Bazel

Prefer Go and Bazel-native rules for repository automation. Do not introduce
shell scripts or shell-based rules unless the task explicitly requires them.
Expose Go automation as `go_binary` targets, and validate it with the pinned
Bazel toolchain rather than host `go`. Write a proper rule, a checked-in Go
binary, or a separate checked-in script instead of encoding build logic in a
shell snippet.

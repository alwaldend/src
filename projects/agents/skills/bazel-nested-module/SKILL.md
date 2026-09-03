---
name: bazel-nested-module
description: >-
  Create or update a standalone nested Bzlmod module under projects/ and wire
  it into this monorepo. Use for projects with their own MODULE.bazel; do not
  use for ordinary Bazel packages or Git submodules.
---

# Create a nested Bazel module

Follow the root `AGENTS.md` and the `repo-bazel` skill. Inspect all existing
`projects/rules_*` modules before changing conventions shared by them.

## Create the standalone workspace

1. Use `projects/<module_name>` with an underscore-separated module name that
   matches `module(name = ...)`. Add `MODULE.bazel`, `MODULE.bazel.lock`,
   `README.md`, and a root `BUILD.bazel`. New modules use `version = "0.0.0"`
   and `bazel_compatibility = [">=8.0.0"]` until release policy requires a
   different value.
2. Symlink `.bazeliskrc` to `../../.bazeliskrc` and `.bazelignore` to
   `../../.bazelignore`. Normally symlink `.bazelrc` to
   `../../tools/bazelrc/bzl_project.bazelrc` too. If the module needs
   workspace-wide flags such as a validation aspect, use a regular `.bazelrc`
   that imports `%workspace%/../../tools/bazelrc/preset.bazelrc`, then
   `%workspace%/../../tools/bazelrc/project.bazelrc`, adds only its custom
   flags, and `try-import`s `%workspace%/../../user.bazelrc` last. Do not copy
   shared settings into it.
3. Declare only dependencies the module needs. For repository documentation,
   depend on `rules_docs` and load `@rules_docs//docs:defs.bzl`; `rules_docs`
   itself uses the canonical `//docs:defs.bzl` label instead of depending on
   itself. This dependency must remain regular because published BUILD files
   load it. Declare `rules_docs_gazelle` as a development dependency when the
   workspace runs the documentation Gazelle extension. Never copy a fallback
   documentation macro into the nested workspace. A sibling
   `local_path_override` may make an unpublished module usable in this
   checkout, but remove or replace that override when preparing a registry
   release.
4. Give the README the repository's documentation frontmatter. Expose a
   public root `docs_filegroup` with the explicit
   `content/docs/projects/<module_name>` prefix so the parent documentation
   aggregate preserves its destination layout.
5. Add focused `build_test` or analysis tests for the module's public rules.
   Add Go, Gazelle, toolchain, or extension setup only when the implementation
   requires it; follow the closest nested module rather than installing a
   generic template.
6. Keep the module self-contained: never reference the parent repository from
   inside a nested module. Do not use root-workspace labels such as
   `//tools/...` in BUILD files or rule defaults, because they resolve to
   `@<module>//tools/...` under Bzlmod and make reusable modules dependent on
   the parent repository. Use module-internal labels, add a mandatory attr for
   a tool that must be supplied by the consumer, or declare a dependency the
   module owns.

## Integrate it with the parent repository

- Add the nested directory to the root `.bazelignore` so root package
  discovery does not cross workspace boundaries.
- Add the module's `bazel_dep` and `local_path_override` beside the other local
  rules modules in `third_party/include.MODULE.bazel`. If parent targets use a
  toolchain supplied by the module, register that toolchain there as well;
  otherwise do not add a repository-wide registration.
- Add its external `:docs` target to `//projects:docs`. Exclude the directory
  from local `subpackages()` expansion and from `//projects:deploy_heads`
  unless the module intentionally provides the required release target.
- Add the workspace to the `full-repo-check` runner and update its expected
  command count. A complete repository check must build and test every nested
  module independently.
- Register any Gazelle language library in the root `gazelle_binary`. Remember
  that root Gazelle ignores nested workspaces; run the plugin against the
  nested repository root as well when it should generate that module's BUILD
  files. Keep development-only Gazelle tooling such as `rules_docs_gazelle`
  out of downstream module graphs with `dev_dependency = True`; do not mark a
  module as development-only when published BUILD files load from it.

## Generate and verify

Update, rather than hand-edit, module locks after dependency changes:

```sh
bazel_agent mod deps --lockfile_mode=update
```

Run that command from every affected workspace root. Then run Gazelle where
needed, review every generated change, and run it a second time to prove the
result is stable. Validate at least:

```sh
bazel_agent bazel test //...
bazel_agent bazel build //...
```

Run those commands inside the new module and run focused parent-repository
checks for its integration, including `bazel_agent bazel test //:buildifier_test`.
Finish with `full-repo-check` when adding a workspace so the updated
nested-workspace matrix is exercised.

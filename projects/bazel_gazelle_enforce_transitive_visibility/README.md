---
title: Bazel Gazelle transitive visibility
description: Generate package dependency restrictions from TOML path mappings
languages:
  - go
---

This standalone Gazelle extension maintains `package(transitive_visibility =
...)` from a TOML mapping of workspace-relative directories to package groups.
The default configuration is `tools/transitive_visibility/transitive_visibility.toml`.

```toml
[paths]
"." = "//:transitive_consumers"
"infra" = "//infra:transitive_consumers"
"tools" = "//tools:transitive_consumers"
"tools/nogo" = "//tools:build_api_consumers"
```

Mappings include descendant directories. The longest directory-prefix match
wins; `infra` does not match `infrastructure`. Paths are canonical relative
directories, with `"."` denoting the workspace root; globs and parent paths are
not accepted. Labels must be absolute. Package groups remain in their owning
BUILD files and define the permitted consumers.

Use `-transitive_visibility_config=path/to/policy.toml` to select another file.
Relative filenames resolve from Gazelle's repository root. A missing default
configuration leaves policy unmanaged; an explicitly selected missing file is
an error. Malformed TOML, unknown fields, empty mappings, invalid paths and
invalid labels fail before BUILD output is written. Obsolete
`gazelle:transitive_visibility` comments are rejected with a migration message.

Only the `transitive_visibility` attribute is managed. Other attributes and
comments are preserved, and managed `package()` calls in existing files appear after loads and
before rules. Unmatched packages retain their manual declarations. The plugin
does not create BUILD files in otherwise empty directories. Register it after
other generators so it covers the packages they create. Gazelle's normal
ignored-directory and nested-workspace boundaries apply; parent configurations
are never searched.

Register `@bazel_gazelle_enforce_transitive_visibility//pkg/gazelle` in the
`languages` of a `gazelle_binary`. `//cmd/gazelle` provides the plugin-only
executable. This module's own configuration is `transitive_visibility.toml`:

```sh
bazel_agent bazel run //cmd/gazelle -- -transitive_visibility_config=transitive_visibility.toml
```

Bazel enforcement is separate from generation. Enable
`--experimental_enforce_transitive_visibility=true` and use Bazel 9.3.1 or
later. Bazel checks target dependencies for the entire package; it cannot
distinguish target roles inside that package and does not govern Starlark
`load()` statements.

## Enforcement limitation in Bazel 9.3.1

Direct target and source-file restrictions work. An integration experiment
also found that an intermediary with its own permissive `transitive_visibility`
declaration can lose an inherited restriction: a forbidden consumer of its
`filegroup` succeeds. Removing the intermediary's declaration makes the same
indirect dependency fail. Consequently, generated declarations are not a
guarantee of complete indirect dependency enforcement in this Bazel release.
The extension emits the requested policy; it does not repair Bazel's analysis.

The native check is also broader than an artifact-only dependency policy.
Its [9.3.1 validator](https://github.com/bazelbuild/bazel/blob/9.3.1/src/main/java/com/google/devtools/build/lib/analysis/CommonPrerequisiteValidator.java#L579)
does not exempt execution-tool attributes; declaring a dependency with
`cfg = "exec"` does not itself exempt it. Package groups cannot express
"allow this executable as a build tool, but forbid it as a runtime input."
Explicit build-support package exceptions are therefore an approximation.
Starlark `load()` edges are separate. Use `--notool_deps` when querying the
ordinary dependency graph; that query option does not change build enforcement.

From this module, run `bazel_agent bazel test //...` to exercise the real
Gazelle command. The tests preserve generated BUILD files and a summary under
Bazel's undeclared test outputs.

Gazelle 0.54.0 appends `package()` after generated rules when it creates a new
BUILD file. A subsequent run moves it after the loads; further runs are stable.
The repository uses unmodified upstream Gazelle and does not patch this behavior.

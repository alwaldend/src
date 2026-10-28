---
title: Transitive visibility policy
---

`transitive_visibility.toml` selects the package groups used by the repository's
Gazelle transitive-visibility extension. The groups themselves remain in their
owning trees' BUILD files.

The `[paths]` table maps workspace-relative directories to absolute group labels.
Use `"."` for the workspace root. A mapping includes descendants; the most
specific directory match wins. Paths are literal directories, not globs.
The `tools:build_api_consumers` entries preserve existing build-tool and
toolchain exceptions. Website, developer, and test exceptions remain expressed
by the package groups.

Run `bazel_agent bazel run //:gazelle -- -lang=enforce_transitive_visibility`
after changing this file. Gazelle validates all mappings before updating BUILD
files and places managed `package()` declarations in existing files after loads
and before rules.
Do not add `gazelle:transitive_visibility` comments; the extension rejects them.

See the [plugin documentation](../../projects/bazel_gazelle_enforce_transitive_visibility/README.md)
for standalone usage and the Bazel 9.3.1 indirect-enforcement limitation.

Gazelle 0.54.0 appends `package()` after generated rules when it creates a new
BUILD file. A subsequent run moves it after the loads; further runs are stable.
The repository uses unmodified upstream Gazelle and does not patch this behavior.

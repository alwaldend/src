## Decisions

- Use the pinned Go SDK's `golang.org/x/mod/modfile` package to parse and
  rewrite the `go` directive. It preserves file ordering and formatting and
  avoids invoking `go mod edit`, whose workspace hooks fail outside the root
  module.
- Discover files with `git ls-files -z go.mod **/go.mod` through the
  repository-pinned `//tools/git` wrapper. This matches the existing
  `tools/bazel_module` and `tools/repo_quality` discovery pattern and skips
  ignored scratch files.
- Configure the version once in `tools/go_mod/cmd/go_mod/BUILD.bazel` and pass
  it as a flag so the fixer and checker share one source of truth.

## Evidence

- `//tools/go_mod/cmd/go_mod:test` covers directive parsing, the rewrite
  behavior, and the whole-workspace freshness check.
- Lowering a module file's version to a value the pinned toolchain does not
  satisfy makes Bazel fail before the fixer can run, so the checker, not the
  fixer, is the durable gate.

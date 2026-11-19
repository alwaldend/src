## 1. Implement the tool

- [x] 1.1 Add `tools/go_mod/cmd/go_mod` with `golang.org/x/mod/modfile`
      parsing, Git discovery, update mode, and check mode.
- [x] 1.2 Add BUILD targets for the update binary and the check test with
      the configured version `1.26.5`.

## 2. Integrate and document

- [x] 2.1 Add `tools/go_mod/README.md` and include the check in
      `//:repo_quality_test`.
- [x] 2.2 Add the `go-mod-version` capability delta.

## 3. Validate

- [x] 3.1 `bazel_agent bazel test //tools/go_mod/cmd/go_mod:test` passes
      (directive parsing, rewrite behavior, workspace freshness).
- [x] 3.2 `bazel_agent bazel test //tools/bazel_module/cmd/update:freshness_test`
      passes.
- [x] 3.3 `bazel_agent bazel test //tools/buildifier:workspace_test` passes.
- [x] 3.4 `bazel_agent bazel test //:repo_quality_test` passes (25 tests).
- [x] 3.5 `bazel_agent bazel build //tools/go_mod/... //tools:docs` passes.

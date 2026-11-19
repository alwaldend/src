## Why

The repository has nine `go.mod` files with no repository-wide rule that keeps
their `go` directive on one configured version. Checking them by hand is easy
to miss when the Go toolchain version changes.

## What Changes

Add `tools/go_mod`, a repository tool that discovers every tracked `go.mod`
file and either rewrites its `go` directive to the configured version or fails
when one differs. The version is configured in the owning BUILD file; the
check target runs as part of `//:repo_quality_test`.

## Capabilities

### New Capabilities

- `go-mod-version`: Keep tracked `go.mod` files on the configured Go version.

### Modified Capabilities

None.

## Impact

- `tools/go_mod`: new Go binary, check test, BUILD files, and README.
- `tools/repo_quality`: the repository quality suite includes the new check.
- `infra/src/openspec/specs/repository`: new baseline requirement and scenario.

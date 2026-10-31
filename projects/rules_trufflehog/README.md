---
title: TruffleHog rules
description: Scan-only Bazel tests with commit-and-line finding allowlists
languages:
  - go
  - bzl
---

`rules_trufflehog` is a standalone Bazel module with a reusable test rule and Go
scan wrapper. It redeclares existing repository dependency pins.
Git metadata uses TruffleHog's generated protobuf Go type at the scanner's
3.99.7 pin. The surrounding CLI output uses ordinary Go JSON, so the wrapper
keeps a small envelope adapter and retains complete original JSON findings.
The consumer supplies its scanner, Git executable, repository marker, and policy.
The monorepo's scan targets and allowlist live in `tools/trufflehog`.

```starlark
load("@rules_trufflehog//:defs.bzl", "trufflehog_test")

trufflehog_test(
    name = "secrets_test",
    scanner = ":trufflehog",
    git = ":git",
    repository_marker = "//:AGENTS.md",
    allowlist = "allowlist.json",
)
```

The marker must be a file at the Git worktree root. The rule resolves its real
location, merges tool runfiles, and runs without generated shell glue.
`HEAD` history is the default; `all_refs = True` includes every local ref.
Both modes disable credential verification and update checks and bypass cached
test results. They do not scan uncommitted or untracked source.

## Exact finding policy

The JSON policy uses schema `trufflehog_allowlist/v1` and a `findings` array.
Each entry contains `detector`, `path`, `commit`, `line`, and `reason`.
All coordinates must match the scanner result exactly. Commit IDs must be full
lowercase Git object IDs; lines must be positive. Paths are canonical
repository-relative Git paths. Reasons must be nonempty.

There are no value hashes, regex matches, mutable refs, or abbreviated commits.
Duplicate tuples and unknown policy fields fail before starting the scanner.
A tuple covers results from that detector at the specified immutable source
line. Findings in a different commit or at another line remain actionable.
An entry absent from a branch's history is valid; stale entries do not fail a scan.

## Output and failures

The wrapper writes newline-delimited JSON finding events containing the complete
original result and a `suppressed` flag, followed by total, suppressed, and
unmatched counts. Approved findings remain available for local inspection.
Scanner diagnostics pass through stderr. Keep original reports in restricted
local storage rather than publishing credential values in a PR.

The wrapper returns zero only after a complete scan with no unmatched findings,
183 for unmatched findings, and 1 for policy, execution, or decoding errors.
A scanner error remains an error when every decoded finding is approved.
Both modes always pass `--no-verification`, `--no-update`, and
`--fail-on-scan-errors`.

Use `@rules_trufflehog//cmd/scan:scan -- --help` for direct flags, or
`//cmd/scan:scan` from the standalone workspace.
Direct invocations require caller-owned `TMPDIR`; Bazel tests use `TEST_TMPDIR`.
The [E2E command fixture](test/e2e/README.md) produces `e2e-report.json` in
undeclared test outputs. `rule_test` exercises the native rule and runfiles
wiring; the repository consumer exercises the pinned scanner over real history.
Maintained work lives in this project's OpenSpec workspace. The shared Bazel
configuration is imported through `.bazelrc`; the scanner pin remains with the
consumer rather than being duplicated inside the module.

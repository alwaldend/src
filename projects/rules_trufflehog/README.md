---
title: TruffleHog rules
description: Scan-only Bazel tests with commit-and-line finding allowlists
languages:
  - go
  - bzl
---

![Scan state machine with policy validation, finding classification, and pass, unmatched, or error outcomes](assets/scan.svg)

## Diagram key

| State                                    | Details                                                                                                                  |
| ---------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| [Scan Git history](#scan-scope)          | Scan `HEAD` history by default, or all local refs on request. Inline ignores apply only to the revision containing them. |
| [Exit status](#output-and-failures)      | Raw scanner status: 0 without findings, 183 with findings. Other statuses fail validation.                               |
| [Read stdout](#output-and-failures)      | Process findings while the scanner runs, continuing through buffered output until EOF.                                   |
| [Apply allowlist](#exact-finding-policy) | Match detector, path, full commit ID, and line. Preserve original JSON, including suppressed findings.                   |
| [Validate scan](#output-and-failures)    | Wait for both stdout EOF and process exit, then check the exit status against finding counts. Any error fails the scan.  |

## Usage

`rules_trufflehog` is a standalone Bazel module with a reusable test rule and Go
scan wrapper. It redeclares existing repository dependency pins.
Git metadata uses TruffleHog's generated protobuf Go type at the scanner's
3.99.7 pin. The surrounding CLI output uses ordinary Go JSON, so the wrapper
keeps a small envelope adapter and retains complete original JSON findings.
The module owns the scanner release pin and binary toolchain. The consumer supplies
a Git executable, repository marker, and policy.
The monorepo's scan targets and allowlist live in `tools/trufflehog`.
The [state-machine source](docs/scan.mmd) is rendered through the shared site pipeline.

```starlark
load("@rules_trufflehog//:defs.bzl", "trufflehog_test")

trufflehog_test(
    name = "secrets_test",
    git = ":git",
    repository_marker = "//:AGENTS.md",
    allowlist = "allowlist.json",
)
```

The default scanner is `@rules_trufflehog//:trufflehog`, selected through the
module's registered execution-platform toolchain. The module registers every
toolchain generated from its lock; archive execution constraints select a compatible
implementation without another registration in the parent repository. The existing
release archive supports Linux x86-64. The [binary lock](third_party/com_github_trufflesecurity_trufflehog/binary_toolchain.json)
owns its version and integrity; the generated Go metadata uses the same version.
Pass `scanner` only for a custom binary or a controlled test fixture. Its executable
and runfiles replace the default scanner. The public CLI target can also run
directly. This module requires Bazel 9.3.1 or later, matching its binary-toolchain dependency.

## Scan scope

The marker must be a file at the Git worktree root. The rule resolves its real
location, merges tool runfiles, and runs without generated shell glue.
`HEAD` history is the default; `all_refs = True` includes every local ref.
Both modes disable credential verification and update checks and bypass cached
test results. They do not scan uncommitted or untracked source.
Inline `trufflehog:ignore` comments apply to the revision containing them. Adding
an ignore now does not suppress a finding in an earlier commit.

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

TruffleHog's raw exit status must be 0 without findings or 183 with findings.
The wrapper checks it after reading stdout through EOF; any other status or
contradictory finding count fails the scan.

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
configuration is imported through `.bazelrc`. The repository CLI alias in
`tools/trufflehog` forwards to the module-owned binary. There is no root third-party
scanner package or separate release pin.

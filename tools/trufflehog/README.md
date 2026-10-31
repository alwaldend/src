---
title: Trufflehog
description: Trufflehog
languages:
  - bzl
tags:
  - bzl_rules
---

The normal `//tools/trufflehog:repo_test` scans the history reachable from the
checked-out `HEAD` and skips unrelated local refs. Run the manual
`//tools/trufflehog:repo_all_refs_test` target explicitly when an audit must
include every local Git ref, such as T3 checkpoint refs from other worktrees.
Both tests use the `external` tag so Bazel reruns them when undeclared Git
history or refs change instead of reusing a previous test result.

The [rules project](../../projects/rules_trufflehog/README.md) supplies the
scan-only wrapper. Both tests disable credential verification and update checks.
They use the rules module's default scanner without specifying a binary label.
`allowlist.json` suppresses only an exact detector, path, full commit ID, and
line tuple with a recorded reason. Files with approved historical findings remain
in the scan. There are no value hashes or regex exclusions.

For a reviewed synthetic credential in current source, prefer a
`trufflehog:ignore` comment on the same physical line, followed by a short reason.
The pinned scanner honors these comments before emitting results. They do not
affect earlier Git revisions: retain exact commit-and-line entries for historical
findings whose source lacks the comment. Do not add a current-commit entry for a
fixture already covered by its inline comment.

Inline-ignored results are absent from the wrapper's events and counts. To inspect
them explicitly, run the pinned CLI with `--no-ignore-tag`, `--no-verification`,
and `--no-update`, selecting the required source scope and retaining output locally.

Output preserves original findings as JSON events with a `suppressed` flag and
ends with total, suppressed, and unmatched counts. A successful repository test
means zero unmatched findings; approved findings are still counted and retained.
Policy errors, incomplete results, and scanner errors fail the test even when
every decoded finding is approved. Add entries only after reviewing the exact
source location and recording the reason.

Use the [repository secret-scan skill](../agents/skills/repo-secret-scan/SKILL.md)
for scan execution, original local finding retention, and triage.

## Links

- Repo: https://github.com/trufflesecurity/trufflehog

The pinned upstream binary and toolchain are owned by
[`rules_trufflehog`](../../projects/rules_trufflehog/README.md). This package's CLI
label forwards to that same binary; the former root third-party package has been removed.

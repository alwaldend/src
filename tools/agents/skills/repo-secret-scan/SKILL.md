---
name: repo-secret-scan
description: >-
  Scan this repository for secrets with the pinned TruffleHog targets, select
  Git-history scope, preserve original local findings, and report actionable
  triage results. Use for repository secret scans and all-ref credential audits;
  use repo-secrets for Vault wiring, rotation, or credential configuration.
---

# Scan the repository for secrets

## Select coverage

Use `al-workspace` before writing scan output and `al-bazel`/`repo-bazel`
for repository commands. Record the current commit and working-tree state.
The [scanner README](../../../trufflehog/README.md) owns the supported
history scopes; its BUILD file owns the consumers, `allowlist.json` owns approved
finding coordinates, and `projects/rules_trufflehog` owns execution and matching.
Inspect these before selecting a target. Both repository tests force
`--no-verification`: they detect credentials without attempting authentication.

- Normal scan: `//tools/trufflehog:repo_test` scans history reachable from
  checked-out `HEAD`, skipping unrelated refs.
- Explicit all-ref audit: `//tools/trufflehog:repo_all_refs_test` includes local
  refs such as other worktrees' T3 checkpoints. It does not fetch remote refs.

Neither Git-history target establishes coverage of uncommitted or untracked
files. If those are requested, use the pinned `//tools/trufflehog:trufflehog`
CLI's filesystem mode with an explicit source scope; do not silently include
ignored build caches, local credentials, or previous scan reports. Check its
`--help` through `al bazel run` before constructing a different scan.

## Run and retain evidence

Run from the Git worktree root. Use ignored, task-owned `out/<task>/` storage
with `umask 077`; point configurable temporary storage there too. For example:

```sh
umask 077
mkdir -p out/secret-scan/tmp
TMPDIR="$PWD/out/secret-scan/tmp" TMP="$PWD/out/secret-scan/tmp" \
  timeout --signal=TERM --kill-after=30s 1200s \
  al bazel test //tools/trufflehog:repo_test \
  --nocache_test_results --test_output=summary > out/secret-scan/bazel.log 2>&1
scan_status=$?
printf '%s\n' "$scan_status" > out/secret-scan/exit-code.txt
```

Use a distinct task directory when earlier evidence must be preserved. For
an all-ref request, substitute only the target. Preserve the test's original
`test.log` in the restricted task directory as `findings.raw.log`, and give
the user its absolute path. Resolve the test log location from this invocation's
Bazel summary; do not assume another checkout's `bazel-testlogs` belongs to it.
Keep raw values out of shell arguments and ordinary progress output.
Force fresh execution when checking changed Git history: a `no-cache` tag
alone does not disable Bazel's local test-result cache.

Record scanner version, commit, invocation, scope and finding policy, completion,
total/suppressed/unmatched counts, scan-only mode, and exit status. The wrapper
preserves each original result inside a finding event with its suppression status;
its final summary records `verification: false`. Confirm that summary
and inspect scan errors: a failed or interrupted invocation is not a clean
scan. The wrapper returns 183 for unmatched findings; Bazel may display
this as `Signal 55`. Execution, policy, and output errors return 1.
Distinguish that finding condition from execution errors
using the actual scan output. A passing invocation without evidence of the
intended coverage is insufficient.

## Inspect public findings without destroying the evidence

A credential in public Git is already exposed, or the match is a false
positive. Redacting the original local evidence does not remove it from Git
and prevents the operator from checking which case applies. Preserve original
values for that inspection; do not automatically overwrite or delete the raw
report before the user has examined it. Public source paths, commit hashes,
detectors, line numbers, and demonstrated synthetic fixture values do not need
blanket redaction. Inspect context before declaring a value synthetic; a test
file name or an unverified detector result is not enough.

Use `repo-secrets` for handling genuine credentials. Local inspection and
publishing a credential to a PR or shared report are different actions: keep
usable credential values in the restricted local report and publish useful
source metadata and conclusions. Redacting a shared report is disclosure
control, not remediation of the repository leak. A real exposed credential
needs revocation or rotation; do not perform those operations or rewrite
history without the user's authorization.

Report each finding's detector, suppression status, source path, line, and
commit, plus its classification and supporting source context when inspected.
Use historical source locations for deleted files rather than inventing links
to current lines. Treat unverified results as unresolved until investigated.
For a reviewed synthetic fixture in current source, prefer a `trufflehog:ignore`
comment on the same physical line as the credential, with a short reason. Native
inline ignores run before JSON output, so these matches are absent from wrapper
events and counts. Use the pinned CLI's `--no-ignore-tag` with `--no-verification`
and `--no-update` for an explicit inspection that includes ignored lines.
Adding a comment does not change older Git revisions: retain historical entries
and do not duplicate an inline-covered fixture with a new current-commit entry.

Approve a historical false positive through an exact `detector`, `path`, full `commit`, and
positive `line` entry in `allowlist.json`, with a nonempty `reason`. Every coordinate
must match. Do not use credential hashes, abbreviated commits, mutable refs, regex
path exclusions, or whole-detector suppression. A later commit or different line
requires its own review. Stale entries may remain when their commits are absent
from the scanned branch. Scanner failures must remain failures after filtering.

Retain the original report for the user's requested inspection, outside tracked
source. Remove task-owned raw copies and secret-bearing generated test outputs
when inspection is complete or the user requests cleanup; keep the safe summary
and leave unrelated output alone.

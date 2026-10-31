## Why

The repository pins TruffleHog 3.94.2. The user requested an update and a new
scan; the publisher's latest stable release observed on 2026-10-31 is 3.99.7.

## What Changes

- Update the existing Linux x86-64 release URL, version, and SHA-256 SRI pin.
- Verify the downloaded binary's version and rerun the existing repository scan.
- Retain original local scan output for the user's inspection.
- Preserve reviewed fixture suppressions through exact finding locations and
  rerun the scan without credential verification.
- Mark both Git-history tests `external` so undeclared Git changes cannot be
  hidden by Bazel's local test-result cache.
- Replace the temporary path exclusions with commit-and-line finding suppression
  through `projects/rules_trufflehog/openspec/changes/commit-finding-allowlist`.
- Share the verified Go build-action sandbox fix in `tools/bazelrc/project.bazelrc`
  while preserving `no-sandbox` on test execution. The user explicitly approved
  including this supporting validation fix in the same pull request.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This dependency refresh preserves scan scope; `skip_specs`
records the absence of a repository requirement change.

## Impact

`third_party/com_github_trufflesecurity_trufflehog/binary_toolchain.json` owns
the pin. `tools/trufflehog` owns the scan targets. The companion procedure is
recorded in `tools/agents/openspec/changes/repository-secret-scan`.
`tools/bazelrc` owns the shared execution-policy setting and its documentation.

Acceptance requires publisher-checksum verification, version 3.99.7 from the
pinned binary, an observed completed scan after the authorized exclusions, and
repository validation. The user classified all four original URI results as
false positives; the original report remains available locally for inspection.
The shared fix must remove `no-sandbox` from Go compilation and analysis actions,
retain it on the freshness test runner, and pass the existing environment-policy
test without inheriting host temporary-directory variables into build actions.

Release evidence: [publisher release](https://github.com/trufflesecurity/trufflehog/releases/tag/v3.99.7)
and [publisher checksums](https://github.com/trufflesecurity/trufflehog/releases/download/v3.99.7/trufflehog_3.99.7_checksums.txt).

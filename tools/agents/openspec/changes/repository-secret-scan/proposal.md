## Why

Repository secret scans need a discoverable procedure that selects the intended
Git scope and preserves original findings for the operator to distinguish
exposed credentials from false positives.

## What Changes

- Add and register `repo-secret-scan` with the repository's Bazel scan commands,
  raw local evidence retention, exit-status interpretation, and finding triage.
- Include offline evaluation cases for normal scans, all-ref audits, and public
  repository reporting.
- Rerun the scan with the separately recorded TruffleHog pin update.
- Document scan-only execution and exact commit-and-line allowlisting through
  `projects/rules_trufflehog`.
- Clarify that first-party nested modules and repeated declarations of existing
  pins do not require new-dependency approval, as explicitly directed by the user.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `project-agents`: Provide a discoverable repository secret-scan procedure.

## Impact

Canonical skill content lives under `tools/agents/skills/repo-secret-scan`;
`.agents/BUILD.bazel` owns its generated discovery link. The supporting
TruffleHog update belongs to the shared repository change
`infra/src/openspec/changes/update-trufflehog`.

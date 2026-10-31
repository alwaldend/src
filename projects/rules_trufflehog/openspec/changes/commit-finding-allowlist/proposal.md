## Why

File-wide exclusions hide unrelated findings. The user requested exact finding
suppression by detector, path, commit, and line, and scan-only repository tests
that never attempt credential verification.

## What Changes

- Replace the path regex file with a versioned JSON finding allowlist.
- Run the pinned scanner through a Go wrapper with verification disabled.
- Preserve original finding records and annotate accepted suppressions.
- Fail on unmatched findings, malformed input, or scanner errors.
- Migrate the ten findings covered by existing path exclusions and update the
  scan skill and repository tests.

## Capabilities

### New Capabilities

- `repository-scan`: Scan Git history and suppress only exact approved findings.

### Modified Capabilities

None.

## Impact

`projects/rules_trufflehog` owns the standalone module, wrapper, public
test rule, and isolated E2E fixture. `tools/trufflehog` owns its repository-specific
allowlist and consumers. The wrapper uses existing Cobra and Go runfiles pins
and TruffleHog's generated Git metadata type at the existing scanner pin.
The new first-party module redeclares existing software without
introducing a new scanner or credential verification, under the user's clarified policy.
The active `tools/agents/openspec/changes/repository-secret-scan` owns the skill
update. This extends the same task PR as the dependency and sandbox fixes.

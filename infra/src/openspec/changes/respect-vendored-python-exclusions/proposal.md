## Why

The first-party import checker currently scans every tracked Python file,
including explicitly vendored source. It needs a vendor boundary independent
of broad formatting exclusions so imported teaching material can stay unchanged.

## What Changes

- Honor explicit `linguist-vendored` and `linguist-vendored=true` classification when selecting Python
  files for import checks, independently of formatting exclusions.
- Keep the existing syntax fixtures and checks for other tracked Python files.

## Capabilities

### New Capabilities

None; this corrects tooling selection under the existing repository boundary.

### Modified Capabilities

None; `skip_specs: true` preserves the existing repository requirements.

## Impact

Only the import checker's file selection changes. First-party Python remains
checked, including `infra/xcp_ng/ansible/files/host_acme.py` despite its formatting
exclusion. Existing syntax fixtures and repository quality gates must pass.

This regular pull request targets `master` directly and supplies the linter
prerequisite for lesson import PR #148. It contains no lesson source or
lesson-specific Git attributes. Marking the imported lesson as vendored belongs
to that separate import.

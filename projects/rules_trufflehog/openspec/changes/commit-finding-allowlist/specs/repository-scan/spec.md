## Purpose

Provide reusable scan-only Bazel tests that retain findings and suppress only
explicitly accepted locations in immutable Git history.

## ADDED Requirements

### Requirement: Exact finding allowlist

The scanner wrapper SHALL suppress a finding only when detector, repository-relative
path, full commit ID, and positive line number all match one approved JSON entry.
Each entry SHALL include a reason. It SHALL NOT match value hashes, patterns,
abbreviated commits, or mutable refs. Malformed and duplicate entries SHALL fail.

#### Scenario: Another finding occurs in an excluded file

- **WHEN** its detector, commit, or line differs from the approved tuple
- **THEN** the finding is reported and the scan test fails

### Requirement: Scan without credential verification

Repository scan tests SHALL pass `--no-verification` to the pinned TruffleHog
executable and SHALL NOT attempt credential verification. HEAD history and
all-local-ref scopes SHALL remain explicit and SHALL bypass cached test results.

#### Scenario: Run a repository scan

- **WHEN** a repository scan test executes
- **THEN** it scans Git history with verification disabled and classifies every
  finding against the allowlist

### Requirement: Preserve evidence and fail closed

The wrapper SHALL retain complete original finding records, mark suppressions,
and report total, suppressed, and unmatched counts. It SHALL fail on scanner
errors, malformed output, missing location metadata, or inconsistent exit status,
including when every decoded finding is allowlisted.

#### Scenario: All findings are approved

- **WHEN** the scanner completes successfully and every finding matches
- **THEN** the wrapper succeeds while retaining the original findings

#### Scenario: A failed scanner prints an approved finding

- **WHEN** the scanner reports an execution error after that finding
- **THEN** the wrapper fails and does not describe the scan as successful

### Requirement: Reuse upstream metadata types

The wrapper SHALL decode Git source metadata with TruffleHog's generated
protobuf Go type at the pinned scanner version. It SHALL retain a CLI JSON
adapter and original finding records rather than treating CLI output as protobuf JSON.

#### Scenario: Decode a Git finding

- **WHEN** the scanner emits `SourceMetadata.Data.Git` in CLI JSON
- **THEN** its Git payload is decoded with `source_metadatapb.Git` and the extracted
  coordinates are validated against the project-owned allowlist contract

### Requirement: Standalone reusable module

The rules_trufflehog module SHALL build and test independently and SHALL accept
consumer-supplied scanner, Git executable, repository marker, and allowlist labels.
It SHALL use existing dependency pins and SHALL NOT reference parent-repository
labels internally.

#### Scenario: Consume the module

- **WHEN** a caller declares a scan test through the project rule
- **THEN** its explicit tools and policy become runfiles without shell-generated
  command execution or hardcoded scanner policy in the rule

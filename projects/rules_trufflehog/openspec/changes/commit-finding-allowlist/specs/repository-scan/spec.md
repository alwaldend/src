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

The rules_trufflehog module SHALL build and test independently and SHALL own the
pinned scanner binary and its default toolchain. It SHALL accept consumer-supplied
Git executable, repository marker, and allowlist labels, with an optional scanner override.
It SHALL use existing dependency pins and SHALL NOT reference parent-repository
labels internally.
The module SHALL register all scanner toolchains generated from its release lock;
archive execution constraints SHALL select compatible implementations. The parent
repository SHALL NOT require a platform-specific scanner registration or retain
a duplicate third-party compatibility package.

#### Scenario: Consume the module

- **WHEN** a caller declares a scan test through the project rule
- **THEN** its explicit tools and policy become runfiles without shell-generated
  command execution or hardcoded scanner policy in the rule

#### Scenario: Use the default scanner

- **WHEN** a caller omits the scanner attribute
- **THEN** the test uses the module-owned pinned scanner through its registered
  execution-platform toolchain, and the module CLI target selects the same binary

#### Scenario: Override the scanner

- **WHEN** a caller supplies an explicit scanner for a fixture or custom build
- **THEN** the test uses that executable and its runfiles instead of the default

#### Scenario: Resolve through module registration

- **WHEN** a parent repository uses the default scanner without registering it again
- **THEN** Bazel selects a compatible toolchain from the module's generated registrations
  and the remaining repository CLI alias resolves to the module-owned binary

### Requirement: State-machine documentation

The project README and landing SHALL begin with a diagram rendered from one
canonical Mermaid state machine. The diagram and its linked documentation key
SHALL describe policy validation, scan-only execution, original finding
preservation, exact allowlist matching, inline-ignore scope, completion validation,
and exit statuses 0, 183, and 1. The site SHALL use
the shared light/dark diagram theme and publish both variants.
The diagram SHALL give the TruffleHog subprocess a nested state machine with
start and end markers, keep wrapper processing outside that boundary, and
distinguish the scanner's exit status from the wrapper's final policy result.
State labels SHALL remain short actions or exit statuses. A compact README key,
linked from the landing, SHALL carry the detailed scope and matching rules and
link to their owning documentation sections.
The diagram SHALL model scanner execution and stdout processing as concurrent
activities. Processing a finding SHALL NOT leave or re-enter the scanner state.
Final validation SHALL require both stdout EOF and process exit, regardless of
which occurs first.

#### Scenario: Inspect the scan workflow

- **WHEN** a reader opens the project README or landing
- **THEN** the first diagram explains scan outcomes and shows that errors still
  fail when findings have been suppressed

#### Scenario: Trace an approved finding

- **WHEN** TruffleHog emits only allowlisted findings and exits with status 183
- **THEN** the diagram shows the wrapper validating that exit and returning zero
  because no unmatched findings remain

#### Scenario: Read workflow details

- **WHEN** a reader needs the scan scope or exact allowlist coordinates
- **THEN** the diagram key supplies those details without paragraph-sized state
  labels and links to the corresponding documentation sections

#### Scenario: Process several findings

- **WHEN** the scanner emits several findings before exiting
- **THEN** the output branch repeats classification without leaving or restarting
  the single scanner subprocess

#### Scenario: Drain output after scanner exit

- **WHEN** the subprocess has exited but stdout contains unread findings
- **THEN** output processing continues through EOF before final validation

#### Scenario: Wait for exit after stdout closes

- **WHEN** stdout reaches EOF before the subprocess exits
- **THEN** final validation waits for the process exit as well

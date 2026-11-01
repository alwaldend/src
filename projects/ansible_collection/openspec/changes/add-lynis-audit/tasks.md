## 1. Scanner and role

- [x] 1.1 Pin the current stable CISOfy source archive and verify its SHA256,
      archive layout, and Bazel archive target.
- [x] 1.2 Implement the bounded audit and evidence collection role; verify the
      supported upstream options, YAML, and packaged role files without a live run.

## 2. Handoff

- [x] 2.1 Validate the OpenSpec change and narrow Bazel packages; record evidence
      and provide the target, variables, and validation limitations to the parent.

## 3. Parent-reviewed role corrections

- [x] 3.1 Record failure cases before implementation.
- [x] 3.2 Preserve private receipts, transport diagnostics, and remote recovery
      paths; collect available evidence and gate cleanup on successful retention.
- [x] 3.3 Require a finished positive-test retained report and accepted exit
      status for completion; retain partial artifacts on failure and timeout.
- [x] 3.4 Apply the shared YUM exclusions and document precise audit coverage,
      noninteractive options, and the supported remote Python prerequisite.
- [x] 3.5 Validate with isolated E2E fixtures where practical, Bazel packaging,
      and strict OpenSpec validation; return the uncommitted diff to the parent.

## 4. PR simplification

- [x] Use native temporary-directory permissions and direct archive transfer,
      and consolidate repeated receipt fields while preserving recovery checkpoints.
- [x] Verify warning, timeout, incomplete-report, and collection-failure behavior
      against the simplified role and retain the exact-candidate evidence.

## Native simplification

- [x] Run Lynis with native async polling, remove its stdout exit protocol, and retain scratch when terminal status is unknown.
- [x] Verify successful and failed execution, private retention, owned cleanup, and the packaged candidate.

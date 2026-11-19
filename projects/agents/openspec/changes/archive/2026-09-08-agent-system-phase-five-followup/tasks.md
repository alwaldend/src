# Historical acceptance tasks

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline. The source criterion revision is 1; each checkbox is supported by the [followup-baseline-001 review](provenance/source/attempts/followup-baseline-001/attempt.yaml).

## 1. Accepted criteria

- [x] 1.1 `baseline-coverage` (revision 1): baseline-coverage: The repository emits a bounded CoverageMatrix from normalized routing cases, binding the exact capability-catalog digest and distinguishing emitted cases from the complete skill universe.
      Evidence: [baseline-coverage.md](provenance/source/attempts/followup-baseline-001/evidence/baseline-coverage.md).
- [x] 1.2 `writable-fixtures` (revision 1): writable-fixtures: Representative Git, Bazel, forge, Terraform, Ansible, secret, and runtime trajectories have deterministic writable fixtures and recorded coverage gaps.
      Evidence: [writable-fixtures.md](provenance/source/attempts/followup-baseline-001/evidence/writable-fixtures.md).
- [x] 1.3 `live-comparison-policy` (revision 1): live-comparison-policy: Scheduled live comparisons remain outside ordinary tests and bind model, skill, catalog, fixture, and judge identities.
      Evidence: [live-comparison-policy.md](provenance/source/attempts/followup-baseline-001/evidence/live-comparison-policy.md).
- [x] 1.4 `optimization-adoption` (revision 1): optimization-adoption: Each adopted optimization traces a measured baseline, predeclared threshold, reviewed proposal, owner-local change, regression, fallback, and retirement rule.
      Evidence: [optimization-adoption.md](provenance/source/attempts/followup-baseline-001/evidence/optimization-adoption.md).

## Preserved limits

The accepted attempt result still names remaining work. Its final review marks all four criteria pass; evidence explicitly records Terraform, Ansible, Vault, and Bazel end-to-end fixture gaps. Import preserves both without treating gaps as implemented.

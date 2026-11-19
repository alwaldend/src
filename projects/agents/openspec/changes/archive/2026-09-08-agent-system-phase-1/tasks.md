# Historical acceptance tasks

This change was imported as completed historical work on 2026-09-08. The source outcome is `achieved` and execution is `paused`. Checked tasks reproduce its accepted review; they are not a new validation of current code. Historical deltas were not applied to the independently documented current baseline. The source criterion revision is 3; each checkbox is supported by the [phase1-completion-001 review](provenance/source/attempts/phase1-completion-001/attempt.yaml).

## 1. Accepted criteria

- [x] 1.1 `shared-contracts` (revision 1): Shared v1alpha1 contracts round-trip deterministically and reject malformed identity, unknown effectful operations, widened authority, and incompatible information flows.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.2 `registered-universe` (revision 1): A cheap report-only completeness check covers the closed registered agent-reachable universe and names missing or unclassified entries.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.3 `information-policy` (revision 1): Public, secret, and personal-information classes remain independent from repository visibility, build-consumer, and publication policy axes.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.4 `scratch-isolation` (revision 1): Task and run scratch is namespaced; concurrent tasks do not collide and stale workers fail expected-revision publication.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.5 `safe-operations` (revision 1): Every removed mutating alias has a safe explicit replacement and provider mutation ownership is unambiguous.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.6 `resource-baseline` (revision 1): A revision-bound Phase 1 baseline records numeric ceilings for correctness, unsafe actions, calls, context size, cold and warm time, universe coverage, and reused checks without estimating unavailable measurements.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).
- [x] 1.7 `exact-candidate-validation` (revision 1): All affected packages, schemas, declarations, generated checks, documentation, and goal records pass focused validation against the exact delivered candidate.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md), [result.md](provenance/source/attempts/phase1-completion-001/result.md).
- [x] 1.8 `encountered-bug-policy` (revision 1): Repository-wide policy requires primary agents to fix small bounded bugs encountered in repository tooling or the affected project instead of silently working around them, while subagents remain scoped and substantial redesigns or rewrites are reported separately.
      Evidence: [phase1-validation.md](provenance/source/attempts/phase1-completion-001/evidence/phase1-validation.md).

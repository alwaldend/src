# Phase 1 completion attempt

## Bindings

- Goal: `agent-system-phase-1`
- Source commit: `5acf09c8ecb1576cc53a11078db4ead847cdc3eb`
- Goal generation: 2
- Lifecycle generation: 7
- Goal resource version: `12`
- Criteria revision: 3
- Criteria digest:
  `sha256:d8c600b75d19ff670e5df73801c1db81761b8890840c78050fab9a534d3a47cd`
- Goal-state digest:
  `sha256:62db8ac202b22a64c98e3967af88bc6142844890e21f6bc03c3023d3deed375a`

## Target defect

Phase 1 has shared types but still lacks a closed registered-universe report,
complete controller/runtime scratch isolation, explicit replacement of the
unnamed Terraform apply alias, and a revision-bound numeric resource
baseline. Those gaps prevent honest Phase 1 acceptance and therefore block
Phase 2.

## Plan

1. Add owner-local declarations and a cheap report-only checker that names
   every registration authority, classifies legacy gaps explicitly, validates
   project lifecycle and generated-artifact declarations, and fails fixtures
   for missing or unknown entries.
2. Replace Cordis worktree-global scratch with explicit task/run namespaces
   and a bounded manifest; retain Bazel-managed action/test temporary storage.
3. Remove the unnamed Terraform `apply` entry while preserving the explicit
   `.apply` replacement and add a regression check.
4. Add a sanitized read-only `bazel_agent doctor` report and tests for source,
   pin, profile, platform, scratch, and stale-install fields.
5. Record and validate a criteria-revision-bound numeric resource baseline.
6. Run every Phase 1 criterion, the fixed regression set, affected package
   tests/builds/lints, Buildifier, goal validation, and delivery verification
   against one frozen candidate.

## Decision review

Proceed with owner-local declarations plus a derived report. Reject a central
hand-maintained fact catalog because it would duplicate component authority
and contradict the roadmap. Phase 1 remains report-only: legacy and unavailable
entries must be explicit classifications, not silently omitted.

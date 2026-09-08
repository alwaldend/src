# Phase 3 acceptance evidence

Candidate: branch `t3code/continue-agent-system-1` (PR 43), rebased onto
`origin/master` `5e7d2802` after the ast-grep base advance.

## impact-planning

- `tools/agents/plan` package: deterministic `ImpactPlan` digest over
  canonical inputs; profiles changed/fast, workspace, fresh/evidence,
  full/audit, diagnose.
- `tools/agents/cmd/agent_system plan` subcommand (`plan_cmd.go`) joins the
  planner with the topology catalog and affected-path resolution.
- `tools/agents/api/v1alpha1/plans.go` defines the plan contract and
  validation.

## effect-admission

- `tools/agents/admission` package: provider gateways admit action contracts
  against exact authority, subject, environment, pre-state, and budgets;
  prepare/validate/authorize/execute/verify for remote writes; cheap read
  path; rejection of unknown actions, stale pre-state, wrong environment,
  weakened safety flags, and missing cancellation budgets.
- Fault-regression coverage in `admission_test.go` (wrong env, stale
  authority, secret output, symlink escape, redirect, infinite loop,
  cancellation, concurrent writers, cleanup failure).

## validation-sets

- `tools/agents/evidence` package: immutable candidate-bound `ValidationSet`
  and `EvidenceAssertion` with exact candidate, profile, check identities,
  sanitized arguments, working scope, provider/config/toolchain/policy
  digests, clean pre/post state, results, coverage, duration, output bounds,
  limitations, and task-local raw-log digests; canonical-JSON digest.
- `tools/agents/api/v1alpha1/validation_sets.go` and `contracts.go` define
  the contracts; `tools/agents/api/v1alpha1/phase3_validation.go` wires
  validation sets into the shared API.

## evidence-assertions

- `EvidenceAssertion` applies validation sets to an exact criterion revision
  and semantic verdict; delivery consumes validation sets and goals consume
  both; the same immutable set supports multiple assertions without mutation.

## receipt-applicability

- Receipt applicability uses relevant dependency-slice digests; tree-bound
  checks survive message-only rewrites; commit-bound checks do not; changes
  to base, tree, config, policy, contract, toolchain, generator,
  environment, or remote generation make evidence stale; global catalog
  churn triggers reevaluation not a rerun.

## structured full-repo-check

- `projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go`
  emits incremental versioned JSON plus Markdown, workspace/phase selectors,
  exact input binding, target-universe count,
  zero/unexpected-reduction rejection, resume-only-identical-inputs.

## Verification

- `bazel_agent test //tools/agents/api/v1alpha1:all
  //tools/agents/plan:all //tools/agents/admission:all
  //tools/agents/evidence:all //tools/agents/cmd/agent_system:all
  //projects/agents/skills/full-repo-check:run_full_repo_check_test` — 7/7
  pass on the rebased candidate.
- `bazel_agent test //tools/repo_delivery/...` — all tests pass.
- Catalog drift gates: goal, index, policy, topology, workspace-check,
  action, capability all regenerate and pass after the rebase.
- `//:buildifier_test` pass; `git diff --check` clean.

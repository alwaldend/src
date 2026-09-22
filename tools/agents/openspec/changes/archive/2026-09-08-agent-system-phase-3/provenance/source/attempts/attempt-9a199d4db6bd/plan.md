# Phase 3: add advisory planning, admission, and reusable evidence

Goal: agent-system-phase-3. Depends on phase-2 (bounded catalogs, runtime
isolation, context capsule).

## Deliverables in dependency order

1. **Shared API types** (`tools/agents/api/v1alpha1/plans.go`,
   `validation_sets.go`, `contracts.go`):
   - `ImpactPlan` (3A): profile, capabilities, effects, targets, checks,
     coverage gaps, cost class, escalation, evidence applicability, digest.
   - `ValidationSet` (3C): exact candidate, profile, check identities,
     sanitized args, working scope, digests, clean pre/post state, results,
     coverage, duration, output bounds, limitations, raw-log digests.
   - `EvidenceAssertion` (3C): applies validation sets to an exact criterion
     revision and semantic verdict.
   - `AdmissionRequest`/`AdmissionDecision` (3B): action contract admission
     against authority, subject, environment, pre-state, budgets.

2. **Planner (3A)**: `tools/agents/cmd/agent_system plan` subcommand + package
   `tools/agents/plan`. Profiles `changed/fast`, `workspace`,
   `fresh/evidence`, `full/audit`, `diagnose`. Deterministic digest over
   canonical inputs. Target selection via topology catalog + affected paths.

3. **Admission (3B)**: `tools/agents/admission` package enforcing
   prepare/validate/authorize/execute/verify for remote writes, cheap read
   path, unknown action rejection, stale pre-state rejection. Unit tests with
   fault fixtures (wrong env, stale authority, secret output, symlink escape,
   redirect, infinite loop, cancellation, concurrent writers, cleanup
   failure).

4. **ValidationSet emitter (3C)**: `tools/agents/evidence` package producing
   immutable candidate-bound validation sets and evidence assertions; digest
   via canonical JSON; applicability by dependency-slice digests. Delivery
   consumes validation sets (roadmap phase 4 wiring).

5. **Structured full-repo-check (3D)**: extend
   `projects/agents/skills/full-repo-check/scripts/run_full_repo_check.go` to
   emit incremental versioned JSON + Markdown, workspace/phase selectors,
   exact input binding, target-universe count, zero/unexpected-reduction
   rejection, resume-only-identical-inputs.

## Acceptance

Each deliverable has focused unit/regression tests plus the goal criteria
acceptance signals from the roadmap. All affected packages pass
`bazel_agent test`, buildifier, and `git diff --check` before delivery.

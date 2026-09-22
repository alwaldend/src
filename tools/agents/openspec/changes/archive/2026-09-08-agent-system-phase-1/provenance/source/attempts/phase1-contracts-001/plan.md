# Phase 1 shared-contract attempt

## Bound input

- Goal: `agent-system-phase-1`
- Goal generation: 2
- Lifecycle generation before activation: 6
- Criteria revision: 3
- Source base: `6780d53a69e32064d648e6a04f1c0cecd7d713fd`
- Prior attempt: none

## Target

Implement the smallest repository-internal `v1alpha1` contract package needed
by later Phase 1 declarations and Phase 2 catalogs. Also retain the bounded
goal-store repair encountered while attaching this attempt.

## Plan

1. Make an attempt-free project goal valid after a clean Git checkout and
   prove that its first checkpoint recreates the untracked directory safely.
2. Add `tools/agents` as the repository-internal owner for shared agent-system
   contracts, with role-based `api/v1alpha1` source placement.
3. Define narrow identity, authority, budget, path-policy, effect,
   information, retention, completeness, availability, and evidence-
   applicability types without component payload schemas.
4. Provide strict validation and deterministic JSON round-trip behavior.
5. Add rejection fixtures for malformed identities, unknown effectful
   operations, authority widening, and incompatible information flow.
6. Run focused package tests, builds, formatting checks, goal validation, and
   inspect the exact diff.

## Affected criteria

- `shared-contracts`
- `information-policy`
- `exact-candidate-validation`
- `encountered-bug-policy`

## Fixed regressions

- Existing goal API and filesystem-store tests.
- Repository Buildifier test if BUILD metadata changes.
- The project-goal record validates from a clean-checkout directory shape.

## Reset condition

Reset the approach if the common types require importing a component schema,
runtime provider, daemon, or production dependency, or if deterministic
round-trip requires a second serialization authority.

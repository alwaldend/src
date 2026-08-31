# Bazel agent scratch-isolation attempt

## Bound input

- Goal: `agent-system-phase-1`
- Goal generation: 2
- Lifecycle generation: 7
- Criteria revision: 3
- Goal resource version before publication: 10
- Source base: `6780d53a69e32064d648e6a04f1c0cecd7d713fd`
- Prior attempt: `phase1-contracts-001`

## Target

Remove the shared host temporary directory from the mandatory Bazel agent
path. Bazel actions and tests must use their Bazel-managed temporary contracts;
task-owned host tools will use explicit task/run manifests in a later slice.

## Plan

1. Remove `bazel_agent` workspace discovery, `out/tmp` creation, and
   `TMPDIR`/`TMP`/`TEMP` rewriting.
2. Remove the agent Bazel configuration that propagates ambient temporary
   paths into repository rules, normal actions, host actions, and tests.
3. Update runner documentation to distinguish Bazel-managed temporary storage
   from explicit task-owned host scratch.
4. Add regression coverage proving the runner preserves the ambient
   environment byte-for-byte while retaining signal-transparent Bazel exec
   arguments.
5. Run the Bazel-agent package tests, affected Bazelrc checks, Buildifier,
   lint, and exact diff hygiene.

## Affected criteria

- `scratch-isolation`
- `exact-candidate-validation`

## Fixed regressions

- The shared-contract and goal-store suites from `phase1-contracts-001`.
- Root Buildifier and repository Bazelrc validation.

## Reset condition

Reset if any supported no-sandbox repository updater demonstrably depends on
ambient `TMPDIR` injection and has no owner-local scratch contract.

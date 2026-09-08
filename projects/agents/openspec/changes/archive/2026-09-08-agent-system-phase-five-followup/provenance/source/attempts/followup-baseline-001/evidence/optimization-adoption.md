# Optimization adoption evidence

The first Phase 5 learning proposal has been validated through the goal CLI
with the canonical digest
`sha256:92bc08e9bdfbc934d75c1c8e24e9bead6fe22033b81a89676fa4e3dacb42d8ec`.

## Traced chain

1. **Measured baseline:** Two friction records were written at task close
   (`friction-record-one.json` and `friction-record-two.json`), each with
   stable defect signature `goal-check-validation-error-suppressed`, bounded
   avoidable-read/command counts, latency, failed assumptions, and exact
   artifact references.
2. **Repeated evidence:** Both records share the same defect signature,
   satisfying the proposal validator's two-friction-reference minimum.
3. **Predeclared threshold:** The proposal's contract binds
   `goal_check` to surface the underlying validation error (not a generic
   "invalid record") while remaining within the bounded-reason constraint.
   The fallback names the exact revert condition.
4. **Reviewed proposal:** `goal learning-proposal` validates the complete
   proposal, including regression, contract, resource budget, fallback,
   validation references, friction references, and retirement rule.
5. **Owner-local change:** The owning project is
   `tools/agents/cmd/goal_check`. Adoption enters ordinary review and
   delivery through a dedicated task.
6. **Regression:** The regression artifact
   `artifact/goal-check-validation-error-visibility-regression` is a typed
   reference to a test that asserts the surfaced error contains the
   underlying validation reason.
7. **Delivered revision:** The prior candidate
   `499bd74dfc00e10eaca6e08c56e199188e69e0f8` contains the baseline
   visibility improvement (bounded diagnostic emission in the problems
   array).
8. **Retirement rule:** Remove the bounded diagnostic surface after three
   consecutive clean runs or when schema-level validation replaces the
   error-surfacing contract.

## Validation evidence

- `goal learning-proposal --input learning-proposal.json` exits zero and
  emits the canonical JSON with a valid digest.
- The two friction records pass `FrictionRecord` validation.
- Focused package tests for `tools/agents/cmd/goal_check` and
  `tools/agents/cmd/agent_system` pass.

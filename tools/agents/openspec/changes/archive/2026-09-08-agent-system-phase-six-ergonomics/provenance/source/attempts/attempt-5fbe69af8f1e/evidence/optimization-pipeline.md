# Phase 6B optimization pipeline evidence

Each adopted optimization below follows the proposal, fallback, and
retirement trace. The primary fallback is to retain the prior behavior and
record the defect in the next friction baseline. Retirement triggers when two
consecutive friction baselines record zero avoidable commands for that defect
signature, or zero avoidable reads where the fix is read-only.

## Adopted optimizations (trace to measured baseline defects)

1. Catalog updater discoverability (60 s baseline):
   this aggregate branch made each stale-catalog check name its exact
   `bazel_agent bazel run` updater target. Regression: the five check
   package tests assert their messages.
   Proposal: name the exact runnable target at the point of stale-catalog
   failure instead of requiring package-target discovery.

2. Delivery validation-command mapping (45 s baseline):
   this aggregate branch added `suggested_validation_commands` to the prepare
   report, derived deterministically from affected labels. Regression:
   `TestSuggestedValidationCommands`.
   Proposal: map delivered labels to the narrowest validation command without
   a separate query discovery pass.

3. Commit convention enforcement (goal-workflow archetype):
   this aggregate branch made delivery validate regular Git subject and
   trailer conventions before publication.
   Proposal: validate the full convention at prepare time so rejection
   happens before remote rewrite.

4. Bazel test-log path discoverability (35 s baseline):
   validated learning proposal
   `learning-proposal/bazel-test-log-path-undiscoverable`
   (canonical digest `sha256:8f939eb64d2b958d410743822fa31d6225d1eeed2c277a45858d643b1f35d6d8`).
   Change: `test:agent --test_output=errors` streams failing test output
   inline. Live regression: deliberately broken assertion produced the
   failure line immediately in streamed output with no log-path hunt
   (`out/phase6-logs-proposal/regression-output.txt`, task scratch).

5. Runner strict subcommand (goal-workflow archetype):
   this aggregate branch made `bazel` the validated Bazel entry point.
   Proposal: require an explicit `bazel` subcommand so runner commands and
   repository commands cannot collide.

## Fallback and retirement

Fallback and retirement rules for the test-log optimization are recorded in
the canonical learning proposal. Retirement triggers when two consecutive
friction baselines record zero avoidable reads for this defect signature.

## Remaining defects (not adopted this attempt)

- goal-check validation error suppression (180 s): validator-level
  improvement, bounded follow-up.
- skill-validation error specificity (30 s): partially addressed by path
  context added in this branch; remaining scope is message-level detail.

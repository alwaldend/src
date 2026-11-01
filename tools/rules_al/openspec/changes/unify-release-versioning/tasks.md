## Implementation

- [x] 1. Consolidate versioning ownership, expose the shared release CLI/bootstrap, and update maintained references and skill discovery.
- [x] 2. Bind portable first-party bundles to source identity and provide per-project bundle targets.
- [x] 3. Implement staged, verified, immutable SSH publication and publication-bound extraction/activation.
- [x] 4. Add explicit deletion with active-site protection and ordinary republishing.
- [x] 5. Correct diverged patch publication and replace the simulated lease with an actual lock and guarded ref updates.
- [x] 6. Validate isolated release lifecycles, format and document the operator workflow, and prepare the branch for delivery.

## Review corrections and module extraction

- [x] 7. Add end-to-end cases for filtered destination activation, payload symlinks, and module-owned Bazel bundle generation.
- [x] 8. Fix item selection and symlink verification, and correct website deployment instructions.
- [x] 9. Extract self-contained rules_release with unprefixed public rules, providers, and defs.bzl.
- [x] 10. Migrate parent consumers, module registration, documentation, and release-prefixed skill discovery.
- [x] 11. Validate standalone and parent behavior, update locks through Bazel, run the workspace check matrix, and prepare guarded delivery.
- [x] 12. Verify linked-worktree metadata through the CLI, read the common Git directory, and prepare repeated delivery gates.

## AL repository lifecycle

- [x] 13. Add built-CLI E2E cases for version inspection, bundle generation, delivery provider inspection, and Bazel forwarding before implementation.
- [x] 14. Extract rules_al and merge AL, release/versioning, delivery, and Bazel runner implementations under direct AL subcommands.
- [x] 15. Migrate Bazel rules, plugin API consumers, integration tests, module registration, and repository entry points.
- [x] 16. Move lifecycle skills to al-prefixed names and regenerate discovery.
- [x] 17. Generate dependency metadata, validate standalone and parent outputs, and exercise the full workspace matrix.
- [x] 18. Review the final implementation, record acceptance evidence, and prepare guarded delivery. Publication state belongs to delivery receipts.

## Approved base integration

- [x] 19. Reconcile the fetched synchronization change while preserving immutable releases, retained-site rollback, and explicit deletion; verify ancestry and mount refusals through the CLI.
- [x] 20. Prepare final candidate checks and updated pull-request metadata for guarded delivery. Delivery receipts own validation and publication state.

## PR review follow-up

- [x] 21. Extend built-CLI fixtures for cancellation, Cobra version flags, mounted files, and transfer-free identical retries before implementation.
- [x] 22. Propagate cancellation through version commands and subprocesses, and replace the legacy version flag parser with native Cobra commands.
- [x] 23. Reject file mountpoints before deletion and verify existing publications before staging or transferring retry payloads.
- [x] 24. Build a minimal external consumer without parent Bazel configuration, inspect its bundle, and record repeatable acceptance evidence.
- [x] 25. Format and review the fixes and select aggregate validation gates for guarded delivery. Delivery receipts own validation and publication state.

## Terminal and recovery follow-up

- [x] 26. Verify foreground Git authentication and cancellation through a controlling terminal before correcting process-group handling.
- [x] 27. Automate the independent consumer build and verify its manifest and payload in a retained result artifact.
- [x] 28. Add checkpointed artifact-operation receipts for completed and uncertain steps, including interrupted and partial operations.
- [x] 29. Select final representative-output and aggregate validation gates, review the changes, and prepare guarded delivery. Delivery receipts own validation and publication state.

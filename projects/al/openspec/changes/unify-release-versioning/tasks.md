Completed tasks below retain the history of earlier implementation stages. The
current design and specifications include the later simplifications; the pending
tasks at the end track that remaining work.

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

## Product description after rebase

- [x] 30. Describe the rebased AL command-line surface in owner-local documentation, with a pilot, foundations, shared vocabulary, and source-linked feature behavior.
- [x] 31. Add verification checklists, run safe local observations, and collect suspected defects without changing product behavior or deploying releases.
- [x] 32. Check documentation consistency and representative output, and select exact-candidate delivery gates. Delivery receipts own validation and publication outcomes.

## Product-description bug fixes and design assessment

- [x] 33. Extend behavior fixtures before implementation for standalone command execution, redirected output, conflicting plugin exports, and rebase without premature publication.
- [x] 34. Fix runfiles requirements, output creation/truncation, and plugin export conflicts while preserving runtime cleanup and existing Bazel consumers.
- [x] 35. Make direct rebase local-only and require the ordinary prepared validation/publication workflow for the rebased candidate; update its procedure and evidence.
- [x] 36. Assess typed Bazel configuration providers against existing Lua consumers and identify demonstrated lifecycle gaps, without implementing an unrequested configuration migration.
- [x] 37. Update affected product behavior and triage, verify representative outputs, and select final delivery gates. Receipts own validation and publication outcomes.

## Release inspection and publication preview

- [x] 38. Define failure and success fixtures before implementation for offline selection preview and read-only local/remote release inspection, retaining repeatable output and unchanged-storage evidence.
- [x] 39. Implement artifact publication preview through the existing selection and validation owner, reporting identity, resolved destinations, selected payloads, and requested steps without contacting destinations.
- [x] 40. Implement read-only release inspection through the existing verification owner, distinguishing absent, matching, conflicting or invalid publications and retained/active website state without remote writes.
- [x] 41. Update operator and product-description documentation, inspect representative JSON output, and include the new behavior in exact-candidate delivery gates.

## Independent AL project

- [x] 42. Move the independent AL project to projects/al, retain the rules_al Bazel module identity, and update owning declarations, consumers, documentation paths, and generated metadata through their workflows.
- [x] 43. Identify and remove concrete unconfigurable repository assumptions from reusable runtime behavior; supply this repository's policy and integration through configuration while preserving current guards.
- [x] 44. Adapt the existing external-consumer check and affected validation to the new location. Architectural portability is required; exhaustive execution of every workflow in another repository is not.
- [x] 45. Reconcile the product description and source references, inspect representative outputs, and prepare the updated aggregate for guarded delivery.

## Public naming and restored landing

- [x] 46. Prefix AL-owned rules and rule macros with al_ and providers with Al, align source filenames, and migrate all current consumers and examples.
- [x] 47. Restore the AL landing content and standalone-module registry mapping, then regenerate the project inventory.
- [x] 48. Render one dark WebP state-machine diagram with horizontal Development, Source publication, Release, and optional Website publishing subgraphs, one bold single-line command per transition between meaningful states, and editable Mermaid source.
- [x] 49. Verify public API consumers, standalone and site outputs, documentation references, and select exact-candidate delivery gates for the existing PR. Delivery receipts own validation and publication results.
- [x] 50. Remove retired consolidation stubs and unnecessary compatibility entry points, migrate remaining current consumers, and verify the surviving owner boundaries.

## Explicit website commands and Git metadata cleanup

- [x] 51. Adapt existing CLI and SSH fixtures before implementation for dedicated website extraction and activation, offline previews, retired step-option refusals, and operation receipts.
- [x] 52. Expose release website extract and activate, keep release publish upload-only, and migrate configured repository entry points and active consumers without compatibility aliases.
- [x] 53. Remove the unused full-history Git JSON command, rule, target, and site input; retain minimal release identity and verify the optional packed Git changelog path.
- [x] 54. Reconcile documentation and the state-machine artifact, inspect representative outputs, and select exact-candidate delivery gates. Receipts own validation and publication outcomes.
- [x] 55. Replace the stamped Bazel wrapper with ordinary release configuration invoking installed AL, verifying clean identity and dirty-source refusal without nested Bazel.
- [x] 56. Remove al bazel and migrate runtime callers, preserving doctor diagnostics and delivery guards; move agent defaults to managed user configuration and document setup.
- [x] 57. Reconcile current skills, product description, and four-phase lifecycle diagram with ordinary Bazel; validate representative outputs and prepare the revised candidate for delivery. Receipts own final validation and publication outcomes.
- [x] 58. Replace version release-start with release create for atomic local branch/tag creation without changing HEAD; use ordinary Git checkout explicitly in workflows and behavior fixtures, with no AL checkout command or old alias.
- [x] 59. Reconcile documentation and diagram labels/layout with release creation, Git checkout, and the retained bundle handoff; validate and prepare the revised candidate for delivery.

## Git-only preparation and affected tests

- [x] 60. Adapt existing real-Git publication fixtures before removing delivery validation plans, the validation ledger, validate/continue/deliver commands, and --validated-head; retain exact preparation, publication, review, and recovery guards.
- [x] 61. Keep prepare Git-only and report checks_required after a publication-time rebase; verify separate ordinary checks and receipt-bound publication, including recovery after a review observation failure.
- [x] 62. Define affected-test behavior fixtures before implementing al bazel test: obtain Git changes, query Bazel reverse dependencies within the selected scope, execute the selected tests, and prevent incomplete graph information from silently skipping tests.
- [x] 63. Reconcile command metadata, canonical skills and generated discovery, operator and product documentation, and the regular four-group lifecycle diagram with independent test execution and simplified delivery.
- [x] 64. Validate current source and representative result artifacts, inspect the final diagram, and prepare the aggregate candidate for the existing pull request. Keep ordinary check evidence separate from Git preparation and publication receipts; those receipts own the final remote publication state.

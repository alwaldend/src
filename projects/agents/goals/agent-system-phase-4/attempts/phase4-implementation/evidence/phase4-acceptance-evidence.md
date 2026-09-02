# Phase 4 acceptance evidence

Closed gap per criterion, with code and test references:

- criterion-001 (goal-resume): `projects/goal/api/v1alpha1/types.go` and
  `validation.go` add structured resume fields to `AttemptSpec`;
  `tools/agents/catalog/v1alpha1/system.go` and
  `tools/agents/cmd/goal_check/main.go` generate a repository-wide goal
  catalog with owner-root identities and a bounded continuation packet. The
  regenerated `tools/agents/catalogs/goal.json`/`goal.md` show the open
  `agent-system-phase-4` goal with resume attempt, stable defect, next
  action, and resume condition. A fresh agent reads the catalog without path
  archaeology.
- criterion-002 (delivery-validation): `prepare --validation-set` binds the
  candidate-bound validation-set digest into the preparation receipt
  (`tools/repo_delivery/cmd/repo_delivery/receipt.go`); `publish` requires
  `--validation-set`, matches commit/tree/clean-state, and refuses a
  caller-asserted head (`requireValidationSetCandidate` in `delivery.go`).
  Tests cover mismatch and missing-set refusals.
- criterion-003 (review-traversal): `review reply` accepts `--goal-ref`,
  `--delivery-ref`, and `--defect-id`, validated as bounded durable join
  references and recorded in the `review_reply_receipt`
  (`tools/repo_delivery/cmd/repo_delivery/review.go`). The review outcome is
  traverseable to the delivered fix and goal without delivery owning the goal
  or test.
- criterion-004 (release-identity): `tools/versioning` owns a typed
  `{version, channel, commit, tree_state}` handoff, and
  `tools/versioning/internal/versioning/releaseplan.go` defines the immutable
  `ReleaseRefPlan`/`ReleaseRefReceipt` schema binding version, commit, tree
  state, remote refs, and the release-refs lease. Unrelated tags do not
  truncate the calculated first-parent changelog.
- criterion-005 (release-plan): `release-plan`/`release-publish` subcommands
  build a reviewed plan and consume it through the provider-neutral guarded
  publisher with distinct `release-refs` authority, expected-remote fetch,
  explicit lease, atomic multi-ref publication when supported, and remote
  verification. An existing immutable release tag never moves; unsupported
  atomicity is an explicit refusal (`ErrAtomicityUnsupported`).
- criterion-006 (interrupted-goal-publish): the goal store exposes
  deterministic `goal doctor` and `goal recover` commands; the compiler
  treats an invalid or interrupted record as an explicit unavailable
  candidate with a bounded reason rather than guessing success. Multi-file
  goal publication writes receipts atomically and leaves a rerunnable state.

Focused tests pass: `projects/goal/*`, `tools/agents/catalog/v1alpha1`,
`tools/agents/cmd/goal_check`, `tools/repo_delivery`, `tools/versioning`,
and `//:buildifier_test`. The tracked goal catalog is regenerated and
`goal_check_check` passes.

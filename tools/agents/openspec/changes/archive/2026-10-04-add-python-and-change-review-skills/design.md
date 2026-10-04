## Context

See proposal.md for motivation and specs/project-agents/spec.md for acceptance.
The existing PR head 0fffe5b489f887020c1b11a4e815f00f14a961ff already packages
repo-python and has passing delivery tests and lint. This change adds the
missing work record and the requested review procedure after that implementation.

## Goals / Non-Goals

The skills must be discoverable and useful for actual bounded repository work.
A review request does not authorize fixes or external messages. This change
adds no review service, external dependency, or automatic merging capability.

## Decisions

- Use change-review, as corrected by the user, so reviews cover PRs, commit
  ranges, and staged or working-tree candidates rather than only forge objects.
- Keep the every-change OpenSpec policy in AGENTS.md. repo-workspace selects
  the required record before editing, while openspec owns the record lifecycle.
- Reuse the existing skill_library and offline Promptfoo pattern. Live evaluation
  needs a tool-capable fixture repository; offline validation proves loading,
  not review quality. A bounded review of this candidate supplies observed
  procedural evidence without claiming independent or exhaustive validation.
- Fix the stale Python support reference found in the earlier review: the
  metadata owner is tools/py/pyproject.toml, while root owns tool settings.

## Risks / Trade-offs

- Retrospective records could imply prior planning: proposals explicitly state
  when the records were created and identify the previous candidate.
- Broad standards claims could invent requirements: change-review traces each
  claim to applicable policy, owning contracts, or reproducible behavior.
- Mandatory records could inflate trivial work: reuse active matching records,
  keep artifacts small, and use skip_specs when behavior does not change.

## Evidence and continuation

Prior validation covered head 0fffe5b489f887020c1b11a4e815f00f14a961ff. New skill
content and policy edits invalidate those checks for the next candidate.
Next: implement and package change-review, regenerate discovery, run OpenSpec
validation and skill checks, review the aggregate, then archive and deliver.

The aggregate change review found that the new direct course project lacked an
owner-local OpenSpec workspace required by the existing repository baseline.
projects/mlops_2026/openspec now owns the placeholder contract and record,
with explicit Bazel documentation and validation wiring. No application
behavior was introduced by this correction.

## Pre-archive acceptance

On 2026-10-04, the pre-archive staged source snapshot was
59b946c9281710a97c0e198d2d2283c8aa7169a1, compared against base
2f1266ab33c053c8cabd82db568c366ccc8f96ae. The working tree matched the staged
source during the final quality and lint checks. The selected quality,
dependency, manifest, and checker run passed all 30 tests; skill/OpenSpec
validation passed eight tests. Documentation and skill builds and semantic
lint passed. Offline eval results establish harness loading, not model behavior.

Using change-review, inspection found and corrected the stale Python support
metadata reference and missing course-owner OpenSpec workspace. No further
actionable defect was found in the reviewed aggregate. This was an agent review
of its own change, not independent behavioral evaluation; automated remote
review was previously unavailable because of its usage limit.

Archive projections and final formatting will change the tree. Final acceptance
must therefore use a new prepared candidate and passing delivery receipt,
including quality, consumers, semantic lint, and archived-spec validation.
Publication state remains with Git and delivery receipts, not this document.

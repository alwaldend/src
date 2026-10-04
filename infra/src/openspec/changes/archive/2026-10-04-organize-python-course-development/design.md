## Context

See proposal.md for scope. The initial implementation and follow-up dependency
move were already delivered in PR #124 at
0fffe5b489f887020c1b11a4e815f00f14a961ff before this record was created.
Source ownership remains in projects/mlops_2026, tools/py, tools/gazelle, and
root policy; the linked agent change owns skill-specific requirements.

## Goals / Non-Goals

Preserve existing Python checker and dependency behavior through ownership
refactoring. No application code, live deployment, or manifest relocation is
part of this change.

## Decisions

- Keep shared tool settings at root for ancestor discovery and explicit checker
  targets; tools/py owns only project metadata, dependency declarations, and lock.
- Retain root gazelle_python.yaml for default discovery. The pinned plugin
  supports python_manifest_file_name with a relative path if relocation is later
  requested; default discovery alone would not cover sibling trees after a move.
- AGENTS.md owns mandatory OpenSpec coverage for every change. The linked agent
  change updates selection procedures and evaluation cases.
- Exclude the documentation-only course placeholder from deploy_heads because
  no releases:head target exists.

## Risks / Trade-offs

- Moving dependency files can leave dangling labels: inspect references and
  validate lock, manifest, affected packages, and existing checker consumers.
- Retrospective evidence can overclaim: distinguish the prior validated head
  from the new candidate and record the omission explicitly.

## Evidence and continuation

At 0fffe5b489f887020c1b11a4e815f00f14a961ff, requirements and Gazelle consistency,
repository quality, skill discovery, Black/isort consumers, and semantic lint
passed. Shared settings existed only at root and dependency versions were
unchanged by the move. Next: validate the new policy delta and final candidate,
archive linked changes only after supported acceptance, and publish PR #124.

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

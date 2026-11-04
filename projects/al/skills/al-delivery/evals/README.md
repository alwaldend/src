---
title: Repository delivery evaluations
---

# Repository delivery evaluations

This suite describes the behavioral contract for safely finalizing a feature
branch and its pull request. Its required offline Bazel target validates the
Promptfoo configuration, referenced cases, and staged skill without making a
model call. It does not establish that an agent follows the procedure.

The ordinary workflow separates Git preparation, caller-run Bazel checks, and
receipt-bound publication. Cases preserve repository quality and relevant
semantic lint as mandatory gates, owner/consumer target selection, representative
output inspection, and failure diagnosis. Bazel owns check configuration and
cache reuse; a new Git commit alone does not require disabling its cache.
Preparation receipts bind Git and forge state, not test outcomes. The caller
must honor failed checks without a parallel AL check-plan format or results
ledger.

The source-entrypoint case covers a host without an installed AL command.
The stale-notes case migrates retired delivery commands to the current API.
The preparation case keeps build passthrough and check execution out of Git
preparation. The publication-rebase case reruns ordinary required Bazel commands
on the updated candidate before another explicit publication attempt. Recovery
first verifies known remote state; it does not blindly push again after an
uncertain outcome.

The base-selection cases distinguish independent PRs with merge prerequisites
from an explicitly requested stack. History cases preserve task ownership,
complete local/remote progress, exact remote leases, and evidence-backed minor
conflict resolution. They retain aggregate commit and type-specific disclaimer
requirements. No check failure authorizes merging a prerequisite or bypassing
a required gate.

Local review cases require explicit `review-agent` selection for the complete
candidate diff before each PR creation or update, including documentation-only
changes and changed publish-time rebases. An unchanged publication recovery
reuses its local review evidence. Local review does not replace remote feedback
inspection, and remote review does not replace local review.

Remote review cases keep inspection and mutations behind `al delivery review`.
`review address` uses the clean prepared candidate and verified publication,
then preserves the existing reply checkpoint, one-use receipt, and expiry.
Cases reject automatic duplicate replies after unknown outcomes, silent
resolution, stale provider expectations, and unsupported completion claims.
A behavior-changing review fix requires fresh proportional scrutiny as well
as the relevant Bazel checks.

Compatibility cases preserve honest Forgejo capability reporting, a clean
worktree/index without autostash, same-repository GitHub topology, and explicit
SSH or gh-authenticated HTTPS transport boundaries. HTTPS cases require an
isolated host-scoped helper and consistent receipt-bound transport selection,
without importing global Git configuration. Remote-review cases require observed completion for the
final head when the provider exposes it; an unobservable result is reported
honestly after bounded inspection.

Authentication cases distinguish API preflight from Git transport, compare
sanitized active-account status across execution contexts, inspect JSON state
rather than exit status alone, and preserve approval rejection. They reject
unnecessary login changes and token copying when the stored login works in an
approved context. Offline configuration validation does not exercise a live
credential store or sandbox; the change record retains the bounded live checks.

A live target is omitted because representative delivery requires inspecting
and changing Git history, pushing a remote ref, interacting with a forge, and
responding to live review state. A read-only, tool-free model evaluation cannot
reproduce those operations safely or verify their postconditions. Delivery
behavior needs isolated integration fixtures or review of a real, authorized
operation in addition to this offline configuration check.

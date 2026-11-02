## Context

See proposal.md for the live failure evidence. Resolution consumes a reply receipt before the provider call. The provider marks transport read failures as pre-mutation errors, but advancing-epoch errors from coherence validation are unmarked, so no-mutation failures destroy usable authority. Thread inventory currently fetches each thread's comments separately. The manual CLI exposes all provider expectation fields to callers.

## Goals / Non-Goals

Retain exact candidate validation, full inventory binding, one-use bounded authority, repository/PR ownership, and push leases. Reduce call count and expectation copying. Never retry a public mutation automatically or treat semantic state changes as transport failures.

## Decisions

Proceed after decision review: fix the existing read-error classification and batch nested comment reads rather than removing consistency checks. Removing guards would reduce code but permit stale-head or changed-thread actions; arbitrary retrying mutations risks duplicate public replies.

Add review address as orchestration over existing primitives. Require a validated published preparation receipt, derive fresh PR/thread expectations, and write its reply receipt beside the body. On re-entry, verify the saved reply's identity and body and resume resolution. Low-level commands remain available for existing consumers. A failed reply without a trusted receipt remains an unknown outcome and must not be automatically reposted.

Embed comments in thread nodes, using the same comment validator and limits. Additional pages retain a targeted reread and metadata consistency checks. Preserve raw-byte limits across all calls.

Wrap only a positively identified advancing epoch before resolution in the existing pre-mutation read error. Restoration preserves original bytes and expiry; malformed data, changed inventories/threads, expired authority, and post-mutation reads remain terminal.

## Risks / Trade-offs

- Larger nested responses → retain byte/node/request budgets and paginate comments.
- Automation hides expectations from the caller → bind the workflow to the published validated candidate and derive guarded mutations from fresh observations.
- GitHub lacks atomic compare-and-swap → preserve pre/post verification and explicitly report unknown outcomes; batching shortens but cannot eliminate the race window.

## Verification and ergonomics evidence

- `review-epoch-authority-loss` (live): two #133 resolutions failed before mutation on advancing review epochs. The typed pre-mutation classification retains the original authority window; unrelated errors remain terminal.
- `review-thread-fetch-amplification` (fixture): the two-thread/two-page inventory fixture expects two requests instead of four, with comment values and outer pagination retained. Direct comment continuation and malformed-state fixtures retain their checks.
- `review-guard-handoff` (fixture): real Git preparation, recorded check results, publication, reply, read-failure recovery, and resolution produce `review-address.json`; retries must leave exactly one public reply. Unknown outcomes retain an attempt checkpoint and stop.
- The preimplementation check attempt was blocked during LLVM extraction by a full Bazel cache disk. Only the already-merged scaffold worktree's disposable Bazel cache was cleared. The first executable run caught relative versus absolute validation-log paths and an obsolete decoder fixture; both were corrected before candidate delivery.
- Promptfoo offline validation covers configuration and skill packaging, not live model decisions. Live GitHub mutation success remains unproven by the deterministic forge fixture.

- `review-thread-fetch-amplification` (live): the new read-only reader returned 34 threads and 67 comments for #128 at head `6393d1c62f225bce5fe12cf05b6cc7e4ed741640`, matching the complete provider projection. No live review mutation was performed for this check.

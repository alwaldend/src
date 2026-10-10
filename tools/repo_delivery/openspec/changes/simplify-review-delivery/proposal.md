## Why

Resolving an unchanged review comment failed twice because a GitHub inventory epoch advanced before any resolution mutation, consuming receipts and requiring duplicate public replies. Routine comment handling also requires copying multiple IDs and digests, and fetching each thread separately increases latency and the window for incoherent reads.

## What Changes

- Add a guarded reply-and-resolve entry point that derives expectations from fresh inventory and the published validated candidate.
- Classify unstable pre-resolution reads as retryable reads and retain the original unexpired receipt when no mutation occurred.
- Batch review-thread comments with their inventory page while preserving pagination and byte/node/request bounds.
- Update documentation and the delivery skill to prefer the simpler workflow.

## Capabilities

### New Capabilities

- `review-delivery`: Review workflow automation with exact-candidate guards, bounded reads, and no duplicate replies on retryable read failures.

### Modified Capabilities

None.

## Impact

Owns tools/repo_delivery Go CLI, provider adapter, fixtures, documentation, and skill guidance. No dependencies or live infrastructure configuration change. Exact-commit validation, branch ownership checks, complete inventory verification, and explicit push leases remain required.

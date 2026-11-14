## Context

See proposal.md for motivation. The existing evaluator submits native CLI jobs ahead of retrieval, but owns rental allocation and bounded case matrices. This tool operates an already-running endpoint through a caller-supplied CLI command and must have no instance lifecycle or subject-specific behavior.

## Goals / Non-Goals

Goals are finite scheduling of ready directions, deliberate evaluation before successors, an owned queue, resumable provenance, configurable meaningful novelty, safe interruption, and reviewable media. Model selection, prompt invention, credentials, server deployment, and artistic acceptance remain with the caller or producer.

## Decisions

- Add `cmd/experiments` and `internal/experiments` under the existing ComfyUI tools owner. Use Go, the existing Cobra dependency, standard-library subprocesses, and the pinned CLI protocol. Do not import the provisioning evaluator or introduce a second authentication flow.
- A JSON configuration supplies literal CLI and producer argv, queue lookahead, timeouts, finite limits or explicit continuous mode. The producer receives a sequence number and private experiment directory on stdin, and returns a workflow path, optional display notes, and JSON pointers to exclude from workflow identity. A producer can change any model, parameter, input, or subject without changing the runner.
- Snapshot each workflow before submission. Hash canonical JSON after declared exclusions, such as sampling seed and save prefix, and reject a repeated identity with a bounded number of producer attempts. The contract is identity after explicit exclusions, not a claim to infer arbitrary workflow semantics.
- Keep one executing job plus the configured pending buffer outstanding. Observe completion independently and replenish before downloads or report generation. Downloads and media hashing must not block submission.
- Persist a private atomic journal and the native submission receipt before proceeding. A single advisory lock prevents two controllers using the same run directory. Resume known prompt IDs without resubmission. An unacknowledged submission is ambiguous and must stop automatic submission instead of retrying it.
- `status` reads the journal; `stop` requests a drain, with an explicit owned-job cancellation option. Cancellation passes only recorded prompt IDs to the existing CLI. Ordinary interruption preserves remote jobs for resumption. No command provisions, restarts, stops, or destroys an instance.
- Write original media, hashes, workflow snapshots, producer notes and a local gallery with collapsed data. Treat execution and artistic quality as separate facts. Authentication wrappers and private diagnostics remain outside published source.
- The skill owns agent decisions: record stable directions and lineage, queue one revision per ready direction, inspect finished outputs and explicitly scrap, refine or fork before producing successors. Retired directions do not recur automatically. Buffering applies only to the finite ready set; idle capacity is preferable to quota-filling generation. Caller-owned direction records and opaque producer notes carry this procedure without a subject-specific scheduler or new runner schema. Its examples are neutral and contain no current artwork or model-specific directions.

## Risks / Trade-offs

- Slow producers, a stopped endpoint or pending review can leave the endpoint idle. Retain configurable lookahead within ready directions, report observed state and bounded failures, and do not bypass evaluation to maintain utilization.
- Identity exclusions can be misconfigured: validate every pointer, record the exclusion set, and document that meaningful parameters must remain in the identity.
- A lost submission acknowledgment cannot be safely retried: retain the ambiguous record and refuse automatic resubmission.
- Arbitrary CLI diagnostics or workflows can contain protected data: use private artifact directories, bounded captures, and summaries that omit command arguments and raw diagnostics.
- Native UI and API workflows differ: preserve their bytes as caller-owned inputs and delegate conversion to the existing CLI.

## Migration Plan

Keep the authorized task-local feed running during implementation. Validate the generic command against isolated CLI fixtures, then use it with an external task-local producer on the existing endpoint. Stop refilling the temporary feed only after the generic runner has a queued buffer; let its already-owned jobs finish and retain their outputs. Publish the reusable source and skill through the normal delivery workflow.

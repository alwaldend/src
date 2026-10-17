## Context

Export projects one native snapshot into config paths and strips native identities. State owns binding construction, including note-type members. See proposal.md for motivation.

## Goals / Non-Goals

Generate portable bindings without installing state or introducing a collection import command. Preserve config-only behavior with an empty state destination.

## Decisions

- Export records native IDs while assigning config paths, avoiding a second snapshot or name-based matching. Pass configs, IDs and that snapshot to `state.Export`, which owns complete binding construction and validation. Note-type member binding remains private to state and is shared with import.
- Extend the regular App `ColExport` and internal exporter `Export` with an explicit state destination argument; an empty string selects config-only behavior. Each service retains one export operation.
- Encode TOML through the existing codec. Stage the state file beside its destination and rename after successful config writes, preventing truncated state publication. Follow existing export semantics for partial config output; no cross-file transaction is promised.
- Reject source aliases at the application boundary and configuration filename/directory conflicts before output writes. Document the state file as belonging to the source collection, not a portable identity mapping for unrelated collections.
- Resolve output paths through existing ancestors and dangling symlink targets with one export-owned helper shared by App and destination validation. Compare canonical generated destinations before writing, even when neither target exists yet.

## Risks / Trade-offs

- Partial config output on filesystem failure → retain existing state output and document fresh-directory export for clean snapshots.
- Existing stale configs remain in reused directories → onboarding examples use a fresh directory.
- Concurrent external writers → retain the offline exclusive-access requirement.

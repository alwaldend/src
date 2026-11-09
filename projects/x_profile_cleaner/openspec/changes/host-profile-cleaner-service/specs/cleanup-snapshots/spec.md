## Purpose

Reuse account-owned discovery data for local policy evaluation and accept bounded
X archive imports as an alternative to spending credits on timeline collection.

## ADDED Requirements

### Requirement: Explicit discovery and zero-request reevaluation

Refresh from X, Import X archive, and Evaluate policy SHALL be distinct actions.
Snapshots SHALL be immutable, persisted per account, and reusable within retention
and validity limits. Evaluation SHALL use cached data without X API calls, including
identity calls. Missing data SHALL require an explicit refresh or replacement import.

#### Scenario: Change a policy repeatedly

- **WHEN** a logged-in user evaluates several policies against the same usable snapshot
- **THEN** each evaluation produces its own plan without any X request or additional X charge

### Requirement: Visible snapshot provenance and coverage

A snapshot SHALL record its account, source, collection/import time, available
fields, discovery scope, and known coverage limitations. The UI SHALL expose source,
age, window, and validity before evaluation and apply. Narrower policy windows SHALL
restrict data before graph evaluation. Incomplete acquisition SHALL NOT replace a
usable snapshot or publish an executable partial plan. No source SHALL imply
complete history or graph membership solely because import or collection succeeded.

#### Scenario: Refresh runs out of budget during pagination

- **WHEN** additional pages are required but cannot be admitted
- **THEN** refresh fails without publishing partial coverage as a complete snapshot
- **AND** the previous usable snapshot remains available and dispatched reads remain accounted

#### Scenario: Preservation seed is missing

- **WHEN** a policy explicitly names a seed absent from the selected snapshot
- **THEN** planning fails rather than silently widening deletion or fetching the seed

### Requirement: Safe archive import

The service SHALL accept supported X archive account, post, and repost records as
snapshot data. It SHALL verify matching account identity, validate normalized data,
and preserve the distinction between repost activity and source post IDs. Archives
SHALL NOT grant X authorization. Import SHALL NOT execute bundled code, follow
external links, or ingest unrelated private archive content. Unsupported, ambiguous,
conflicting, or over-limit input SHALL fail without replacing an existing snapshot.

#### Scenario: Import without X authorization

- **WHEN** a logged-in user imports a valid matching archive without an active X grant
- **THEN** the service permits local planning without X requests
- **AND** applying the resulting plan still requires matching temporary cleanup authorization

#### Scenario: Wrong account or malicious archive

- **WHEN** an archive has another account ID, unsafe entries, executable payloads, or invalid records
- **THEN** import rejects the unsafe or unsupported input without executing code or exposing another account's data

### Requirement: Invalidation after mutations

Dispatching a mutation SHALL invalidate affected cached snapshots and other
unconsumed plans, including uncertain outcomes. Evaluation of dirty snapshots SHALL
require explicit refresh or replacement import. Invalidation SHALL NOT erase the
original consumed execution or prevent explicit recovery of its unfinished actions.

#### Scenario: Lost deletion response followed by reevaluation

- **WHEN** a deletion might have succeeded but its response is unknown
- **THEN** other plans cannot apply using the affected stale snapshot
- **AND** explicit resume retains the original selection and known completions

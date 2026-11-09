## Purpose

Use bounded standard X export archives as the only inventory for policy evaluation.
Workers retrieve archives directly from configured sources without payload uploads.

## ADDED Requirements

### Requirement: Explicit discovery and zero-request reevaluation

Archive replacement and policy evaluation SHALL be distinct actions. The cleaner
SHALL NOT provide live timeline acquisition or refresh.
Snapshots SHALL be immutable, persisted per account, and reusable within retention
and validity limits. Evaluation SHALL use cached data without X API calls, including
identity calls. Missing data SHALL require an explicit replacement import.

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

#### Scenario: Archive acquisition fails

- **WHEN** a configured source cannot supply a valid bounded matching export
- **THEN** evaluation fails without publishing partial coverage as a complete snapshot
- **AND** previous usable data remains unchanged and no X request is sent

#### Scenario: Preservation seed is missing

- **WHEN** a policy explicitly names a seed absent from the selected snapshot
- **THEN** planning fails rather than silently widening deletion or fetching the seed

### Requirement: Configured archive source references

The job API SHALL accept a configured source name and an opaque provider-specific
archive key. The worker SHALL retrieve the object directly after acknowledging the
job ID; the backend SHALL NOT stream archive or normalized inventory payloads.
Provider configuration SHALL remain operator-owned. The initial provider SHALL be
a confined local directory. The provider interface SHALL permit future object-store
support without changing job submissions; S3 implementation is outside this increment.

#### Scenario: Select one of several configured sources

- **WHEN** a job selects a valid source name and object key
- **THEN** that provider supplies the archive directly to the worker
- **AND** unknown sources, unsafe local paths and source escapes fail closed

### Requirement: Safe archive import

The service SHALL accept supported X archive account, post, and repost records as
snapshot data. It SHALL require `data/account.js`, `data/tweets.js` and supported consecutive
multipart tweet files, verify matching account identity, validate normalized data,
and preserve the distinction between repost activity and source post IDs. An unresolved repost SHALL NOT be treated as an own post or assigned a guessed
source target; a requested repost view containing one SHALL fail. Archives
SHALL NOT grant X authorization. Import SHALL NOT execute bundled code, follow external links, or decode unrequested
archive content. Generic scripts SHALL be able to read other top-level YTD datasets,
including multipart likes and following, within shared byte/record limits. These
datasets SHALL remain read-only input for plans of X-account API operations. Unsupported, ambiguous,
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
require explicit replacement import. Invalidation SHALL NOT erase the
original consumed execution or prevent explicit recovery of its unfinished actions.

#### Scenario: Lost deletion response followed by reevaluation

- **WHEN** a deletion might have succeeded but its response is unknown
- **THEN** other plans cannot apply using the affected stale snapshot
- **AND** explicit resume retains the original selection and known completions

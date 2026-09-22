# Goal publication and recovery audit

## Audit binding and verdict

- Goal: `agent-system-phase-2`
- Goal generation: 2
- Lifecycle generation: 4
- Criteria revision: 2
- Active attempt: `goal-publication-001`
- Source base / merged PR 39:
  `1ca06a0fc9697e8fb212b32d99d9ad3b996ea76e`
- Audit mode: source-and-test inspection only; no maintained source or canonical
  goal record was changed and no test was run.

**Verdict: the current store fails the `goal-recovery` criterion.** It has
atomic individual-file replacement, a per-goal lock, an expected-resource-
version guard, preflight validation, and reliable fail-closed validation of
many torn states. It does not have a durable publication intent, a structured
incomplete-state diagnosis, or a recovery entry point. Once a checkpoint is
interrupted after its first canonical rename, ordinary `checkpoint`, `render`,
`show`, `attach`, `promote`, and `validate` all enter through
`loadAndValidate` and cannot repair an invalid record. The CLI exposes no
`doctor` or `recover` command (`projects/goal/cmd/goal/command.go:73-85`).

The smallest compatible repair is a filesystem publication-intent protocol,
not a database or daemon: stage exact after-images, atomically publish one
bounded intent containing before/after digests and operation identity, replay
or classify it under the existing goal lock, then remove it only after the
complete record validates. `doctor` must classify `stable`,
`incomplete-recoverable`, `committed-projection-stale`, and `conflict`; a
`recover` operation must idempotently finish the staged after-image or, before
any canonical target differs from its before digest, discard the intent and
retain the exact prior record.

## Existing primitive and observability

`atomicWrite` creates a sibling `.goal-write-*`, writes, closes, calls the
single test hook, and renames it over one target; it does not sync the file or
parent directory (`projects/goal/internal/fsstore/store.go:948-978`). The hook
therefore injects failure immediately before any individual file rename, and
the same hook is called before immutable snapshot and directory publication
renames (`projects/goal/internal/fsstore/integrity.go:153-187`,
`projects/goal/internal/fsstore/checkpoint.go:542-552`). It cannot inject after
a rename, at process exit, or at a durability barrier. Atomic visibility for a
process interruption is supported; power-loss durability is not claimed by
the implementation.

Every normal goal read first deletes recognized `.goal-write-*`,
`.goal-immutable-*`, and `.goal-attempt-*` residue, then demands a completely
valid record (`projects/goal/internal/fsstore/integrity.go:21-79`,
`projects/goal/internal/fsstore/store.go:809-912`). This cleanup is useful for
prepublication residue, but it destroys a staged new attempt after a crash and
has no operation intent from which to decide whether to roll forward. The
validator's observable torn-state errors include:

- criteria current/snapshot mismatch, a missing snapshot, or snapshots beyond
  the goal's current revision (`integrity.go:190-253`);
- missing or malformed attempt resources and artifact digest mismatch
  (`store.go:848-867`);
- a goal active pointer different from the open attempt set
  (`store.go:888-897`);
- an open attempt against a non-open/non-active or stale lifecycle
  (`store.go:899-906`); and
- an achieved goal whose acceptance pointer does not resolve to the exact
  accepted result (`checkpoint.go:617-639`).

These are generic validation errors, not a stable incomplete-state schema.
Store methods can return a `GoalReference` together with a post-commit error,
but command wiring discards the reference and returns only the error
(`command.go:262-266`). The resource version is consequently observable only
because selected errors interpolate it as text.

## Complete publication-boundary matrix

### Whole-directory creation operations

| Operation | Publication sequence and injectable boundaries | Observable interruption state | Current behavior and recovery |
| --- | --- | --- | --- |
| `init` | Build `goal.yaml`, `criteria.yaml`, revision 1, and `README.md` in hidden `.goal-init-*`, then hook and rename the directory (`store.go:218-267`). Each staged atomic write and the final directory rename are hookable. | Before final rename: target absent; SIGKILL can leave a hidden staging directory. After rename: whole valid record visible. | Ordinary errors remove staging and return no reference. Retry can recreate when target is absent; after a lost acknowledgement it reports “already exists,” not idempotent success. No crash-residue adoption/GC. |
| `promote` | Copy and rewrite a complete hidden `.goal-promote-*`, validate it, hook, then one rename (`migrate.go:126-195`). | Target absent or whole validated target; source unchanged. Hidden residue can survive a crash. | A matching existing target is idempotent by promotion digests (`migrate.go:93-124`). No residue recovery. |
| `migrate` | Build and validate `.goal-migrate-<id>-*/<id>`, hook, recheck target, redigest source, then rename staged goal (`migrate.go:541-625`). | Target absent or whole target; source unchanged. Hidden staging parent/candidate can survive a crash. | Matching target is treated as success; absent target is rebuilt on retry. Invalid target is preserved and blocks retry. Details and gaps are below. |

For these operations the internal multi-file writes are not canonical until a
single directory rename. They already meet process-interruption atomic
visibility at the target, apart from lost-acknowledgement semantics for
`init`, hidden residue, and migration's no-overwrite race.

### Checkpoint of an existing goal

The common checkpoint preflights a complete prospective record and renders
the future projection before any canonical write (`checkpoint.go:219-241`).
It then publishes `goal.yaml` first as the declared commit point, attempt
content next, and `README.md` last (`checkpoint.go:242-297`). Every listed
rename is independently hookable.

| Mutation | Ordered canonical boundaries | State after interruption | Current error / retry behavior |
| --- | --- | --- | --- |
| Lifecycle or outcome only | `goal.yaml`; `README.md` | Before goal rename: exact prior valid record. After goal rename: valid new canonical record with absent/stale noncanonical projection. | Goal failure returns empty reference. README failure returns committed reference and error text with new version (`checkpoint.go:257-295`; test `integrity_test.go:247-280`). `render` repairs the projection. |
| New open attempt | hidden complete attempt; `goal.yaml` with new active pointer; attempt-directory rename; `README.md` | Before goal: prior record plus disposable hidden stage. After goal/before attempt: invalid, active pointer has no open attempt. After attempt: valid intended canonical state; projection may be stale. | Post-goal failure reports committed version but `loadAndValidate` blocks retry. Tests assert the invalid state but no recovery (`integrity_test.go:132-245`). |
| New attempt closed immediately | hidden complete closed attempt; intermediate `goal.yaml` with active pointer; attempt-directory rename; final `goal.yaml` at the same resource version; `README.md` | After first goal: invalid missing attempt. After attempt rename: invalid because active pointer names a closed attempt. After final goal: valid intended state. | Both gaps report the committed version in-process; after interruption, validation blocks every normal retry. Tests cover both gaps and only assert fail-closed behavior (`integrity_test.go:155-168,186-239`). |
| Existing open attempt content update | `goal.yaml`; optional `result.md`; each new evidence file in sorted order; `attempt.yaml`; `README.md` (`checkpoint.go:555-585`) | After goal but before any attempt byte: often a valid record with an advanced token but none of the intended attempt change. After result/evidence but before manifest: invalid artifact manifest. After manifest: valid intended state. | A pre-attempt failure can silently leave a different valid identity that is neither the exact prior state nor the intended state. Later failures report committed version and validation fails. Test coverage injects goal and final attempt-manifest failures with one result, not every result/evidence boundary (`integrity_test.go:12-130`). |
| Close existing attempt, optionally close/achieve goal | `goal.yaml` already clears active pointer and may changes lifecycle/acceptance; result/evidence; closed `attempt.yaml`; `README.md` | Immediately after goal: invalid because the old attempt is still open but no longer the active pointer (and an achieved pointer may not resolve). Sidecar gaps additionally produce digest mismatch. After manifest: valid. | No recovery operation; the same validation gate blocks retry. No boundary-complete fault test for existing-attempt close/achieve. |

Staging a new attempt writes `plan.md`, `result.md`, sorted evidence, then
`attempt.yaml` inside a hidden directory (`checkpoint.go:465-523`), so failures
there are noncanonical. The final attempt-directory rename is separately
hookable (`checkpoint.go:542-552`). Existing attempts instead expose every
sidecar rename directly before the new manifest makes their digests coherent
(`checkpoint.go:555-585`).

### Criteria replacement

Criteria replacement publishes four independent files in this order:

1. immutable `criteria-revisions/<new>.yaml`;
2. `criteria.yaml`;
3. `goal.yaml` advancing criteria and lifecycle revisions; and
4. noncanonical `README.md`
   (`projects/goal/internal/fsstore/store.go:727-763`).

The observable states are:

- before snapshot rename: exact prior valid record;
- after snapshot, before `criteria.yaml`: invalid because snapshots extend
  beyond the goal's current revision;
- after `criteria.yaml`, before `goal.yaml`: invalid because current criteria
  and goal status disagree;
- after `goal.yaml`: valid intended canonical record; and
- after a README failure: valid canonical record with stale projection.

Only the final projection failure has a fault test and committed-version
assertion (`integrity_test.go:282-330`). Snapshot and current-criteria boundary
failures return an empty reference because the implementation has not yet
declared the goal commit, even though canonical files have changed. A direct
in-process retry could reuse an identical immutable snapshot, but after a
crash all retries fail at `loadAndValidate` before reaching that idempotence.

### Relationship replacement

Relationship replacement writes `goal.yaml` then `README.md`
(`projects/goal/internal/fsstore/graph.go:168-196`). The first rename is a
complete valid canonical update; projection failure returns the committed
reference and is tested (`graph_test.go:342-373`). `render` can repair it.
There is no torn canonical multi-file state because README is noncanonical.

### Single-file operations outside the goal record

`render` changes only replaceable `README.md` (`render.go:12-45`) and `attach`
changes only one session-binding YAML under its own path lock
(`store.go:435-528`). They do not need a multi-file recovery transaction.

## Fault-injection and test coverage assessment

Current tests establish useful but insufficient properties:

- before-`goal.yaml` failure preserves that file (`store_test.go:612-638`);
- new-attempt failures at first goal, directory publication, and immediate-
  close finalization advance or preserve the expected token and fail closed
  (`integrity_test.go:132-245`);
- one existing-attempt result/manifest split fails closed
  (`integrity_test.go:12-130`);
- checkpoint, criteria, and relationship projection errors report the
  committed version (`integrity_test.go:247-330`, `graph_test.go:342-373`);
- recognized temporary residue is deleted (`integrity_test.go:332-373`); and
- concurrent writers serialize and one stale token loses
  (`store_test.go:200-239`).

Missing tests correspond directly to unsupported guarantees:

- every result and evidence rename, especially multiple evidence files;
- existing-attempt close, achieve, abandon, and lifecycle-change combinations;
- criteria snapshot and `criteria.yaml` failures;
- process death rather than an ordinary returned hook error;
- lost acknowledgement after each successful final rename;
- stable machine-readable doctor classification for every intermediate state;
- repeated recovery after interruption at every recovery rename;
- conflicting/tampered before or staged after images; and
- file and directory durability (if power loss enters the contract).

The existing test that rejects lock/WAL/journal/transaction names in the goal
catalog (`store_test.go:882-944`) encodes the old design assumption and must be
replaced or narrowed if a durable publication intent is introduced. The
intent is canonical recovery metadata, not a lock and not a background
service.

## Smallest recoverable filesystem protocol

Use the current per-goal lock and resource-version CAS. Add one bounded,
versioned `GoalPublication` intent plus a hidden same-filesystem staging
directory. No daemon, database, shared scheduler, or generated README input is
needed.

1. Under the goal lock, require and validate the exact prior resource version.
2. Construct all exact after-images in hidden staging and validate the complete
   prospective record, as today. Record an operation ID, operation kind,
   expected prior and intended resource versions, each affected path, its
   before digest or explicit absence, its after digest, and the staged relative
   path. Bound file count and bytes using existing record limits.
3. Atomically install the intent only after staging is complete. Until this
   rename, canonical state is the exact prior valid record; orphan staging is
   safe to garbage-collect by verified operation ID.
4. Publish files in the current deterministic order. On any call, `doctor` or
   mutation preflight reads the intent before full-record validation and
   classifies every target as `before`, `after`, or `conflict` by digest.
5. `recover` refuses any conflict, otherwise replays every remaining exact
   after-image. It is safe to interrupt and rerun because an already-after file
   is skipped and no input is regenerated. Once the full record validates at
   the intended version, refresh README, remove the intent/stage, and return a
   structured committed receipt. If all targets are still `before`, recovery
   may instead remove the intent/stage and preserve the exact prior record.
6. Normal mutation on a goal with an intent returns a stable structured
   `publicationIncomplete` error with operation ID, intended version, phase,
   and the explicit `doctor`/`recover` action; it must not report an ordinary
   stale token or generic validation failure.

This protocol closes all checkpoint and criteria gaps with one mechanism. A
generation-directory plus atomic pointer would also work but is a larger
format/layout migration. Rollback images are unnecessary for the minimum
guarantee because every nonconflicting partial state can roll forward from
staged exact after-images; preserving the prior state is only needed before
the first canonical rename. Add file and parent-directory `fsync` only if the
stated guarantee includes machine/power failure rather than injected process
interruption, and document that boundary explicitly.

## Legacy migration guarantees

### What is supported and preserved

The only accepted legacy format is a flat unversioned directory containing
exact `README.md` plus optional root Markdown files. Directories, symlinks,
unsafe names, non-Markdown entries, invalid UTF-8, NUL bytes, excess count,
and oversize content are rejected (`migrate.go:672-715`;
`markdown.go:9-20`). No older structured API version is migrated.

Migration derives the goal ID from the source directory basename, extracts an
unambiguous H1 title unless overridden, and extracts list items only from an
acceptance-criteria section unless overridden (`migrate.go:335-341,376-390,
718-754`). It creates a new v1alpha1 Goal, revision-1 criteria, and a closed
`imported-unversioned` investigation attempt whose criterion verdicts are
`unverified` (`migrate.go:416-531`). The README bytes become `plan.md`; every
other Markdown file is copied byte-for-byte as evidence
(`migrate.go:565-598`). The source digest is deterministic length-framed
SHA-256 over sorted filename/byte pairs (`migrate.go:757-769`). The source is
re-read immediately before publication and is never written
(`migrate.go:615-625`). Tests confirm byte preservation, source
non-modification, and ordinary idempotence (`store_test.go:826-880`).

### Why the `legacy-migration` criterion is not yet met

1. **Source path is not retained.** `MigrationStatus` contains only
   `sourceFormat`, `sourceDigest`, and `migratedAt`
   (`projects/goal/api/v1alpha1/types.go:70-74`). Neither a normalized
   workspace-relative source reference nor equivalent provenance is stored.
2. **Field mapping is not explicit.** The record does not say whether title
   and criteria were extracted or supplied as overrides, what parser/mapping
   version was used, or which source spans produced which structured fields.
   Preserving the README as `plan.md` retains unmapped prose, but does not make
   the mapping inspectable.
3. **Idempotence is weakly bound.** Existing-target matching validates the
   target and compares migration metadata and high-level options/current
   criteria (`migrate.go:628-669`), but does not recompute the legacy digest
   from the imported attempt's `plan.md` and evidence. An internally valid
   record with copied migration fields can be misclassified as the completed
   import.
4. **Crash staging is unmanaged.** Deferred cleanup removes staging on an
   ordinary return, but SIGKILL can leave `.goal-migrate-*`; regular temporary
   cleanup does not recognize it (`migrate.go:544-553`; `integrity.go:21-79`).
   Goal listing skips dot directories, so residue is filesystem-visible but
   absent from the catalog (`store.go:305-316`).
5. **No-overwrite has a race.** The code checks `pathExists(target)` and then
   calls ordinary `os.Rename` (`migrate.go:612-623`). The test creates a target
   before that check (`migrate_test.go:515-537`), not between check and rename.
   Use an atomic no-replace directory publication primitive on supported
   platforms or narrow the guarantee.
6. **Invalid partial target is terminal.** A retry preserves and refuses an
   invalid existing target (`migrate.go:638-640`) because there is no receipt
   or staged after-image from which to finish it.

The migration tests cover prepublication hook failure/source preservation
(`migrate_test.go:452-485`), source change (`487-513`), pre-check target
appearance (`515-537`), matching and mismatched retries (`539-618`), unrelated
target preservation (`620-674`), and invalid Markdown
(`markdown_test.go:109-138`). They do not cover crash residue, lost
acknowledgement, post-check target races, provenance-to-snapshot binding, or
power-loss durability.

### Minimum migration repair

- Add normalized workspace-relative `sourceRef`, mapping format/version, and
  extraction-vs-override facts to migration provenance. Keep host-absolute
  paths out of portable state.
- On an existing candidate, require the exact canonical
  `imported-unversioned` attempt and recompute `sourceDigest` from its
  `plan.md` as `README.md` plus evidence basenames/bytes before declaring
  idempotent success.
- Under the already ordered source/target locks, recognize only structurally
  valid staging candidates for this goal and operation. With target absent,
  re-read the current source digest, then adopt a matching complete candidate
  or safely discard/quarantine stale candidates and rebuild. With a matching
  target, clean matching empty/stale roots.
- Replace check-plus-rename with true atomic no-replace publication. Add an
  after-publication failure hook to prove lost-acknowledgement idempotence.
- Preserve and refuse nonmatching or conflicting targets. If the contract
  covers only process interruption, say so; otherwise sync staged files,
  directories, and the destination parent before claiming power-loss
  durability.

## Acceptance conclusion

The current system provides good detection and partial committed-version
reporting, but not recovery. It cannot truthfully advertise resumability for
an invalid record. Phase 2A should implement the bounded publication intent,
`doctor`, `recover`, and boundary-complete injected-failure tests before any
goal catalog exposes rich continuation state. Legacy migration additionally
needs source-reference and mapping provenance plus stronger snapshot binding;
its existing non-destructive, full-directory staging design should be kept.

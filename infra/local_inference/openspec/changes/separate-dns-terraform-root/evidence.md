## Observed migration

On 2026-11-03, the user authorized extracting DNS, applying the independent root,
and preferring separate DNS roots in repo-infra. Source was committed before
state mutation. The provisioning branch's legacy DNS commands and Cloudflare
configuration were removed before migration. A bounded search found no local
inference writer in checked-in GitHub/Forgejo CI; external schedulers were not
inventoried. Saved-plan serial checks guarded concurrent state changes.

The legacy backend had exactly three managed records and no other managed
resources. The destination backend was empty. Private snapshots supplied the
exact import IDs; no record IDs or secret values are checked in here.

Both saved plans were reviewed before mutation. The legacy plan contained only
three `forget` actions under `module.dns`, with `destroy=false`. The destination
plan imported the same three IDs with no provider creates, updates, or deletes.
Removal completed first, and a state read confirmed zero legacy managed
resources before destination apply. Import completed with `3 imported,
0 added, 0 changed, 0 destroyed`.

Post-apply state comparison confirmed the original IDs, zone, names, record
types, values, TTLs, proxy flags, and priorities were unchanged. Both roots'
fresh plans exited 0 with no changes, without supplying the import map again.
Both authoritative Cloudflare nameservers returned the two Yandex NS values,
and a recursive query returned the unchanged T3 Code CNAME.

The new backend is `/tf_backend/dns`; `/tf_backend/tf_setup` now owns no DNS
records. No GPU host, child-zone provisioning, folder, certificate, or host
configuration was applied during this migration.

## Source validation

The independent root, migration and skill branch passed 38 offline tests before
apply, including formatting, ownership, both OpenSpec workspaces, skill
discovery and the Promptfoo configuration. Affected semantic lint passed.
These offline skill checks validate packaging/configuration, not model answers.
Final prepared candidates receive the mandatory repository checks before push.

Private snapshots, import maps, saved plans and raw operational output are
removed after verified ownership. Sanitized action summaries remain in ignored
`out/dns/`; delivery receipts own exact source and publication revisions.

`DNS-UNTRACKED-FORMAT`: initial repository formatting omitted the new untracked
Lua file. The prepared-candidate Lua check caught its layout; the configured
Lua formatter corrected it before final validation/publication. This changed
formatting only, not the reviewed backend or credential selection.

## Current project simplification

At the user's direction, all one-time state-transfer scaffolding was removed:
the migration directory, legacy backend call and import-only inputs are absent.
Only the active independent DNS root remains. Its plan reported no changes
against the already-adopted records; no new apply or state mutation was needed.

## 1. Acceptance fixtures before implementation

- [ ] 1.1 Add a hosted-service E2E fixture with fake OAuth/X endpoints, controllable time, two ordinary accounts, the owner, temporary SQLite state, and a request/charge journal; verify the fixture independently records reads, grants, mutations, and reservations without live credentials.
- [ ] 1.2 Write browser failure/acceptance cases for account isolation, CSRF, escaped text, SSR/htmx forms, cached reevaluation, budgets, stale plans, and exact subset confirmation; verify the cases expose the missing hosted behavior before production changes.
- [ ] 1.3 Write authorization cases for state replay, scope/account mismatch, token expiry/rotation, concurrent refresh, 24-hour inactivity, passive polling, logout, and restart; verify the fake provider observes expected exchanges and no credential appears in persisted evidence.
- [ ] 1.4 Write spending/recovery cases for concurrent last-credit admission, identity lookup, partial discovery, cancellation, unknown responses, UTC rollover, account recreation, and stale-backup restore; verify fixture assertions distinguish confirmed spend, reservations, and dispatch.
- [ ] 1.5 Write supported archive and hostile policy/archive fixtures, including edited IDs, repost mapping, missing seeds, bomb inputs, malformed wrappers, excessive output, and worker OOM; verify expected rejection and acceptance artifacts are specified before parsers or worker code.

## 2. Hosted configuration and durable account state

- [x] 2.0a Adopt the pinned chi router for the existing HTTP backend and expose it as `server serve`; preserve exact routes, capability/origin checks, strict JSON, bounded responses, cancellation and exact-action workflows through process/browser acceptance. This increment does not complete hosted configuration or account persistence.

- [ ] 2.1 Add the hosted command and strict deployment configuration, reusing existing SQLite/OAuth dependencies and the approved pinned htmx asset; verify missing configuration fails safely and the operator CLI/private serve retain their existing behavior.
- [ ] 2.2 Implement schema migrations and account-scoped persistence for sessions, policies, snapshots, plans, executions, and reservations; verify two-account isolation and migration/restart behavior through the E2E fixture.
- [ ] 2.3 Implement bounded retention, upload cleanup, ledger-preserving account deletion, disk admission, and consistent backup/restore procedures; verify expiry, full-disk rejection, and restoration without quota reset.

## 3. Temporary X authorization

- [ ] 3.1 Implement PKCE S256 login, one-time state/callback checks, identity binding, secure sessions, and CSRF/origin controls; verify spoofed identities, replayed callbacks, and forged actions fail before X mutation.
- [ ] 3.2 Implement separate cleanup authorization with verified read/write scopes and matching X identity, keeping access/refresh tokens and verifiers solely in memory; verify login alone cannot authorize deletion and inspect fixture database, cookies, logs, and backups for token absence.
- [ ] 3.3 Implement serialized on-demand refresh with rotated-token replacement, intentional-activity tracking, 24-hour idle expiry, disconnect/logout, and bounded revocation attempts; verify polling does not extend access and expiry or refresh failure stops dispatch.
- [ ] 3.4 Implement conservative admission for identity lookups and protected owner recovery when public funds are exhausted; verify a public request cannot bypass the pool by claiming owner identity and legitimate owner recovery still verifies the X account.

## 4. Durable spending admission

- [ ] 4.1 Implement integer-USD tariff configuration, conservative endpoint/response bounds, UTC monthly periods, stable owner exemption, and shared/account/cleanup/user ceilings; verify exact boundary values and reject invalid or negative budgets.
- [ ] 4.2 Implement atomic reserve/dispatch/settle records across all applicable allowances, separate from X request-window pacing; verify concurrent endpoint workers cannot overspend and accounting failure prevents dispatch.
- [ ] 4.3 Preserve unknown exposure through cancellation, reauthorization, restarts, account recreation, month rollover, and backup restoration; verify only proven-undispatched work releases reservations and explicit reattempts need new allowance.
- [ ] 4.4 Expose read estimates, selected-action cost, spent/reserved/remaining amounts, and user budgets, reserving the exact selected batch at confirmation; verify insufficient headroom is rejected before mutation and concurrent work cannot consume that batch's allowance.

## 5. Snapshot discovery and archive import

- [ ] 5.1 Persist explicitly selected archive source/key snapshots with bounded acquisition and no live read mode; verify acquisition failure preserves previous usable data and never issues X reads.
- [ ] 5.2 Persist immutable account-owned snapshots with provenance, version, fields, requested window, completeness, and retention; verify narrower-window evaluation preserves graph semantics and missing coverage never triggers implicit HTTP.
- [ ] 5.3 Implement supported X archive parsing inside the limited worker, with strict data-only wrappers, account/record validation, entry/path/decoded-size limits, and no media/DM ingestion; verify supported fixtures normalize correctly and hostile or ambiguous input leaves existing data unchanged.
- [ ] 5.4 Expose snapshot selection, age, coverage warnings, import, and dirty-state replacement; verify archive and cached evaluation work without an active X grant and issue zero fixture API requests.

## 6. Isolated policy execution and admission

Policy-worker increment (2026-11-09): implement the requested gRPC process boundary
first. The broader tasks below remain open until archive IPC, deployment isolation,
additional work counters, and hosted admission are complete.

- [x] 6.1a Add process-level gRPC acceptance cases before production code, using only synthetic snapshot data and no X credentials.
- [x] 6.1b Implement credential-free worker mode, private-socket RPC, and backend result validation.
- [x] 6.3a Enforce source/snapshot/record/output bounds and cancellable cached evaluation; verify exact semantics, busy rejection, deadlines, and recovery.
- [x] 6.5 Validate the worker increment plus existing CLI/browser behavior and inspect its repeatable plan artifact.
- [x] 6.6 Cross-check worker decisions against actions and forced indicators; verify forged post/repost contradictions fail and consistent responses remain accepted.
- [x] 6.7 Use typed protobuf for the current job contract and complete results; source-reference submission in 6.12 supersedes the initial snapshot stream.
- [x] 6.8 Expose the process as `worker serve` and verify it launches without credentials through the built command fixture.
- [x] 6.9 Remove unreleased RPC version labels, fields, checks, and compatibility stubs; verify the single current service contract through process-level acceptance.
- [x] 6.10 Replace blocking evaluation with detached submit/get/cancel/delete jobs, retaining backend-owned IDs and bounded transient state; verify connection independence, failure status, account isolation, limits, expiry, and restart through process acceptance artifacts.
- [x] 6.11 Remove collection-presence protobuf wrappers and normalize absent/empty attachment and entity collections in the shared Lua projection; verify CLI and worker policies retain populated metadata, missing counter keys, known zero values, and isolated nested tables.
- [x] 6.12 Replace live identity/timeline acquisition with standard X export ZIP input fetched directly by the worker from configured local directory providers behind an interface that supports future S3; use unary source/key job submissions, retain bounded archive/result validation, and verify offline planning, mutation-only authorization, hostile archives and asynchronous job behavior through process acceptance.

- [x] 6.13 Validate actual central-directory count, byte span and record lengths before standard-library allocation; retain one validated metadata view, honor cancellation during parsing, and verify forged footers and comment signatures with bounded-read and process acceptance artifacts.

- [x] 6.14 Stabilize exported reference order, retain bounded policy stage/line diagnostics without content disclosure, separate result consistency from archive membership, use typed tweet/entity decoding, and make terminal-result cleanup best-effort; verify failure cases before implementation and preserve parser limits and process artifacts.

- [x] 6.15 Generalize the worker to bounded Lua scripts over archive datasets with a complete JSON object result, without deletion-specific RPC messages or compatibility aliases. Keep source retrieval, account/time bindings, asynchronous lifecycle and no X access. Put cleanup semantics and exact-action validation in the cleanup application/helper; verify non-deletion plans, multipart likes/following, JSON conversion failures, failed dataset reads and cleanup behavior through process acceptance before delivery. Plans describe X API operations; archives remain read-only inputs. Synchronize the separate protocol PR.

- [x] 6.16 Replace generic JSON results with typed protobuf operation plans. Declare mutually exclusive operation payloads for post deletion, unrepost, unlike, unfollow, post creation/read and profile-description updates; keep cleanup explanations in an optional typed report. Convert Lua tables directly through the protobuf schema, reject missing/ambiguous/unknown operations and invalid fields at worker/client boundaries, remove JSON result and compatibility paths, and synchronize PR #282. Preserve generic archive scripting, exact cleanup validation, asynchronous source-reference jobs and existing execution limits through failure-first process acceptance.

- [x] 6.17 Generate Lua editor plan/payload declarations from the protobuf descriptor, including mutually exclusive operation variants; type standard archive account/tweet/like/following reads with literal-name overloads and a recursive JSON-value fallback. Annotate the shipped example and cleanup helper return. Write generation/example process acceptance first, verify source freshness and stable output through the owning generator, run existing lint/runtime checks, and disclose that LuaLS diagnostics are unavailable without an approved new tool dependency.

- [x] 6.18 Document the application in a concise left-to-right Mermaid flowchart: Actor points to Backend through HTTP and Upload archive; Backend submits worker jobs, reads results, writes the archive and applies plans through X (x.com); Worker reads the archive. Group Backend, Worker and Archive under x_profile_cleaner and omit implementation components. Render SVG and dark WebP through normal repository-theme targets and inspect the output. Identify upload/write as planned functionality and keep changes local until publication is requested.

- [x] 6.19 Replace annotation-only dataset overloads with named, validated account/tweet/like/following accessors and explicit raw reads. Generate editor records from the runtime Go schemas, return unwrapped records, reject malformed data even across pcall, and verify multipart data, exact IDs, fresh projections and failure recovery through process fixtures written first. Remove the ambiguous legacy read API.
- [x] 6.20 Replace policy evaluation functions with an evaluator struct and explicit options, preserving worker-owned deadlines, stateless evaluations, cleanup provenance and existing callers. Verify configurable evaluation timeout and parent cancellation before implementation; retain the existing CLI/browser/process acceptance.
- [x] 6.21 Move configurable generic evaluation into script.Evaluator, rename policy to cleanup without aliases, and remove test-only archive.Read and Snapshot.AccountID. Produce protobuf plans directly in the planner and retain them throughout worker/backend/CLI/browser/apply with a trusted account/time document envelope. Remove the old JSON version/actions schema and bidirectional plan converters. Write saved-plan acceptance and failure cases first; retain exact cleanup provenance, subset/resume, executor admission, source limits and credential boundaries.
- [x] 6.22 Define canonical account, following, like, tweet and nested archive records in protobuf. Decode X JSON at the importer boundary and expose only protobuf records to the generic evaluator. Derive archive Lua bindings and editor declarations from descriptors; remove exportedTweet aliases, DatasetTypes and Go-reflection projection/generation. Use normalized field names, typed timestamps/enums/counters, exact IDs and safe large integer projection without compatibility aliases. Write normalization, metadata, optional values, mutation isolation and malformed-input acceptance first; preserve cleanup and raw dataset behavior. Keep ZIP handling in the worker.
- [x] 6.23 Remove the remaining Go JSON archive record structs, custom protobuf json_name aliases and model.Item-to-protobuf round trip. Import standard X source fields directly into canonical protobuf records; accept one source spelling per field, with no canonical/case-fold fallback. Move cleanup's graph projection and snapshot into cleanup, leaving the generic archive API protobuf-only. Write source-name rejection cases first and preserve supported standard-export data, exact counters, metadata, provenance, parser bounds and credential-free worker behavior.
- [x] 6.24 Give saved JSON documents their own serialization-aware bound, check the final formatted output before writing, and retain the decoded protobuf limit. Normalize optional null source counters to absent values while preserving zero and rejecting invalid numbers. First reproduce both failures through real CLI/worker fixtures, then verify large saved reports reload and execute only their selected operations, and both typed accessors and cleanup accept null counters.
- [x] 6.25 Replace returned Lua plan tables with per-evaluation schema-derived plan methods. Collect protobuf operations directly, validate arguments and aggregate bounds before appending, retain caught API failures, reject script return values and publish only complete successful evaluations. Derive runtime argument shapes and generated editor declarations from protobuf; support scalar targets/text, optional flags and structured profile updates. Append native cleanup results directly without a protobuf-to-Lua round trip. Remove the old return decoder and editor return types. Write process acceptance first for all operations, ordered results, options isolation, failures, cancellation, limits, cleanup and recovery.
- [x] 6.26 Prevent a policy suffix from combining with the wrapper's closing delimiter. First reproduce valid trailing-comment policies with several delimiter levels through the real CLI; preserve source content, exact selected actions and zero X requests during planning.
- [x] 6.27 Add representative archive-script examples for filtered deletion, repost removal, relationship selection, explicit reads, and profile/post updates. Test every maintained example through a dedicated Bazel target whose JSON configuration owns archive fixtures, fixed evaluation time, entrypoint and exact expected protobuf plans. Cover populated and empty exports, missing metadata, large IDs, selection exclusions and operation order; fail when any example lacks a configured case. Write the configured acceptance cases before the new scripts.

Decision for 6.25: proceed with one injected plan collector per evaluation. The
worker has one output, so an explicit returned object adds construction and
conversion without supplying execution authority or useful composition. Method
names come from the Operation oneof. A single required scalar (including target
IDs) is positional; remaining payload fields are options, while presence-sensitive
profile updates use a table. The same descriptor rule owns runtime and editor
signatures. Generic calls and native cleanup helpers append to the same bounded
protobuf result. This replaces the unreleased return contract without fallback.

PR #280 review follow-up (2026-11-13): failure-first process fixtures reproduced
unloadable saved reports and all six rejected null counters. The corrected CLI
round-trips 40,000-decision documents of 17,269,348 bytes (ASCII) and 30,069,348
bytes (escaped text), applying only the selected post deletion and issuing no X
reads. Worker acceptance preserves absent versus zero counters and rejects
negative, fractional, boolean and overflowing counters. The JSON document limit
is 128 MiB, independently checked before final file replacement and during bounded
loading; decoded protobuf bounds remain enforced. Evidence is under ignored
`out/profile-cleaner-worker/publish/`.

Local acceptance for 6.19–6.20 (2026-11-13): 32 project, artifact-boundary and root
quality tests pass; semantic lint passes for 12 affected package/consumer targets.
`typed-datasets-plan.json` verifies exact multipart like/following IDs and a typed
four-operation plan. `evaluator-options.json` verifies option timeout, parent
deadline/cancellation, invalid configuration and isolated interpreter state.
Fixtures preceded implementation; the initial process fixture failed on the missing
named accessors. Logs and inspected artifacts are under ignored
`out/profile-cleaner-worker/typed-datasets/`. Changes remain local and unpublished.

Session ergonomics: the bounded acceptance log and two small output artifacts
support resumption without reading the full process transcript. One oversized
combined log read exceeded context; subsequent reads use bounded tails and explicit
artifact paths. No shared workflow change is needed for that corrected tool use.

Evidence: `//projects/x_profile_cleaner:test` and the command's
`artifact_boundary_test` pass. `worker-plan.json` records the inspected exact
unrepost of source `900` and forced deletion of reply `2`, preserving edited root
`101`. Acceptance rejects foreign-account inputs and forged response bindings,
targets, malformed typed messages, and output size. A separate cached-window evaluation selects only
recent post `3`. Deadlines, active cancellation, and concurrent admission are
exercised against the built process. No live credentials, X access, infrastructure
changes, or hosted-service completion are implied by this increment.

PR #280 review follow-up: action/decision consistency is now checked in both
directions. The gRPC forged-worker fixture rejects unchanged/preserved actions,
missing actions or explanations, and either forced-flag mismatch for posts and
reposts, while accepting consistent ordinary and forced responses. Existing
cached-window and exact-plan acceptance still passes. The pre-fix fixture exposed
the reported contradictions before the validator change.

- [ ] 6.1 Add fixed worker executable/IPC contracts for archive parsing and policy evaluation, excluding credentials and validating returned data in the backend; verify malformed or oversized worker output cannot publish a plan.
- [ ] 6.2 Package systemd worker units with separate cgroup OOM behavior, memory/CPU/task/runtime restrictions, a sanitized environment, and no network/general host-file access; verify real resource exhaustion in a disposable VM leaves backend login/status responsive.
- [ ] 6.3 Implement policy source/rule/identifier/callback/graph/output bounds and cancellable evaluation over one cached snapshot; verify infinite loops, allocation growth, and excessive traversal fail without partial executable plans.
- [ ] 6.4 Implement account/global active and waiting limits plus OAuth/action rate controls; verify one account's repeated cached evaluations cannot create unbounded jobs or monopolize all admission.

## 7. SSR review and durable exact apply

- [ ] 7.1 Replace hosted JSON-driven UI state with escaped server templates and bounded htmx fragments while preserving the private interface; verify light/dark/narrow browser screenshots and keyboard-accessible review/forms.
- [ ] 7.2 Implement per-evaluation immutable plan publication, strict subset selection, stale-plan invalidation, and atomic apply consumption; verify duplicate requests, fabricated indexes, editor changes, and another account's identifiers cannot broaden execution.
- [ ] 7.3 Integrate ephemeral OAuth, monetary reservation, and existing endpoint pacing with durable mutation intent/outcomes; verify exact post/repost targets, dependency order, cancellation, absent targets, and no automatic retry.
- [ ] 7.4 Invalidate affected snapshots/unconsumed plans after dispatch and implement explicit resume from retained selection/progress; verify restart dispatches nothing, unknown outcomes remain visible, and confirmed completions never replay.
- [ ] 7.5 Implement bounded account-scoped progress/log polling and authorization/budget/capacity error fragments; verify sensitive values stay absent and polling never renews inactivity authority.

## 8. Deployment packaging and integrated acceptance

- [ ] 8.1 Establish the owning infrastructure workspace/change and project boundary documentation for the selected deployment, reusing existing Terraform/Ansible/Vault conventions; verify target ownership and secret references without provisioning or configuring live systems.
- [ ] 8.2 Package the single-VM backend, Traefik TLS/origin routing, SQLite state directory, worker services, secret environment injection, no-swap/core-dump controls, logs, health checks, and backups; verify packaged configuration and offline checks, mirroring any intended bot-host change into host_bot if applicable.
- [ ] 8.3 Run fixture-backed full browser/service acceptance in a disposable VM with real systemd limits and restart/restore exercises; inspect a repeatable acceptance report, screenshots, exact HTTP journal, and sanitized accounting export.
- [ ] 8.4 Run the existing CLI/browser suites, hosted acceptance, strict OpenSpec validation, repository quality, and semantic lint against the candidate; inspect representative artifacts and resolve in-scope findings before implementation delivery.
- [ ] 8.5 Complete implementation review and PR delivery with source/evidence references; verify publication receipts and report live OAuth enrollment, billing observations, and rollout as separate explicitly authorized operations.

Archive-source verification (2026-11-12): process acceptance verifies unary
source/key submissions, multiple configured local providers, multipart exports,
wrong-account and unsafe-path rejection, malformed wrappers, CRC and entry bounds,
non-regular sources, offline CLI/browser planning and mutation-only API use.
`archive-cli.json` and `source-archives.json` record repeatable observations.
No S3 provider/dependency, hosted persistence or hard OS resource isolation is
included. Backend structural validation relies on the source-owning worker for
archive membership, as documented in design.md.

Directory-allocation review r4550014478: the pre-fix read observer consumed
1,366,603 bytes for 25,000 real headers behind a two-entry footer and continued after
cancellation. The corrected preflight rejects count/size/length contradictions
before standard-library allocation, captures only bounded validated metadata,
uses a canonical footer to avoid comment-signature reinterpretation, and blocks
ZIP64 locator redirection. `archive-directory-bounds.json` records bounded reads
and cancellation; source-process acceptance verifies rejection and later recovery.

Review simplification acceptance: process fixtures prove 64 stable dual-reference
evaluations, content-free syntax/callback locations, successful plan publication
when DeleteEvaluation is unavailable, and unchanged forged-response rejection.
Typed decoding retains strict counters, required text, entity offsets and all
archive limits. `policy-diagnostics.json` and `cleanup-result.json` record the
repeatable outcomes. Hard OS limits remain deployment work in task 6.2.

Generic-script acceptance: archive reads, arbitrary Lua computation and source/job
bindings remain from task 6.15. Task 6.16 supersedes its JSON result format with
`typed-plan.json`, demonstrating typed unlike/unfollow/create/read/update payloads,
exact string IDs, optional cleanup reports and explicit empty-description updates.
Failure cases cover unknown/ambiguous operations, targets/types/fields, missing
payloads, invalid protobuf metadata, size/complexity bounds and caught archive
failures. The original exact cleanup, lifecycle and archive-bound cases remain.

Typed-plan verification: malformed operation plans accepted by the former JSON
boundary are rejected by direct Lua/protobuf conversion. Process cases retain
exact IDs and explicit empty profile updates, reject unknown/multiple/missing
payloads, wrong targets/types/fields and excess operation counts, and preserve
existing cleanup/lifecycle/source behavior. A failure-first report-mutation case
also proves final cleanup provenance is checked after Lua returns, using the same
cached normalized snapshot. The independent protocol PR contains identical schema
sources. No new X executor, credential, external dependency or OS deployment is
introduced by this increment.

Lua typing acceptance: the owning generation command produces stable declarations
matching the checked-in api/types.lua projection. Generation failures preserve
existing output; generated Lua parses successfully. The annotated operation example
executes through the built worker and yields the expected unlike/unfollow plan.
`lua-api.json` records generation checks and explicitly reports that LuaLS diagnostics
were not run; `lua-example-plan.json` records the example result. No new external
checker dependency or runtime API behavior is introduced.

Architecture documentation: `docs/architecture.mmd` is the maintained source for
the README's current component map. The `architecture`, `architecture_webp` and
`docs` targets render and package repository-themed SVGs and a dark WebP directly.
The five-node overview labels all relationships, including the backend's archive
write path and the worker's ZIP reads. The upload/write path is design intent;
this diagram change does not implement its endpoints or complete hosted,
durable-storage or deployment tasks.

Local acceptance for 6.21 (2026-11-13): the planner and application now keep one
protobuf operation payload inside the account/time Document; the old Plan,
Action, Decision and Reason models and both conversion functions are removed.
Generic options live in script.Evaluator and cleanup.New creates per-call library
state. No policy package/timeout alias or archive.Read wrapper remains. Saved-plan
acceptance was written first and exposed the legacy output. All 32 project/root
quality/credential-boundary checks and semantic lint for 16 affected targets pass.
The browser suite verifies exact subsets, stale/fabricated selections and resume
without replay; inspected protobuf-plan.json and the dark review screenshot match
the new contract. Logs and representative outputs are retained under ignored
out/profile-cleaner-worker/plan-migration/. No publication or live X access occurred.

Migration review: retain trusted account/time outside the script payload; reject
unknown/ambiguous protobuf fields and obsolete files; check every executor payload
before dispatch. Selection and resume filter deeply cloned operations, preserving
the original review object. The protobuf schema itself is unchanged. Scoped Gazelle
needed an explicit public-wrapper resolution and a preserved proto embed; the two
custom fixture targets remain manually maintained. A broad initial inventory and
text replacement caused avoidable reads and build repair; keep subsequent discovery
and edits bounded to declarations and owning targets. No shared skill change is
needed for these corrected uses of existing procedures.

Local acceptance for 6.22 (2026-11-13): api/archivepb/archive.proto owns the
canonical records, including normalized timestamps, enums, optional counters and
metadata. Archive accessors return generated protobuf messages. The shared
descriptor projection and declaration generator replace both Go-reflection paths;
only private JSON importer structs and the cleanup graph's domain projection
remain. Process cases were updated first and failed on the missing normalized API.
They now verify known zero versus missing counters, exact large values/IDs,
multipart records, fresh nested tables and fatal validation errors across pcall.
All 32 project/root-quality/credential-boundary tests and semantic lint for 19
targets pass. Generated declarations and typed-datasets-plan.json were inspected;
evidence is under ignored out/profile-cleaner-worker/archive-protobuf/. No new
dependency, RPC record streaming, backend extraction or publication was introduced.

Session review: scoped Gazelle reused the existing wrapper-resolution convention
without build repair; generation succeeded before acceptance and the first worker
run passed. Reuse those results for unchanged inputs. The remaining final check
covers only completed evidence and current-behavior documentation; no shared skill
change is warranted.

Local acceptance for 6.23 (2026-11-13): the alternate account/like/following keys
and case-folded tweet key were accepted before the change; all four process
fixtures failed as expected. They now fail closed without a plan. All Go JSON
archive record structs and protobuf json_name attributes are removed. Standard X
source fields construct protobuf messages directly; the importer no longer
round-trips through model.Item or exposes a cleanup snapshot. Cleanup's graph
projection is private to cleanup, and obsolete model JSON tags and the unused
source-test model dependency are removed. All 32 project/root-quality/boundary
checks and semantic lint for 19 targets pass. Failure evidence, logs and the exact
typed operation artifact are retained under ignored
out/profile-cleaner-worker/remove-compatibility/. Changes remain unpublished.

Session review: failure-first cases identified four real alternate decoding paths,
so this removal changes accepted input rather than merely renaming types. The first
full behavioral run passed. Keep generic archive parsing and cleanup graph concerns
with their current owners; no shared skill update is needed.

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

- [ ] 5.1 Implement explicit budget-aware API refresh with bounded page sizes, response bytes, records, and pagination; verify partial failure preserves the previous snapshot while dispatched reads remain charged/reserved.
- [ ] 5.2 Persist immutable account-owned snapshots with provenance, version, fields, requested window, completeness, and retention; verify narrower-window evaluation preserves graph semantics and missing coverage never triggers implicit HTTP.
- [ ] 5.3 Implement supported X archive parsing inside the limited worker, with strict data-only wrappers, account/record validation, entry/path/decoded-size limits, and no media/DM ingestion; verify supported fixtures normalize correctly and hostile or ambiguous input leaves existing data unchanged.
- [ ] 5.4 Expose snapshot selection, age, coverage warnings, import, refresh, and dirty-state replacement; verify archive and cached evaluation work without an active X grant and issue zero fixture API requests.

## 6. Isolated policy execution and admission

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

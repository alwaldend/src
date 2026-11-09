## Context

See [proposal.md](proposal.md) for motivation and capability scope. The existing
loopback server owns one application, credential, plan, execution, and log stream.
Its Lua adapter has restricted libraries and cancellation, but runs inside the web
process. API request-window pacing exists; durable billing limits do not.

The baseline safety contract remains in
[profile-cleanup](../add-x-profile-cleaner/specs/profile-cleanup/spec.md). Hosted
behavior is additive: the operator CLI and private loopback workflow remain usable.

## Goals / Non-Goals

**Goals:** Reuse the planner's plain-data boundary, isolate account state and
untrusted work, separate persistent data from temporary authority, and make every
paid request attributable before dispatch.

**Non-Goals:** A distributed queue, Redis, multiple application replicas, user-supplied
X apps, browser-held X credentials, or a rewrite of policy semantics. Import accepts
supported archive layouts without claiming complete history or graph coverage.

## Decisions

### One backend and resource-limited workers on one VM

Run Traefik and a Go backend as systemd services. Bind the backend to loopback and
expose only the configured HTTPS origin. Serve escaped Go templates and locally
packaged htmx; normal forms and fragment responses share handlers and authorization.
Use bounded htmx polling for job state and logs. Polling does not keep X authority
alive. Avoid a second frontend build application or client-side credential store.

Add a hosted command/configuration boundary rather than relaxing the private
`serve` command's capability and loopback checks. Keep the CLI's OAuth 1.0a adapter
separate from hosted OAuth 2.0 user-context clients. The hosted command must not
inherit the Article uploader's personal user credential as a fallback.

Run policy parsing/evaluation and archive parsing in separately limited systemd
worker units. A worker receives only the selected account's bounded source data,
policy, fixed evaluation time, and operation identifier through bounded local IPC.
It has no OAuth environment, network access, arbitrary host filesystem access, or
permission to manage units. Backend-owned admission starts a fixed executable/unit
template; user input cannot choose unit properties or executable arguments.

Require per-worker cgroup memory limits and separate OOM handling so a worker's
failure does not terminate the backend. Disable swap and core dumps for services
holding credentials. Worker environments explicitly exclude deployment secrets.
In-process timeouts alone were rejected because a Lua allocation or parser can
exhaust backend memory before cooperative cancellation.

### OAuth 2.0 login and temporary cleanup authorization

Use authorization code flow with PKCE S256, unpredictable one-time state, exact
redirect URIs, bounded callback lifetime, and backend confidential-client exchange.
The deployment supplies `X_OAUTH_CLIENT_ID` and `X_OAUTH_CLIENT_SECRET` from Vault;
public origin and owner X user ID are configuration, not credentials.

Login requests `tweet.read`, `users.read`, and `offline.access`. Resolve identity
once for a newly obtained grant with `/2/users/me`; store the stable X user ID and
an opaque application session. Read/write cleanup uses a separate authorization
request adding `tweet.write`, verifies the actual returned scopes and the same X
user ID, and rotates the local session identifier. A write-capable grant is never
silently attached to another logged-in account. Separate flows need not imply two
developer apps; restrict backend endpoints by the verified grant's permissions.

Keep access tokens, refresh tokens, PKCE verifiers, and grant state in process
memory only. SQLite sessions contain no X tokens, encrypted token copies, or data
that can recreate them. Use Secure, HttpOnly, host-only SameSite cookies, explicit
CSRF protection for actions, origin checks, and a fixed allowlist of X endpoints.
Exclude callback query strings, headers, token responses, and credential values
from application, proxy, diagnostic, and crash logs.

Refresh on demand under one lock per grant. If X rotates the refresh token, replace
the pair atomically in memory; do not replay an uncertain refresh exchange blindly.
Refresh failure stops X operations and requests reauthorization. Automatic refresh,
polling, and job progress never extend the 24-hour inactivity deadline. Intentional
authenticated navigation and user-submitted actions do; reject already expired
grants before recording activity so a late visit cannot resurrect them.

Disconnect, application logout/session expiry, inactivity expiry, and shutdown stop
new dispatch and cancel active X work. Attempt revocation where possible, then
discard local credentials regardless of revocation success. A restart cannot
restore a grant. Persisted sessions can still access their account's cached data;
new X operations require authorization again. X-side consent can outlive local
token deletion, and in-flight requests can complete after cancellation.

OAuth 1.0a was considered to obtain an ID in its token response, then rejected for
the hosted service: avoiding an identity lookup does not justify changing the
agreed OAuth 2.0 flow. Its existing use in the operator CLI is unaffected.

### SQLite owns durable state, never X authority

Use the existing pinned SQLite driver with foreign keys, WAL, a busy timeout, short
transactions, and one coordinated writer. Persist accounts, hashed application
session identifiers, policies, versioned snapshots, immutable plans, submitted
selections, per-action outcomes, and budget reservations. Every query and job lookup
is scoped to the authenticated internal account ID. Handles are display metadata;
only configured stable X identity controls owner exemption.

Store intent before dispatch and outcomes after confirmation. A crash between
those records leaves an unknown outcome. On startup, mark formerly running jobs
interrupted; do not restart them. Resume requires the original account, original
submitted selection, fresh permission where needed, explicit confirmation of
unfinished/unknown actions, and allowance for each new request. Confirmed actions
are never replayed. Preserve endpoint order and same-target dependencies. A new
policy or plan cannot replace an execution's retained selection.

### Separate discovery from evaluation

Refresh from X and Import X archive produce immutable snapshots. Evaluate policy
uses exactly one authorized cached snapshot, performs no HTTP identity check or
lazy network read, and publishes a new plan ID. Preserve Lua's lazy query behavior
against the local snapshot. Freeze a per-evaluation time so relative dates and ages
are consistent without pretending cached metrics or text are fresh.

Record owner, source, creation/import times, available fields, requested discovery
window, and completion/coverage limitations. Cache keys include those properties
and normalization version, not a username alone. Narrower windows filter before
building the policy graph. Unbounded policies use the available snapshot with its
coverage warning; they do not establish complete history. Missing required
coverage or fields returns an explicit refresh/import choice. Never silently merge
archives and API data or spend credits during an evaluation.

Budget-aware pagination can reduce page size within X's supported range. Before a
page, reserve its maximum possible returned resources and requested expansions.
An incomplete refresh never replaces a usable snapshot or produces an executable
partial plan. Reads already dispatched still count even if collection or a later
policy fails. Discovery estimates include user-selected headroom for deletion;
when complete discovery cannot fit, suggest a narrower window or archive import.

Read supported account/post/repost archive data without executing the bundled site
or scripts. Treat script-wrapped JSON as data with a strict parser; use the existing
dependency set where possible. Parse bounded selected entries rather than extract
the whole archive. Reject path traversal, links, nested archives, unsupported or
conflicting records, missing identity, and invalid repost mappings. Verify archive
identity against the logged-in account, but never treat upload metadata as proof of
ownership. Store only supported own records, not DMs, contacts, media binaries, or
unneeded archive files. Missing preservation seeds fail as in API discovery.

The backend validates the worker's returned snapshot/plan before publication.
After a dispatched mutation, mark affected snapshots and other unconsumed plans
dirty, including when the result is unknown. Do not alter the consumed execution's
original plan or disable explicit recovery. A dirty snapshot needs an explicit
refresh or replacement import before another evaluation; display its status.

### Money limits are independent of request windows

The monetary contract is [cleanup-spending](specs/cleanup-spending/spec.md).
Represent USD as integer microdollars, never binary floating point. Each explicit
cleanup owns one budget account: refresh, repeated evaluation, apply, cancellation,
and resume retain its accumulated spend and outstanding reservations. A user may
set a lower run and/or personal monthly ceiling. A new cleanup does not reset
account or shared monthly usage. Raising a budget is explicit and remains bounded
by service caps. Budget reductions cannot go below already committed exposure.

In one SQLite transaction, reserve the upper bound against the run, personal,
account, and shared allowances before dispatch. Wait for request-window admission
without holding a database transaction; recheck authorization, cancellation,
period, and reservation immediately before dispatch. Each concurrent endpoint
request needs its own reservation. At apply confirmation, reserve the entire
selected batch and allocate child request reservations from that hold without
double counting, so another job cannot spend its approved allowance. On explicit
resume, reserve unfinished actions again while retaining unresolved earlier
exposure. If accounting fails, fail closed.

Settle only when evidence establishes the charge. Keep ambiguous outcomes reserved
at their upper bound, including identity reads and uncertain mutation responses.
Crash recovery never refunds an intent whose dispatch outcome is unknown. Release
an undispatched reservation only when the durable state proves no dispatch could
occur. Do not assume missing-target responses, errors, retries, or same-day read
deduplication are free. Conservative accounting can leave usable provider credit
unavailable until reconciled; it must not oversubscribe the public allocation.

Monthly periods are UTC calendar months. Attribute a dispatched request to its
dispatch period; retain unresolved reservations and conservatively constrain new
admission by unresolved exposure across rollover. Undispatched work crossing a
boundary must be admitted again against the current period. Account deletion and
recreation cannot clear usage: retain a minimal identity-keyed ledger through the
applicable period and unresolved-exposure lifetime, separate from user content.

The owner bypasses service monetary caps, including the exhausted public pool, but
still honors any user-selected budgets, valid permission, and request-window
and resource admission. Do not set an X-wide $5 provider ceiling that also blocks
the owner. Attribute identity requests before their X identity is known to the
public pool by default. Initial identity acquisition has no verified account or
cleanup yet: account/run caps start after identity is established, while that
anonymous overhead remains charged to the public pool under strict rate admission.
Show this distinction in accounting rather than promise per-account control over
an identity that is not yet known. Use an operator-only login mechanism
for recovering owner access when that pool is exhausted. Gate that mechanism with
a deployment-managed recovery secret, separate rate limits, and post-exchange owner
ID verification; never expose a public "owner" switch or bypass by supplied ID.
Its incidental failed identity checks have their own small deployment-configured
recovery allowance outside the public pool; a verified owner then has the agreed
uncapped service allowance.

Price configuration carries its observation date, endpoint/category mapping,
conservative rounding, and maximum billable response size. As researched on
2026-11-09, public Post reads are $0.005/resource, owner-app timeline reads
$0.001/resource, and ContentManage is $0.005/request. The deletion category mapping
comes from X's Playground, not a live invoice. `/2/users/me` free treatment remains
unverified: reserve the published User category's $0.01 until confirmed. These are
dated implementation inputs, not permanent tariffs. Keep the provider ceiling and
application reservation semantics explicit rather than promise exact invoices.

### Proposed configurable operating defaults

These values complete unspecified implementation details for review; they are not
host configuration changes. Display applicable limits in the UI and reject startup
configuration whose aggregate worker/storage bounds cannot fit the VM.

| Limit                       | Initial value                                    |
| --------------------------- | ------------------------------------------------ |
| Policy source               | 64 KiB UTF-8, bounded before parsing             |
| API non-upload request body | 256 KiB                                          |
| Archive upload              | 64 MiB compressed                                |
| Archive decoded entries     | 256 MiB total; 10,000 entries; no nested archive |
| Snapshot                    | 50,000 records and 32 MiB normalized data        |
| Cached snapshots            | 64 MiB/account; 512 MiB globally; seven days     |
| Policy declarations         | 256 rules; 10,000 explicit IDs total             |
| Callback/graph work         | 1,000,000 callbacks and 1,000,000 edge visits    |
| Planning worker output      | 16 MiB, including decisions and actions          |
| Worker unit                 | 256 MiB memory; 100% CPU quota; 32 tasks         |
| Worker lifetime             | 90 seconds, including parsing and IPC            |
| Admission                   | One active operation/account; two workers/VM     |
| Waiting queue               | 20 jobs globally; one queued job/account         |
| Unconsumed plan retention   | 24 hours; cannot outlive its snapshot            |
| Stopped execution retention | Seven days from last explicit apply/resume       |
| Upload staging              | Removed after import/failure; one-hour sweeper   |
| Application sessions        | Seven-day absolute lifetime                      |
| Browser progress logs       | Latest 1,000 entries; 4 KiB/entry                |
| Durable operational logs    | Seven days; bounded rotation; no post content    |
| Rate control                | Five OAuth starts/IP/minute and 20/account/hour  |
| Total runtime data disk     | 2 GiB including database, WAL, uploads, and logs |
| Accounting disk headroom    | 128 MiB reserved; stop data admission below it   |
| Consistent backups          | Daily; seven-day retention in bounded storage    |

Limits apply to the owner too. Enforce total disk use across database, WAL, retained
execution records, uploads, and logs, reserving free space for ledger commits.
Evict only eligible unreferenced cache records; preserve a running or retained
execution's required data until completion/expiry. If protected data fills the
limit, reject new work rather than silently delete progress. Account deletion clears
content and cancels work; explain the minimal anti-reset billing ledger retention
and backup expiry. Keep credential storage absent from backups.

### Deployment and evidence ownership

Future service deployment belongs in an owning `infra/x_profile_cleaner` component
with its own OpenSpec change, README, provisioning stage if needed, packaged Ansible,
DNS root, and Vault references. Reuse existing reusable roles. Mirror any intentional
change to the current bot host in `users/simeonwarren/host_bot`; do not provision the
service there implicitly. Application source remains owned by this project.

Offline acceptance uses local OAuth/X fixtures, a disposable SQLite database, and
isolated worker units. Record exact read/mutation requests, budget reservations,
cross-account denials, restart outcomes, and responsive SSR screenshots in a
repeatable acceptance artifact. Resource-isolation acceptance must exercise real
systemd/cgroup limits in a disposable VM, rather than mock a successful unit launch.
Write failure cases before production changes. No live credentials or deletion
belong in fixtures. Live OAuth enrollment, billing verification, and rollout require
the separately authorized environment and scope.

### Evidence and decision verdict

Proceed with the single-VM design. SQLite and systemd satisfy the agreed durability
and isolation needs without another service. Retain OAuth 2.0 plus memory-only
refresh instead of OAuth 1.0a or durable tokens. Retain a global public ceiling
instead of an allowlist or user-supplied apps. Accept that public funds can be
exhausted by a small set of users, and enforce separate CPU/storage admission.

External references, observed 2026-11-09:

- [X OAuth 2.0 and offline access](https://docs.x.com/fundamentals/authentication/oauth-2-0/authorization-code)
- [X identity endpoint](https://docs.x.com/x-api/users/get-my-user)
- [X billing and owned reads](https://docs.x.com/x-api/getting-started/pricing)
- [Public category rate card](https://console.x.com/api/credits/pricing)
- [X Playground endpoint mapping](https://github.com/xdevplatform/playground/blob/main/internal/playground/credits.go)

## Risks / Trade-offs

- A VM/backend failure drops temporary grants: retain data and progress, show
  reauthorization, and require explicit Resume.
- X consent and in-flight requests can outlive local cancellation: avoid promising
  immediate remote revocation or exactly-once effects across network failures.
- Cached/archive data can be old or incomplete: show provenance and scope, fail on
  missing explicit seeds, and never describe successful planning as full coverage.
- A tiny public allowance can be spent collecting an unusable partial refresh:
  bound reads, preserve the old snapshot, expose headroom, and offer archive import.
- Price changes and uncertain responses can make estimates differ from invoices:
  use conservative dated tariffs, durable reservations, and operator reconciliation.
- Account isolation and archive intake expand the attack surface: enforce every
  read/write boundary on the server, escape text, and test hostile fixtures.
- Single-VM storage failure affects all users: protect the ledger, back up SQLite
  consistently with bounded retention, and exercise restoration before rollout.

## Migration Plan

1. Keep the existing CLI and private `serve` command unchanged while adding hosted
   configuration, schema migrations, and offline acceptance fixtures.
2. Implement account boundaries, ephemeral grants, and durable request admission
   before enabling discovery or mutation routes.
3. Add snapshots/import, isolated evaluation, SSR review, and durable exact apply.
4. Package deployment and run acceptance in a disposable single VM with fixtures.
5. For an explicitly authorized rollout, configure the selected hostname, OAuth
   callbacks, owner identity, Vault credentials, and backups through the owner.
6. On rollback, disable new admissions, cancel outstanding work, discard grants,
   preserve the ledger/progress, and restore only a schema-compatible application.
   Restoring an older backup must not reset spending: reconcile dispatch records or
   conservatively exhaust public allowances for the affected period before reopening.

## Open Questions

- Deployment input: VM provider/location, hostname, sizing above the configured
  resource floor, and backup destination.
- Live billing observation: confirm `/2/users/me` treatment and deletion category
  mapping for the deployed developer account before relaxing conservative tariffs.

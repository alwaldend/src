## Context

See [proposal.md](proposal.md) for motivation and capability scope. The existing
loopback server owns one application, credential, plan, execution, and log stream.
Its restricted Lua adapter now runs in a separate archive-owning worker. API request-window pacing exists; durable billing limits do not.

The README's current architecture diagram is maintained in `docs/architecture.mmd`
and rendered through the project's documentation targets. It describes implemented
actors, application processes, archive storage and external services; the hosted
design below includes further work tracked by the remaining tasks.

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

Use chi for explicit HTTP routes and standard net/http middleware. The current
private interface launches through `server serve`; its local-authority and
capability checks precede routing, including unmatched paths and methods.
Keep strict bounded JSON decoding and application job ownership outside the
router. The user selected chi; the pinned module has no additional dependencies.

Run Traefik and a Go backend as systemd services. Bind the backend to loopback and
expose only the configured HTTPS origin. Serve escaped Go templates and locally
packaged htmx; normal forms and fragment responses share handlers and authorization.
Use bounded htmx polling for job state and logs. Polling does not keep X authority
alive. Avoid a second frontend build application or client-side credential store.

Add a hosted command/configuration boundary rather than relaxing the private
`server serve` command's capability and loopback checks. Keep the CLI's OAuth 1.0a adapter
separate from hosted OAuth 2.0 user-context clients. The hosted command must not
inherit the Article uploader's personal user credential as a fallback.

Run policy parsing/evaluation and archive parsing in separately limited systemd
worker units. A worker receives a configured source reference, expected account, policy, fixed
evaluation time, and operation identifier through bounded local IPC. It fetches
the archive directly through an operator-configured provider. It has no OAuth
environment, general network access, arbitrary host filesystem access, or
permission to manage units. Backend-owned admission starts a fixed executable/unit
template; user input cannot choose unit properties or executable arguments.

Require per-worker cgroup memory limits and separate OOM handling so a worker's
failure does not terminate the backend. Disable swap and core dumps for services
holding credentials. Worker environments explicitly exclude deployment secrets.
In-process timeouts alone were rejected because a Lua allocation or parser can
exhaust backend memory before cooperative cancellation.

The policy worker protocol uses gRPC over a private Unix socket,
with no transport credentials or TCP listener. `worker serve --socket <absolute path>`
is a mode of the direct cleaner executable and bypasses OAuth initialization.
Use a clean supervisor environment; the command also clears inherited variables.
The unary `SubmitEvaluation(Evaluation)` request carries operation/account binding,
Lua script, fixed time, a logical source name and an opaque provider-specific archive
key. It contains no archive bytes or normalized records. Submission acknowledges
the backend-provided ID before fetching, parsing or computation. Accepted work
survives RPC disconnection under a worker-owned deadline. `GetEvaluation` returns
running or terminal status by job/account ID, with a complete typed operation plan only on
success; `CancelEvaluation` stops work and `DeleteEvaluation` releases terminal data.
There is one current unversioned contract, without streaming or compatibility APIs.

Only local sources are implemented: repeated `--source NAME=/absolute/directory`
flags open directory capabilities at startup. A source-provider interface owns key
validation, seekable object access and cleanup. Local keys are relative paths;
`os.Root` confines symlink and rename traversal, and nonblocking opens reject
non-regular files. Future S3 providers may interpret the same key as an object name
under operator-configured storage, with no job API change. No S3 implementation or
AWS dependency is included. Jobs cannot choose arbitrary roots, buckets or endpoints.

Open standard X ZIP exports, requiring `data/account.js` and `data/tweets.js`.
Validate ZIP metadata and the sole matching account before script execution. Expose
lazy `archive.datasets()` and named `account()`, `tweets()`, `likes()` and
`following()` accessors over top-level YTD dataset files, including consecutive
multipart files. Validate and normalize records into protobuf, remove YTD row wrappers and
retain string IDs. Optional null fields become absent Lua fields. Expose original
JSON structure only through explicit `archive.read_raw(name)`, including its
`json.null` sentinel. Remove the ambiguous `archive.read` API without a compatibility
alias. Parse literal array assignments as data; never
execute bundled JavaScript, extract files, follow links or read unrequested datasets.
A failed dataset read remains fatal across Lua `pcall`.

Scripts call a per-evaluation `plan` API whose methods derive from the protobuf
Operation oneof. Required scalar strings and target identifiers are positional;
remaining fields form options, while presence-sensitive profile updates use a
table. One descriptor rule drives runtime argument handling and editor generation.
The collector validates each contribution and aggregate bounds before mutation,
copies values, preserves operation order and retains caught API failures. The
interpreter accepts no returned values; successful completion publishes the native
protobuf plan, and any failure discards it. An empty script yields no operations.
There is no returned-plan decoder or compatibility path.

The generic job service imports no cleanup evaluator. The composition root installs
the cleanup application's `cleaner.plan` helper, which appends native deletion/repost
operations and a typed cleanup report for the existing review interface without
a protobuf-to-Lua round trip.
Other scripts need no report. The helper normalizes tweets, preserves cleanup/graph
semantics and verifies its report against actual archived records. The cleanup
backend retains the same protobuf payload and checks exact
operation/decision agreement. The generic worker/client share operation-payload,
unknown-field and message-bound validation. Every backend still owns authorization
and its supported executor set; planning a read or write performs neither.

Saved JSON documents have an independent 128 MiB encoded-size limit, including
file indentation, escaping and default fields. Validate the final formatted bytes
before replacing a file, bound reads before decoding, and retain the protobuf
payload limit after decoding. Optional null archive counters normalize to absent
values; numeric zero remains present and invalid counters still fail import.

Decision verdict: revise JSON results to typed protobuf operation plans. The result
is specifically executable X-operation intent, so an unconstrained return object
is broader than the requirement. Extending the protocol alongside a new executor
is appropriate. Typed payloads constrain targets and parameters while Lua remains
arbitrary computation over archive datasets. No JSON-result fallback, generic
parameter map, opaque bytes or compatibility service is retained. Preserve up to
50,000 operations and a complete result below 16 MiB, reserving metadata space;
conversion retains 64-level/one-million-value work bounds.

Lua editor declarations follow the same owners as runtime data. A project-owned
Go command generates plan methods and their exact argument types from operation
descriptors, and account/tweet/like/following types from archive descriptors.
The handwritten API declares named validated accessors, an explicit recursive JSON-value
raw read and the void cleanup helper. Typed projections include only
declared fields and allocate fresh tables. Dataset validation failures remain fatal
even if Lua catches the immediate error. The example calls plan methods directly,
without an operation array or returned result. Process acceptance verifies deterministic generation,
freshness, invalid-path behavior, Lua syntax and real example execution. Existing
Selene/StyLua checks remain available; LuaLS diagnostic checking is unavailable and
no new external dependency is introduced for this increment.

The composition root supplies a `script.Evaluator` struct through the worker's
`Evaluator` interface. Its `Options.Timeout` configures evaluation without adding
trusted options to job payloads. Zero inherits the parent context, a positive value
shortens the deadline and a negative value is rejected. The parent worker deadline
continues to bound the complete fetch/parse/evaluate job. Each call owns a fresh
interpreter; `Library: cleanup.New` installs a fresh cleanup helper and snapshot.
Cancellation and deadline errors retain their causes. No legacy policy evaluator
or function alias remains. The operator and worker use `--evaluation-timeout`.

Bounds: 128 KiB unary request, 64 KiB policy, 1 GiB ZIP, 10,000 entries, 8 MiB central
directory, 64 MiB per selected file, 256 MiB total selected decoded data, 50,000
records, 2 MiB per row, 32 MiB encoded inventory, 16 MiB complete result. ZIP64 and
multi-disk archives are unsupported. Typed account/tweet/entity decoding replaces per-field raw JSON maps. Explicit
and inferred references use stable reply/quote/repost array order. Cleanup result validation checks action/decision consistency directly; only the
source-owning cleanup helper performs membership checks against actual records. Policy error summaries contain
only the stage and source location, retaining original causes internally. The
client performs one terminal cleanup attempt under an independent bounded context;
cleanup failure warns without discarding a complete validated result.

The backend does not fetch the archive again. It relies on the trusted source-owning
cleanup helper for membership and cannot independently detect fabricated but
consistent inventory. The private socket is the process trust boundary.

The CLI and private server use this worker for all Lua evaluation. Planning has no
X client or credentials requirement. Explicit apply lazily initializes the mutation
client, requires a matching operator-selected account and rejects non-DELETE HTTP
operations. Archive identity is not proof of credential ownership; X enforces write
authority. Hosted OAuth identity verification remains future authorization work.

Decision verdict: proceed with pinned gRPC/protobuf and standard-library ZIP and
filesystem capabilities. This satisfies direct worker retrieval without a duplicate
backend parse or payload transport. Hard OS restrictions remain deployment work.

Keep one active fetch/evaluation, no waiting queue, at most 16 retained job
records, and 32 MiB retained successful payloads. Terminal results expire after
five minutes by default (`--result-ttl`, at most one hour); polling never renews
retention. Reject submissions when job retention is full and fail a result that
cannot fit, preserving existing retained results. Duplicate retained IDs cannot
replace work. The backend chooses globally unique IDs and owns durable status;
the worker owns only bounded transient attempts. Connection loss after acceptance
does not cancel them. Worker exit loses all transient jobs, which the backend
must mark interrupted rather than resume automatically. Long-lived worker units
must support result retrieval; the 90-second limit applies to fetch/parse/computation,
not retention. Hard operating-system enforcement remains deployment work.

### Canonical plan and evaluator ownership

Local simplification decision (2026-11-13): proceed with one protobuf operation
plan from planner through execution. A Go document envelope owns trusted account
and evaluation time and serializes its payload with strict protobuf JSON; it adds
no second operation schema or version compatibility. Adding these fields to the
script-controlled Plan would blur metadata ownership; retaining the old model
would preserve the redundant conversions identified in review. Saved files and
browser responses use the new envelope without accepting the old actions/version
format. Generic evaluation options belong to script.Evaluator. A per-evaluation
library factory installs cleanup behavior and checks its report; cleanup state
must never survive between calls. Rename the cleanup-specific package rather
than retaining a policy alias. The cleanup DSL and its graph semantics remain
active behavior, while only the existing deletion/unrepost executor is admitted.

### Canonical archive records

Archive schema decision (2026-11-13): proceed with api/archivepb/archive.proto as
the canonical typed input contract. The importer maps X's third-party JSON into
normalized protobuf account/tweet/like/following records, with typed timestamps,
reference/activity enums and optional exact counters. The importer constructs protobuf records directly through primitive source readers;
there are no Go JSON record structs or protobuf JSON-name aliases. The cleanup
graph's derived domain projection belongs only to cleanup. Lua uses
canonical protobuf field names and enum names; integers outside Lua's exact range
use decimal strings. Descriptor-driven projection and editor generation replace
the separate Go-reflection schema. Only literal standard X field names are read;
canonical-record and case-folded input spellings do not act as fallbacks. Records remain worker-local; the job protocol
still carries source references and typed operation results. ZIP retrieval and
decompression stay in the worker; no extraction responsibility moves to the backend.

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

Configured X export archives are the only inventory source. Live timeline reads and
refresh are removed. Evaluate policy uses one matching archive snapshot with no X
requests. Preserve Lua's lazy query behavior against the parsed local snapshot and
freeze one evaluation time for consistent relative dates and ages.

Hosted persistence will record owner, source/key, import time, available fields,
normalization version and coverage limitations. Narrower windows filter before
graph construction. Missing seeds, fields or coverage require explicit archive
replacement; they never trigger a network read or merge sources. Only explicitly requested YTD datasets are decoded; media binaries and external
links are never loaded. The cleanup helper retains only its account/tweet projection. An export does not grant X authorization or establish
complete history or graph coverage.

The backend validates the worker's returned snapshot/plan before publication.
After a dispatched mutation, mark affected snapshots and other unconsumed plans
dirty, including when the result is unknown. Do not alter the consumed execution's
original plan or disable explicit recovery. A dirty snapshot needs an explicit
replacement import before another evaluation; display its status.

### Money limits are independent of request windows

The monetary contract is [cleanup-spending](specs/cleanup-spending/spec.md).
Represent USD as integer microdollars, never binary floating point. Each explicit
cleanup owns one budget account: archive acquisition, evaluation, apply, cancellation,
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
2026-11-09, ContentManage is $0.005/request. The deletion category mapping
comes from X's Playground, not a live invoice. `/2/users/me` free treatment remains
unverified: reserve the published User category's $0.01 until confirmed. These are
dated implementation inputs, not permanent tariffs. Keep the provider ceiling and
application reservation semantics explicit rather than promise exact invoices.

### Proposed configurable operating defaults

These values complete unspecified implementation details for review; they are not
host configuration changes. Display applicable limits in the UI and reject startup
configuration whose aggregate worker/storage bounds cannot fit the VM.

| Limit                       | Initial value                                        |
| --------------------------- | ---------------------------------------------------- |
| Policy source               | 64 KiB UTF-8, bounded before parsing                 |
| API request body            | 256 KiB                                              |
| Archive object              | 1 GiB, fetched by worker                             |
| Archive decoded entries     | 256 MiB total; 10,000 entries; no nested archive     |
| Snapshot                    | 50,000 records and 32 MiB normalized data            |
| Cached snapshots            | 64 MiB/account; 512 MiB globally; seven days         |
| Policy declarations         | 256 rules; 10,000 explicit IDs total                 |
| Callback/graph work         | 1,000,000 callbacks and 1,000,000 edge visits        |
| Planning worker output      | 16 MiB, including decisions and actions              |
| Worker unit                 | 256 MiB memory; 100% CPU quota; 32 tasks             |
| Evaluation deadline         | 90 seconds, including fetch, parsing and computation |
| Admission                   | One active operation/account; two workers/VM         |
| Waiting queue               | 20 jobs globally; one queued job/account             |
| Unconsumed plan retention   | 24 hours; cannot outlive its snapshot                |
| Stopped execution retention | Seven days from last explicit apply/resume           |
| Application sessions        | Seven-day absolute lifetime                          |
| Browser progress logs       | Latest 1,000 entries; 4 KiB/entry                    |
| Durable operational logs    | Seven days; bounded rotation; no post content        |
| Rate control                | Five OAuth starts/IP/minute and 20/account/hour      |
| Total runtime data disk     | 2 GiB including database, WAL, snapshots, and logs   |
| Accounting disk headroom    | 128 MiB reserved; stop data admission below it       |
| Consistent backups          | Daily; seven-day retention in bounded storage        |

Limits apply to the owner too. Enforce total disk use across database, WAL, retained
execution records, snapshots, and logs, reserving free space for ledger commits.
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
- Missing archive coverage can produce an incomplete graph: expose coverage, fail
  missing seeds, and require explicit replacement rather than live read fallback.
- Price changes and uncertain responses can make estimates differ from invoices:
  use conservative dated tariffs, durable reservations, and operator reconciliation.
- Account isolation and archive intake expand the attack surface: enforce every
  read/write boundary on the server, escape text, and test hostile fixtures.
- Single-VM storage failure affects all users: protect the ledger, back up SQLite
  consistently with bounded retention, and exercise restoration before rollout.

## Migration Plan

1. Keep the existing CLI and private `server serve` command unchanged while adding hosted
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

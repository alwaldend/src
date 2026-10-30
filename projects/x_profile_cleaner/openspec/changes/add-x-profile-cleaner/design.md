## Context

See proposal.md for motivation. The repository pins Cobra and GopherLua.
The project is new; no existing behavior requires migration. The inventory
boundary and destructive operation distinction make an explicit design useful.

## Goals / Non-Goals

**Goals:** Keep Lua as a configuration language with predicate registration,
documented editor shims, bounded evaluation, and predictable plan/apply behavior.

**Non-Goals:** Historical search or archive import, Articles, Highlights,
deployment, credential provisioning, persistent execution journals, and live
validation against X.

## Decisions

- Separate command adaptation, application orchestration, X HTTP operations,
  and policy/graph evaluation into project-owned internal packages.
- Use the shared SDK OAuth 1.0a signer and its four credential fields. The
  cleaner runner reuses the Article uploader's existing injection configuration;
  its executable consumes those environment variables directly. Extend the
  existing credential-consumer package group rather than duplicating Vault paths
  or provisioning a separate credential. Sign each request after quota admission
  so waits cannot age its timestamp; retain credentials outside Lua and plans.
- Use the shared generated X SDK, with an authenticated transport that bounds
  responses and checks mutation confirmation presence before SDK decoding.
  Keep a 30-second network deadline separate from cancellable quota admission.
- Load activity kinds only when selected or queried. Cache mixed timeline pages
  requested at the maximum size of 100, without excluding replies or reposts.
  Normalize only requested author-owned kinds; keep repost action and source
  identifiers distinct. Include note-post full text; ignore article records.
  A caught Lua collection error still invalidates the entire plan.
- Build reply and quote adjacency independently on demand by reversing own
  references. No foreign parent node is required for an own reply to remain a
  candidate. Snapshot per-kind predicate IDs before callbacks so cross-kind
  queries cannot alter active iteration.
- Allow set_time_window to declare independently optional RFC3339 discovery
  bounds before inventory loads. Use inclusive-start/exclusive-end activity
  timestamps and pass the same bounds to every SDK page. Scope cached views,
  queries, and graph traversal to that window; reject late changes and do not
  perform an unbounded fallback. Reject out-of-window server records locally.
- Index owned edit-history identifiers as aliases of the collected current
  activity. Resolve selection seeds, queries, and graph references through
  that index; reject ambiguous aliases. Actions retain the collected ID.
- Describe each Lua operation's item kind, preservation scope, and deletion
  level once. Use one breadth-first preservation walk with a visited set;
  post preservation follows reply descendants; thread preservation additionally
  follows later quotes and own reply parents. Thus a quoting reply connects the
  thread it replies to with the quoted thread, without following quoted ancestors.
  Edge callback endpoints describe traversal direction, including parent steps.
  Each reached node and eligible outgoing edge is evaluated once per request.
- Register predicates through `delete_post_when`, `force_delete_post_when`,
  and equivalent repost and preservation functions. Direct preservation accepts
  an optional inclusion function or an options table with `include_post` and
  `follow_edge`. Query helpers return fresh Lua records, so caller changes cannot
  mutate the Go inventory. Callbacks must return booleans.
- Resolve selection sets after the complete configuration executes. Preservation
  request callbacks are applied uniformly through recursive reply traversal.
  Graph visitation is separate from preservation state. Deletion never cascades.
- Accumulate ordinary and forced deletion as one ordered deletion level, with
  preservation tracked independently. Use named kinds and derive Lua reply/quote
  flags from references instead of keeping duplicate state.
- Direct selection functions accept an ID or a dense array of ID strings.
  Validate the complete array before registering requests. Empty arrays are
  no-ops; each valid entry has the semantics of an individual call with the
  supplied callbacks. Keep editable ID arrays at the start of the example.
- Share one application instance between command constructors. The command
  initializes it after flag parsing and handles terminal input/output; the
  application owns choosing between a configuration and a saved plan.
- Use one fixed evaluation time and stable action ordering. Plan records include
  version, account ID, evaluation time, and typed actions. Applying validates all
  action kinds, identifiers, and duplicates before HTTP mutation.
- The CLI has plan and apply subcommands. Inline apply shares the planner and
  prompts on stdin; saved-plan apply does not rerun policies. Both use the
  authenticated identity, refuse redirects, and cancel outstanding work on a failed action with
  a count of all confirmed completed actions. Rate-limit failures are reported without replay.
- Track independent endpoint quotas, using documented user-context defaults
  until valid rate-limit headers override them. Wait before dispatch after
  exhaustion, honor cancellation, and retain per-endpoint order and same-target dependencies. The timeline
  batches records in pages; deletion and unrepost operations have single targets.
- Establish command-level fixture failure cases before production code. A local
  HTTP service provides account identity, paginated inventory, and observable
  delete/repost requests; its emitted plan is the repeatable evidence artifact.

## Risks / Trade-offs

### Lua API review decision

The language-independent internal/planner package owns discovery, graph indexes,
rule descriptions, selections, and plan construction. internal/policy owns the
Lua VM and callback registry, translating calls immediately into Go rule data.
The planner invokes opaque callback identifiers through a supplied evaluator
using plain Go values. Lua values and interpreter phases never enter the planner.
Predicates remain deferred and sequential; shared Lua closure behavior is retained.

Proceed with small additions to the existing function API. Selection functions
accept an optional options table with a name; preservation options retain their
callback shorthand. Reject unsupported arguments and options before loading data
or registering a request. days_ago uses the fixed evaluation time and returns a
whole-second RFC3339 timestamp suitable for set_time_window.

Plans may carry informational decisions for loaded activities: dates, text,
final outcomes, matching rule identifiers/names, and one preservation seed per
matching declaration. Apply continues to consume only actions; explanations do
not authorize execution. CLI --explain prints these decisions separately from
the action list. Keeping explanations with the plan makes a saved plan
reviewable without fetching again; older plans without decisions remain valid.

Expose optional language, attachments, and normalized text entities from the
same timeline pages. Use note-post entities whenever note-post text replaces
the ordinary text. No foreign expansion or additional endpoint is introduced.
Preservation callbacks receive a fresh seed_id context; traversal edges carry
child/parent direction. A callback-free declaration shares visited nodes across
its seeds and records its first preserving seed for each item. Declarations with
callbacks keep independent per-seed traversal, because callback state and seed
context can change their results. Check cancellation inside Go traversal too.

This retains the accepted configuration style and existing graph semantics.
A lifecycle framework or raw mutable SDK objects would add complexity without
improving the supported cleanup rules.

- Limited discovery can omit old posts or descendants. The README and plan
  command document that coverage; complete discovery is explicitly deferred.
- A discovery time window can omit a preservation seed or a quote/reply
  connecting thread members. An explicit window bounds API reads and the graph;
  age predicates select candidates inside it. An absent direct
  seed fails planning, while absent connecting posts can leave a partial graph.
- A policy can loop or fail. Evaluate with a deadline and restricted Lua
  libraries; fail the entire plan before any live mutation.
- X schemas evolve. Keep adapter field handling local and reject incomplete
  owned records rather than invent ages or silently drop candidates.
- API failures can leave a partially applied list. Report the exact completed
  count and return failure; resumable execution state is outside initial scope.

### Independent deletion workers

The user's concurrency instruction revises global serial execution: use one
worker per mutation endpoint. Every action remains copied from the validated
plan; worker queues preserve relative list order and same-target dependencies
always point to earlier list entries. On failure, cancel both workers and wait
for them, reporting all confirmed completions. Requests already accepted by X
may complete despite cancellation. No new dependency or retry mechanism is needed.

### Operator diagnostics and confirmation

Use a command-owned standard-library slog text handler writing to stderr. Pass
the logger explicitly to application and API orchestration; keep interpreter and
planner data independent of logging. Report phase boundaries, timeline page counts,
quota waits/resumption, and individual planned mutations with serialized confirmed
completion counts. Do not log post text, credential values, headers, or response
bodies. Keep errors at the CLI boundary and retain stdout's existing data output.

Wait for confirmation input in a buffered-result goroutine while the command
selects on its context. On cancellation return immediately without closing the
caller's input; the blocked reader ends when stdin closes or the CLI exits. Keep
EOF and negative answers non-destructive and put the prompt on stderr.

### Local browser interface

Proceed with a Go HTTP adapter in internal/web and embedded HTML, CSS, and
vanilla JavaScript. Use the existing standard library, pinned Sass, Bootstrap,
and Puppeteer rather than adding dependencies or a separate frontend service.
The site's assets remain the owner of shared palette and typography; export a
compiled application stylesheet from those sources, without copying theme values.

The serve command binds only a loopback address, defaults to an available port,
and prints a private access URL once on stdout. A fresh random capability travels
in the URL fragment, is removed from the address bar, and is retained in browser
session storage for authenticated API headers. It never enters operational logs,
request URLs, Lua, or plans. Check Host and Origin, reject cross-origin requests,
restrict framing and executable content, and render account data as text.

Read --config once at startup when supplied; otherwise embed the maintained
example policy. The browser file picker reads a local file into an editable
textarea. API requests carry policy text and never arbitrary server file paths.
Reuse application planning from in-memory source and retain the existing deadline.

Keep one active asynchronous operation and one authoritative plan on the server.
Starting a new plan invalidates the previous plan immediately. Apply accepts its
plan identifier and unique action indexes, validates the complete selection, and
copies only those server-owned actions. It consumes the plan before dispatch, so
duplicate apply requests cannot replay it. Retain the submitted subset separately
with confirmed completion indexes. After errors or cancellation, Resume projects
only unfinished actions, rechecks the account, and requires the latest stopped job
ID. Unknown outcomes remain unfinished; a later HTTP 404 confirms the target is
already absent. Confirmed actions are never resubmitted by Resume. A new plan
replaces execution state, and a server restart discards it. Editing browser policy
disables initial apply until replanning but cannot change a resumed selection.
Expose optional repost source IDs in informational decisions for item descriptions.

Poll operation status and a bounded, cursor-based log buffer; fetch the plan
separately when its identity changes so polling does not repeatedly transfer post
contents. Status includes completion indexes to restore progress after page reload.
Use the same logger for stderr and the browser. Cancellation stops the job's
context; stopping serve cancels active work. Already accepted X mutations retain
the existing unknown-outcome caveat. A browser confirmation describes the exact
selected count before submitting apply. Runtime state is in memory only.

Reject the alternative of accepting an editable plan JSON from the browser:
server-owned actions make exclusions a strict subset and prevent fabricated
targets, stale plan replay, and accidental expansion at the HTTP boundary.

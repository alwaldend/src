---
statuses:
  - active
title: X profile cleaner
description: CLI and local browser interface for Lua-configured X profile cleanup
languages:
  - go
  - lua
---

# X profile cleaner

![Apply state machine](/docs/projects/x_profile_cleaner/assets/apply-state.svg)

Posts and reposts use separate workers and quota buckets. Resume keeps the
confirmed selection and skips completed actions within the same server session.

This Go CLI prepares and applies Lua-configured cleanup from an account's X export archive, covering posts, replies, quote posts, and repost actions. Planning never
deletes content. Lua configuration registers predicates and preservation
requests; Go owns collection, graph traversal, and execution. `internal/planner` consumes
plain rule descriptions and exchanges plain callback values with an evaluator.
`internal/script` owns generic evaluation; `internal/cleanup` supplies its cleanup
library and callback registry. The planner produces protobuf operations directly;
interpreter values never enter it.

## Application flow

![Application flow](/docs/projects/x_profile_cleaner/assets/architecture.svg)

The diagram shows the intended upload and evaluation flow. The actor uploads an
archive to the backend, which stores it, submits worker jobs, reads their results
and applies confirmed plans through X. The worker reads the archive's ZIP data.
Backend, worker and archive are grouped under `x_profile_cleaner`.

The current backend selects preconfigured archives; the upload and archive-write
endpoints shown here remain planned.

The diagram source is `docs/architecture.mmd`. Build its dark WebP with
`bazel build //projects/x_profile_cleaner:architecture_webp`.

## Mandatory safety rules

The maintained [safety requirements](openspec/changes/add-x-profile-cleaner/specs/profile-cleanup/spec.md)
are hard constraints:

- **Planning never mutates the X profile**, even for forced deletions. It only
  reads a configured X export archive and writes the requested plan output.
- **Policy selection is exact:** delete only when explicitly forced, or when
  ordinary deletion selects an item which is not preserved. Unselected items
  stay unchanged. Errors fail planning rather than broaden deletion.
- **The generated plan is the complete selected action list**, with the correct
  account, operation kinds, target IDs, and forced indicators. Repost removal
  targets the repost relationship rather than deleting its source content.
- **Apply executes only that list.** It validates the entire plan and account
  first, adds no actions, performs no deletion cascade, and does not rerun Lua.
  Independent endpoint workers preserve order within each endpoint and between
  actions on the same target. A failure cancels outstanding work and stops further
  dispatch; requests already accepted by X may still complete.

These rules specify behavior for the supplied policy and collected snapshot;
they cannot determine intent which is not expressed in that policy. Review
the generated plan before applying it. Historical discovery limits remain
documented below.

## Policy

Start with [`examples/policy.lua`](examples/policy.lua):

```lua
local function older_than_30_days(item)
    return item.age_days > 30
end

delete_post_when(older_than_30_days)
delete_repost_when(older_than_30_days)
```

The post predicate includes **all your replies**, including replies to other
accounts. Foreign content is ignored, never deleted, and never connects own
graph branches. Reposts are your own activities: their age is time since your
repost, and execution removes that relationship rather than its source post.
Articles and Highlights-specific handling are outside scope.

For chunked runs, set an optional discovery window at the top of the policy:

```lua
set_time_window("2024-01-01T00:00:00Z", "2024-02-01T00:00:00Z")
```

For a rolling period, use `set_time_window(days_ago(90), days_ago(30))`.
`days_ago` accepts finite nonnegative fractional days and returns a whole-second
UTC timestamp based on the same fixed evaluation time as `age_days`. It performs
no inventory read. Unsupported timestamp ranges fail planning.

The start is inclusive and the end exclusive, so adjacent periods do not
overlap. Either bound may be `nil`: use
`set_time_window(nil, "2024-02-01T00:00:00Z")` for an end-only window or
`set_time_window("2024-01-01T00:00:00Z", nil)` for a start-only window.
Omitting the function or both bounds leaves discovery unbounded. Bounds must
be RFC3339 timestamp strings at whole seconds; offsets
are normalized to UTC. Set the window before any query or nonempty direct
selection loads inventory. It can be revised before that first read, but not
inside callbacks or after data has loaded.

All Lua queries, graph traversal, and deletion candidates use that archive
window. Missing context never triggers an API request. Even forced deletions cannot
select an activity absent from that inventory. See the
[discovery boundary](#discovery-boundary) for the effect on thread preservation.

`preserve_post(id)` preserves the selected post and recursively its own reply
descendants. `preserve_thread(id)` follows all later quote branches and walks
own reply chains in both directions. If a reply to your `post1` quotes `post2`,
preserving the `post2` thread also preserves that reply, `post1`, and their
reachable reply/quote branches. It never walks back to a quoted ancestor.
Repost preservation uses
`preserve_repost(activity_id)`; source post IDs remain separate.

Direct deletion and preservation functions accept one ID string or a dense
array of ID strings, for example `preserve_post({"111", "222"})`. Empty arrays
do nothing. The complete array is validated before any request from that call
is registered. Each entry has the same behavior and callbacks as an individual
call. The example keeps editable ID arrays together at the top of the file.

Every preservation function also has a `_when(predicate)` form. It accepts an
optional inclusion function, or a table with `include_post` and `follow_edge`:

```lua
preserve_thread("1111111111111111111", {
    include_post = function(post)
        return not post.is_reply or (post.metrics.like_count or 0) >= 10
    end,
    follow_edge = function(edge)
        return edge.to_id ~= "2222222222222222222"
    end,
})
```

The IDs above are examples. Earlier edit-history IDs resolve to the collected
current post for selections, lookups, and graph edges. Plans use the collected
ID once. An unknown or ignored seed fails planning.
Inclusion filters apply uniformly to replies reached indirectly. Rejecting a
post skips its preservation, while traversal continues. Rejecting an edge
prunes that path; another accepted path can still reach the same node.
Traversal visitation and preservation state are separate. Requests with callbacks
walk breadth first independently for each seed, evaluating each reached post and
eligible outgoing edge once per request. Callback-free seeds within one declaration
share visited nodes, preserving the same union without repeating overlapping walks. Callbacks should use item and edge data rather than depend on the
order of visits. Edge `from_id` and `to_id` describe the traversal direction;
a `replied_to` step in thread preservation can point to an own reply parent.
`edge.direction` explicitly distinguishes `child` from `parent` traversal.
Inclusion and edge callbacks receive a fresh second argument with `seed_id`,
resolved to its collected current ID:

```lua
preserve_thread_when(function(post)
    return (post.metrics.like_count or 0) >= 100
end, {
    name = "Replies near popular seeds",
    include_post = function(post, context)
        local seed = get_post(context.seed_id)
        return math.abs(post.created_at - seed.created_at) <= 7 * 86400
    end,
})
```

All selection functions accept an optional `{ name = "Description" }` table.
Preservation options can combine that name with callbacks. Repost preservation
accepts inclusion only, because repost activities have no traversable edges.
Unsupported arguments, unknown options, and empty names fail before selections
are registered; a deletion function never silently ignores an extra filter.

Ordinary `delete_post`, `delete_repost`, and their `_when` predicates request
deletion unless preserved. The corresponding `force_delete_*` functions
override preservation for that selected activity only; deletion never
cascades. Selection requests accumulate independently of declaration order:

```text
delete = forced_deletion OR (ordinary_deletion AND NOT preserved)
```

An item receiving no deletion request remains unchanged without being marked
preserved. Multiple matching predicates combine as requests, not last-write
assignments.

## Lua data and editor declarations

[`api/cleaner.lua`](api/cleaner.lua) describes the runtime functions, callback
types, and fields for editor completion. Add that file to your editor's Lua
workspace library; do not execute it as part of the policy. Runtime bindings belong to
`internal/cleanup`; selection and graph behavior belong to `internal/planner`.

Callbacks receive fresh item records with `id`, `author_id`, `kind`, `text`,
`created_at` (Unix seconds), fractional `age_days`, `is_reply`, `is_quote`,
`metrics`, and typed `references`. Repost records additionally have `source_id`;
their `id` is the own repost activity ID. IDs are strings throughout. Missing
metric values are nil, not invented zeroes. The complete `full_text` supplied by the export is used. `language`, `attachments`, and `entities`
expose exported language codes, media identifiers, and hashtag, cashtag,
mention, and URL spans. Attachment/entity objects and their collection fields
are always tables; missing and empty collections both produce empty tables.
Missing scalar metadata, including language and individual entity fields,
remains nil. Entities describe the exported text; unrelated files and foreign records are never
expanded. Poll identifiers are unavailable from supported tweet records and yield
an empty collection.

`get_post`, `get_repost`, `get_posts`, `get_reposts`, `get_replies`, and
`get_quotes` expose only collected owned records. Reply and quote queries
return all **direct** own children of an own graph node. Mutating a returned
Lua table does not change inventory. There is no global profile object or
mandatory lifecycle entry point.

Registration creates Go rule descriptions immediately. The top-level Lua
configuration finishes before selection predicates run. Predicates execute
sequentially in declaration order against stable per-kind candidate lists;
preservation walks then execute. Lua closures retain their shared variables.
Queries block until the required inventory is available. This is deferred,
lazy evaluation; it does not run Lua callbacks concurrently.

Callbacks must return actual booleans. They can query data but cannot register
new selections while evaluation is underway. Configuration errors, callback
errors, and timeouts abort the entire plan before execution. Lua has base,
table, string, and math libraries, without filesystem, process, package-loading,
credential, or HTTP access. `print` and dynamic code-loading globals are absent.
The worker validates the entire selected archive before Lua runs. Queries load
post/repost views lazily from that snapshot, including inside callbacks; they
perform no X requests. Reply/quote indexes are built only when needed. Collection
errors abort planning even when Lua catches them with `pcall`.

The default client deadline is 10 seconds and includes source fetching, ZIP parsing,
and policy evaluation; `--evaluation-timeout` changes it. The worker separately enforces
its own deadline, at most 90 seconds.

## Authentication

The default Bazel runner reuses the Article uploader's
[Vault-backed injection configuration](../x_article_uploader/BUILD.bazel) and
OAuth 1.0a user-context credential. Its four environment variables are documented
under [Article uploader credentials](../x_article_uploader/README.md#credentials).
The existing AppRole and Vault path remain owned by that configuration; the
cleaner does not duplicate secret references or provision credentials. The
authorization is needed only for explicit apply. Planning requires no X credentials.

Requests use the shared SDK's OAuth 1.0a signer, with a fresh timestamp and nonce
after each quota wait. Credential values never enter Lua, plans, source, logs,
or command-line arguments. To run the executable with credentials already
injected, use `//projects/x_profile_cleaner/cmd/x_profile_cleaner` directly;
it expects the same four environment variables. Credential issuance and host
configuration remain outside this project.

The default endpoint is `https://api.x.com/2`. `--api-base-url` accepts another
HTTPS endpoint or loopback HTTP for isolated fixtures. Redirects are refused.
Run the CLI only with an endpoint you trust to receive authenticated requests.

## Credential-free archive worker

Launch a separate worker with the direct executable target, which has no Vault
injection wrapper:

```sh
bazel run //projects/x_profile_cleaner/cmd/x_profile_cleaner -- \
  worker serve --socket /absolute/private/directory/policy.sock \
  --source exports=/absolute/archive/directory
```

Create the socket directory with mode `0700` first. The socket has mode `0600`;
an existing socket is never replaced. Interrupt or SIGTERM stops the worker and
removes its socket. Supervisors should start the built executable with an empty
environment; worker mode also clears inherited environment variables and bypasses
the operator's OAuth initialization. It needs no X or transport credentials.
The default `//projects/x_profile_cleaner` runner still injects operator credentials
and is intended for the operator commands.

Configure one or more named local sources at worker startup:

```sh
bazel run //projects/x_profile_cleaner/cmd/x_profile_cleaner -- \
  worker serve --socket /absolute/private/directory/policy.sock \
  --source exports=/absolute/archive/directory \
  --source other=/absolute/other/directory
```

The [`ArchiveWorker` contract](api/workerpb/worker.proto) accepts one unary
`SubmitEvaluation` request containing operation/account IDs, fixed evaluation time,
Lua script, `source_name`, and `archive_key`. The key is opaque at the RPC boundary.
Each configured provider interprets its own keys. The local provider accepts relative
paths beneath its directory, confines symlinks with `os.Root`, and accepts only
regular files. Jobs cannot supply source roots or endpoints. A provider interface
owns validation, object access and cleanup; future S3 support can use object keys
without changing this RPC. Only local directory providers are implemented today.

The worker fetches the ZIP itself after acknowledging the job ID. Archive payloads
and normalized inventories never pass through the backend or RPC. Standard X ZIP
exports must contain `data/account.js` and `data/tweets.js`; consecutive
`data/tweets-part1.js`, `data/tweets-part2.js`, etc. are supported. Literal
`window.YTD.<category>.partN = [...]` wrappers are parsed as data, never executed.
Datasets are decoded only when requested. `archive.account()` returns the validated
account record; `archive.tweets()`, `archive.likes()` and `archive.following()` return
validated arrays without their YTD row wrappers, including consecutive multipart
files. IDs remain strings and optional null fields become absent Lua fields. Each
call returns fresh Lua tables. `archive.datasets()` lists available dataset names;
`archive.read_raw("like")` reads a dataset with its original JSON structure and
`json.null` sentinel, without promising a record schema. Raw dataset names use the
filename spelling; YTD wrapper names replace hyphens with underscores. Unrequested
datasets and bundled media remain unread; no external links are followed. Archives
are read-only inputs.

[`api/archive.lua`](api/archive.lua) and [`api/types.lua`](api/types.lua) supply LuaCATS
editor declarations. Add both files to the editor's Lua library. The generated `plan` methods collect typed operations without returning a value.
The named archive accessors return the canonical records from
[`api/archivepb/archive.proto`](api/archivepb/archive.proto). Lua uses protobuf
field names: account `id`/`display_name`, like `post_id`/`text`/`url`, following
`account_id`/`user_link`, and tweet `id`/`author_id`/`text`. Tweet `created_at` is a
timestamp table with `seconds` and `nanos`; activity/reference kinds use enum names.
Metrics distinguish missing values from zero. Exact archive integers outside
Lua's number range become decimal strings. Metadata collections are always tables.
`archive.read_raw` returns `ArchiveJSONValue[]` in X's original format.

`api/types.lua` is generated entirely from the operation and archive protobuf
descriptors. Archive records use the descriptor projection at runtime; plan method names and
argument shapes use those same schema owners. The importer reads literal X source keys directly into protobuf messages;
there is no parallel Go archive record model or protobuf JSON-name alias path. Update declarations from the repository root after either schema changes:

```sh
bazel run //projects/x_profile_cleaner/cmd/generate_lua_api -- \
  --out "$PWD/projects/x_profile_cleaner/api/types.lua"
```

The generated declarations describe archive records, exact method argument tables
and schema-derived `plan` methods. Runtime bindings and editor signatures use the
same positional-argument rule: a required scalar string or target identifier is
positional, remaining fields are an optional options table, and presence-sensitive
profile updates take a table. The declarations use
[LuaCATS classes](https://luals.github.io/wiki/annotations/).
Acceptance checks generation freshness/stability and runs the shipped example
through the real worker. The repository has no LuaLS checker; these checks do not
claim editor diagnostic verification.

Scripts build the current evaluation's plan through methods:

```lua
for _, like in ipairs(archive.likes()) do
    plan.unlike(like.post_id)
end
for _, account in ipairs(archive.following()) do
    plan.unfollow(account.account_id)
end
plan.delete_post("123", { forced = true })
plan.unrepost("456")
plan.create_post("A new post")
plan.read_post("789")
plan.update_profile({ description = "Updated description" })
```

Each method validates its exact arguments and appends one protobuf operation in
call order. IDs remain strings. Optional forced flags default to false. Profile
updates require an explicit `description`; an empty string requests clearing it.
Argument tables are copied, so later Lua mutations do not alter planned operations.
All methods return no values and perform no X requests. New operation families
need an explicit protobuf payload and a corresponding backend executor.

An empty script produces an empty plan. A bare `return` can end a script early;
returning any value, including a table or nil, fails evaluation. There is no
returned-plan decoder, legacy return fallback, JSON result field or generic
parameter map. Unknown methods/fields, incorrect types, invalid identifiers,
argument-count errors and exceeded operation/byte bounds fail the complete job.
Plan API failures remain fatal even when Lua catches them with `pcall`. Earlier
collected operations are never published after an evaluation fails.

The [example collection](examples/README.md) covers archive filtering, relationship
selection, repost source handling, reference reads and profile/post updates.
Every example runs through the real worker in
`//projects/x_profile_cleaner/test/worker:examples_test`; its JSON configuration
supplies input archives, fixed time and exact expected protobuf plans.
Scripts receive `context.account_id` and
fixed `context.evaluated_at` Unix seconds. RPC status binds account/time/source
independently of script code. Failed dataset validation or projection also aborts
even when caught. The interpreter provides base/table/string/math libraries and
no filesystem, process, package, dynamic-code-loading, credential or HTTP API.

The CLI and browser submit `cleaner.plan(<policy source>)`. That helper retains
the cleanup selection/preservation and graph semantics, adding native protobuf
operations and a cleanup report to the collector without a Lua result object.
The report's provenance and exact operation/decision agreement are checked before
publication; contradictory mixtures of cleanup and generic operations fail.
Generic scripts need no cleanup report. The current backend executor supports
post deletion and unrepost only; other operation families can be planned but
require an executor before use. The backend still authorizes and confirms the
selected plan before executing it.

The cleanup helper uses typed tweet and entity decoding. Reference arrays have stable reply,
quote, repost order, including inferred explicit quote links.

Go callers configure `script.Evaluator{Options: script.Options{Timeout: limit}}`
and call its `Evaluate` method. The worker additionally sets `Library: cleanup.New`
to install the cleanup DSL; a generic evaluator needs no cleanup library. The
factory creates fresh helper state per call. Each call owns fresh interpreter state. A zero
timeout inherits the caller's context; a positive timeout can shorten it but cannot
extend an earlier parent deadline or the worker's complete fetch/parse/evaluate
limit. Negative timeouts are rejected. Worker startup supplies the existing
`--evaluation-timeout` setting through these options.

The sole archive account must match `--account`. Post IDs, edit aliases, timestamps,
metadata and references are validated before cleanup policy execution. Reposts remain
separate activities: removal needs an explicit original post ID. Many standard
exports provide only `RT @...` text. Such unresolved reposts are excluded from the
own-post view; requesting a repost view in a window containing one fails rather
than inferring a target from text. The export is data, never proof of X authority.

Backend callers use `internal/worker.Dial`, `Client.Submit`, and `Client.Get`.
Submission returns the backend-provided operation ID without waiting for fetching
or evaluation. `GetEvaluation` returns `RUNNING`, `SUCCEEDED`, `FAILED`, or
`CANCELLED` by job/account ID. Only success contains a complete typed protobuf
`Plan`. `Client.Get` validates operation/account/time/source/key bindings, terminal
structure, message bounds and declared operation payloads. Unknown protobuf plan
fields and malformed nested messages are rejected. The backend still validates its
supported executor set, authorization and exact selected actions. The cleanup helper
checks its report against the actual archive; the backend relies on that trusted
helper because it does not fetch the archive again.

Accepted work survives submission disconnection and remains subject to the worker
deadline. Script failures report a bounded stage and source line when available, without
source text, error-object contents or archive contents.

`CancelEvaluation` stops running computation. `DeleteEvaluation` releases
a terminal result; clients cancel abandoned work and delete results after copying. Terminal-result cleanup is best-effort: a deletion RPC
failure emits a warning and leaves TTL cleanup to the worker, while preserving
the validated plan.
Duplicate retained IDs are rejected. Restart forgets all transient jobs without
replay; durable identity/status remains backend-owned.

Limits are 64 KiB UTF-8 policy source, 128 KiB unary request, 1 GiB ZIP object,
10,000 entries, 8 MiB central directory, 64 MiB per selected data file, 256 MiB
selected decoded data, 50,000 activities, 2 MiB per record, 32 MiB encoded activity
inventory, and 16 MiB complete protobuf result. ZIP64 and multi-disk ZIPs are
unsupported. Unsafe names, links, duplicate entries, malformed wrappers, CRC
failures, conflicting records and exceeded limits fail without a partial plan.
Directory count, byte span and variable record lengths are verified before the
standard ZIP reader allocates entries. It receives an immutable validated directory
and canonical footer, preventing alternate footer or ZIP64 locator interpretations.
Archive reads check cancellation throughout parsing.

Saved plan JSON has a separate 128 MiB document limit for field names, escaping
and indentation. The final formatted file is checked before writing, and loading
is bounded before decoding. The decoded plan retains the protobuf payload limit.

One evaluation runs per worker with no queue. It retains at most 16 jobs and
32 MiB successful payloads. `--result-ttl` sets terminal retention (default five
minutes, maximum one hour); polling never renews it. Full retention rejects new
jobs; excessive results fail without evicting existing results. Expiry and deletion
release capacity. `--evaluation-timeout` on `worker serve` bounds fetch/parse/evaluation
(default and maximum 90 seconds).

These cooperative limits do not establish hard memory or CPU isolation. Packaged
systemd/cgroup restrictions and hosted admission remain open in the
[hosted-service plan](openspec/changes/host-profile-cleaner-service/tasks.md).

## Local browser interface

```sh
bazel run //projects/x_profile_cleaner -- \
  server serve --config /absolute/path/policy.lua \
  --worker-socket /absolute/private/directory/policy.sock \
  --archive-source exports --archive-key export.zip --account 100
```

The HTTP backend uses [chi](https://github.com/go-chi/chi) for routing and
standard net/http middleware. Open the private link printed on stdout.
`server serve` binds to an available loopback
port; `--listen 127.0.0.1:8080` selects a fixed port. Non-loopback addresses are
rejected. The page uses the site's shared light/dark theme and embedded assets,
without fetching scripts, fonts, or styles from external services.

Edit the Lua policy in the page, or choose a file with **Load file**. `--config`
loads the initial policy from disk once; omitting it starts with the maintained
example. File selection reads browser-local text, never an arbitrary path on the
server. Loading and editing a policy do not execute it or modify its file.

**Run plan** evaluates the policy without mutations. Review each candidate's
text, date, target, and matching rules; uncheck items to exclude them from this
run. **Apply selected** asks for confirmation, then executes exactly the selected
subset of the server-held plan. Editing the policy requires another plan. The
browser submits action indexes, never replacement action objects. A new plan
invalidates the old one; apply consumes its plan before starting, preventing
replay after reload or a repeated click. Exclusions are temporary and do not
rewrite the policy. An interrupted run retains its submitted selection and confirmed
completion indexes. **Resume** continues only unfinished actions, rechecks the
account, and uses the original selection even if the editor has changed. It does
not fetch the archive again or rerun the policy.

Search filters the visible rows without changing their selection. **Select
matching** and **Exclude matching** affect only those rows; the page reports
selected items hidden by the filter because apply still includes them. Excluded
rows are marked explicitly. Long post previews are truncated; inline **Show more**
reveals the complete collected text before applying. Rule names share the post-ID
row. The editor stays visible, with **Load file** in its header. On desktop,
Policy and Logs fill the left column beside a full-height Review panel. The editor,
logs, and removal list scroll independently; the review header and controls stay
visible. On narrow screens, the panels stack vertically. Active progress and
results appear above the logs.

The page shows the same progress and quota-wait logs as stderr, retaining the
latest 1,000 lines with a 4 KiB limit per line. Only one operation runs at a time.
**Cancel operation** or Ctrl-C cancels pending work; X may already have accepted
in-flight deletions. Closing or reloading the page leaves the job running while
`server serve` stays open. A reload restores server status and logs, but unsaved
editor
changes and exclusions for an unsubmitted plan are local to that page. After an
interrupted apply, reloading restores its submitted selection and progress, and
**Resume** remains available. **Run plan** replaces that execution record. Resume
state exists only while the same `server serve` process remains running.

Policies, plans, and job status are held in memory for this server process.
The private link contains a per-run access key; keep it private. The page removes
the key from its address bar and retains it for that tab's session. HTTP APIs
require it and reject foreign hosts and origins. OAuth credentials remain in Go
and are never sent to the page. Restarting `server serve` requires its new access link.

## Commands

Run the Bazel command from the repository root, supplying absolute input and
output paths because Bazel may execute within runfiles:

```sh
bazel run //projects/x_profile_cleaner -- \
  plan --config /absolute/path/policy.lua --out /absolute/path/plan.json \
  --worker-socket /absolute/private/directory/policy.sock \
  --archive-source exports --archive-key export.zip --account 100

bazel run //projects/x_profile_cleaner -- \
  apply --plan /absolute/path/plan.json --account 100

bazel run //projects/x_profile_cleaner -- \
  apply --config /absolute/path/policy.lua \
  --worker-socket /absolute/private/directory/policy.sock \
  --archive-source exports --archive-key export.zip --account 100

bazel run //projects/x_profile_cleaner -- \
  apply --config /absolute/path/policy.lua --yes \
  --worker-socket /absolute/private/directory/policy.sock \
  --archive-source exports --archive-key export.zip --account 100
```

`plan` writes a private, atomically replaced JSON file; omit `--out` to emit
JSON on stdout. The document contains the archive `account_id`, fixed `evaluated_at`
and a `plan` payload serialized with protobuf JSON field names. That payload has
`operations` and an optional `cleanup_report`, matching the worker result and
browser response. The old `version`/`actions` format is rejected. Post deletion
targets the own post ID; an unrepost action targets the source post ID. Repost
removals appear before post deletions, with stable ID ordering in each group.
Execution can interleave the two endpoint groups.

For example, a saved document with one deletion has this shape:

```json
{
  "account_id": "100",
  "evaluated_at": "2026-11-13T00:00:00Z",
  "plan": { "operations": [{ "delete_post": { "post_id": "123" } }] }
}
```

Account and time belong to the document rather than the script-controlled payload.
The backend retains protobuf operations through selection, resume and execution;
there is no separate cleanup action model or bidirectional plan conversion.

Saved-plan apply validates the entire plan, requires `--account` to match, and executes the list without rereading the archive
or Lua file. There is no identity lookup; X enforces ownership for each mutation.
Providing a saved plan is explicit execution, without an additional prompt.
Inline apply uses the same planner, displays the exact list, and asks for `y`
or `yes`. Declining or reaching EOF performs no mutations. `-y` and `--yes`
skip only that inline confirmation. Config and saved-plan inputs are mutually
exclusive. Ctrl-C interrupts confirmation immediately even while stdin remains
open, exits unsuccessfully, and performs no mutation.

Timestamped progress logs and the confirmation prompt go to stderr. Logs show
archive policy evaluation, record counts,
prepared actions, and execution progress. Each mutation reports its action index,
operation, and target ID; confirmed completion counts increase across both
workers. Quota waits report the affected endpoint, wait duration, and resumption
time, followed by a message when the allowance becomes available. A stopped apply
reports confirmed progress; outstanding actions can still have an unknown outcome.
Progress logs omit credentials, headers, response bodies, and post text. API
failures retain HTTP status and up to three problems, with only selected diagnostic
fields truncated to 160 characters each. Known credential values are redacted.

Stdout retains plan JSON or the command's action list and result summary. For
example, `plan --config /absolute/path/policy.lua > plan.json 2> progress.log`
keeps a parseable plan and a separate progress log. `--explain` deliberately
includes snapshot text on stderr in addition to these operational logs.

Plans also contain informational `decisions` for each loaded activity: its text,
date, outcome (`unchanged`, `preserved`, `delete`, or `force_delete`), matching
rule IDs/names, and one witness seed per preserving declaration. `--explain` on
plan or apply prints these on stderr, separate from executable actions. A rule
ID distinguishes declarations even when names repeat. Decisions never authorize
execution; older plans without them remain valid. Plan files contain snapshot
text and should be handled as account data.

Apply runs at most one post-deletion request and one repost-removal request
concurrently. Each worker keeps its actions in plan order and obeys its endpoint
quota. Actions sharing a target ID wait for their earlier planned operation to
succeed, even across endpoints. A mutation returning HTTP 404 is logged as already
absent, counted as completed, and skipped. On other failures, apply cancels both
workers, stops further dispatch, waits for outstanding calls to return, reports the
confirmed completed-action count, and returns nonzero. Already accepted server operations
may complete despite cancellation, and an interrupted response can leave their
outcome unknown. These actions remain unfinished for explicit UI resumption;
confirmed completions are never repeated by Resume. The completed count is not
necessarily a prefix of the plan. Requests are not automatically replayed. Rate-limit
failures include HTTP status, Retry-After, and reset metadata without dumping
the response body. This version has no persistent execution journal.

The application owns plan-source validation, saved-plan loading, worker jobs,
and execution. Commands share one application instance, configured after flags
are parsed. The browser adapter retains plans and selects subsets; it uses the
same planner, executor, and endpoint quota state as the CLI.

## Requests and rate limits

API operations use the repository's generated [`x_sdk`](../../internal/x_sdk/).
The transport supports only DELETE operations. Planning performs no account,
timeline or individual-post lookups. Post deletion and repost removal each accept
one target per request, so apply uses one sequential worker for each independent
endpoint bucket.

The documented OAuth user-context buckets are independent per endpoint:

| Operation              | Endpoint                                | Requests per 15 minutes |
| ---------------------- | --------------------------------------- | ----------------------: |
| Delete a post or reply | `DELETE /2/tweets/:id`                  |                      50 |
| Undo a repost          | `DELETE /2/users/:id/retweets/:post_id` |                      50 |

These are client defaults, verified against [X's rate-limit documentation](https://docs.x.com/x-api/fundamentals/rate-limits).
Valid response limit, remaining, and reset headers update each bucket. When its
allowance reaches zero, the next request waits until reset with a one-second
margin. Missing or invalid headers retain the local allowance and 15-minute
window. Quota state belongs to one CLI process; other clients using the account
can consume the same server quota. An unexpected 429 stops the operation without
automatic replay. Ctrl-C cancels quota waits. Each dispatched network request
has a separate 30-second deadline, which starts after admission.

At 50 mutations per 15 minutes, a large initial cleanup can take hours; the two
endpoint buckets can make progress concurrently. Billing and access limits are
separate from request windows. Apply has no overall policy deadline.

## Discovery boundary

Coverage is limited to the selected export and window. Exports may omit records,
fields, repost targets or relationships. Successful parsing does not establish
complete historical or graph coverage. Acquire a newer export explicitly when
fresh data is needed; the cleaner has no live read or refresh mode.
Restricting discovery to a time window can omit a preservation seed, a
connecting quote/reply, or other members of its thread. A directly requested
seed absent from the inventory
fails planning; a seed that is present can still have an incomplete graph if
connecting posts are missing. Thread preservation therefore applies only to
the graph discovered inside the configured window. Age predicates filter
deletion candidates within that inventory; they do not fetch missing context.
Chunked discovery intentionally accepts this limitation and does not establish
complete thread membership.

References: [post deletion](https://docs.x.com/x-api/posts/delete-post), and
[repost removal](https://docs.x.com/x-api/users/unrepost-post).

## Validation

```sh
bazel test //projects/x_profile_cleaner:test
```

The browser suite drives the built server with pinned Chromium and a local X
fixture. It verifies policy loading, non-mutating planning, exact exclusions,
confirmation, stale-plan rejection, access checks, cancellation, API failures,
reload behavior, escaped content, and visible quota logs. It produces light,
dark, and narrow-screen screenshots plus `serve-acceptance.json` in Bazel's test
undeclared outputs.

The CLI E2E suite runs the built command against a local HTTP fixture. It checks
archive collection, branching graph preservation, filtering, own replies to
foreign posts, repost age and identifiers, forced precedence, saved and inline
apply, confirmation and both yes flags, account mismatch, invalid policy and
plan input, direct ID arrays, edit-history aliases and conflicts, exact mutation
paths, concurrent endpoint dispatch, same-target dependencies, cancellation of
in-flight work, strict arguments, named decision reports, relative windows,
optional metadata, seed context, large overlapping threads, restricted Lua libraries, lazy cached queries, reply parents reached
through quotes, independent quota waits, cancellation, and missing or null
mutation confirmations. It also interrupts confirmation with stdin left open,
checks stdout/stderr separation, and emits `operator-progress.log` for inspection.
It independently verifies OAuth signatures on all
mutation requests, rejects incomplete credentials
before HTTP, and checks that output contains no fixture credential values. It emits
`cleanup-plan.json` in Bazel's test undeclared outputs. No live X credentials
or live deletion are used. Maintained change state lives in `openspec/`.

The CLI suite also emits `metadata-plan.json` after checking normalized
collections, optional scalars, and isolated nested policy tables.

The worker suite launches the built command without X credentials and exercises
the private gRPC socket, exact policy semantics, cached time windows, malformed
requests, bounded output, cancellation, deadlines, recovery, and forged worker
responses. It also checks multipart ZIPs, source confinement, malformed wrappers, CRC errors,
entry bounds, unresolved reposts, detached jobs, account-bound management,
retention, expiry and restart. It emits `worker-plan.json`, `source-archives.json`,
`async-lifecycle.json`, and `async-capacity.json`
in the test undeclared outputs. Its short
temporary socket directories use the test sandbox's temporary directory because
long runfile/workspace paths can exceed the Unix socket address limit; they are
removed after each test.

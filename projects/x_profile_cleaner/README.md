---
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

This Go CLI prepares and applies Lua-configured cleanup of an authenticated
X account's posts, replies, quote posts, and repost actions. Planning never
deletes content. Lua configuration registers predicates and preservation
requests; Go owns collection, graph traversal, and execution. `internal/planner` consumes
plain rule descriptions and exchanges plain callback values with an evaluator.
`internal/policy` owns the Lua interpreter and callback registry; interpreter
values never enter the planner.

## Mandatory safety rules

The maintained [safety requirements](openspec/changes/add-x-profile-cleaner/specs/profile-cleanup/spec.md)
are hard constraints:

- **Planning never mutates the X profile**, even for forced deletions. It only
  reads API data and writes the requested plan output.
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
be RFC3339 timestamp strings at whole seconds, on or after 2010-11-06; offsets
are normalized to UTC. Set the window before any query or nonempty direct
selection loads inventory. It can be revised before that first read, but not
inside callbacks or after data has loaded.

Every timeline page carries the same bounds. All Lua queries, graph traversal,
and deletion candidates use the resulting window; the CLI makes no unbounded
fallback request to retrieve missing context. Even forced deletions cannot
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
`internal/policy`; selection and graph behavior belong to `internal/planner`.

Callbacks receive fresh item records with `id`, `author_id`, `kind`, `text`,
`created_at` (Unix seconds), fractional `age_days`, `is_reply`, `is_quote`,
`metrics`, and typed `references`. Repost records additionally have `source_id`;
their `id` is the own repost activity ID. IDs are strings throughout. Missing
metric values are nil, not invented zeroes. The complete available long-post
text is used when returned by X. Optional `language`, `attachments`, and `entities`
expose returned language codes, media/poll identifiers, and hashtag, cashtag,
mention, and URL spans. Missing objects and fields remain nil. Long-post entities
correspond to the exposed long text; ordinary text entities are not substituted
when long-post entities are absent. Nested tables are fresh copies too. These
fields use the same timeline pages and do not expand foreign records.

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
Data loads on demand. An empty policy or empty ID arrays perform no timeline
requests. Post and repost views are normalized only when selected or queried;
reply and quote indexes are built independently when needed. Queries share one
cached timeline snapshot, including queries made inside predicates. Collection
errors abort planning even when Lua catches them with `pcall`.

The default policy deadline is 10 seconds and includes lazy reads and their
quota waits; `--policy-timeout` changes it. For example, use
`--policy-timeout 20m` when planning may need to wait for a 15-minute quota reset.

## Authentication

The default Bazel runner reuses the Article uploader's
[Vault-backed injection configuration](../x_article_uploader/al.lua) and
OAuth 1.0a user-context credential. Its four environment variables are documented
under [Article uploader credentials](../x_article_uploader/README.md#credentials).
The existing AppRole and Vault path remain owned by that configuration; the
cleaner does not duplicate secret references or provision credentials. The
application's user authorization needs read/write permissions.

Requests use the shared SDK's OAuth 1.0a signer, with a fresh timestamp and nonce
after each quota wait. Credential values never enter Lua, plans, source, logs,
or command-line arguments. To run the executable with credentials already
injected, use `//projects/x_profile_cleaner/cmd/x_profile_cleaner` directly;
it expects the same four environment variables. Credential issuance and host
configuration remain outside this project.

The default endpoint is `https://api.x.com/2`. `--api-base-url` accepts another
HTTPS endpoint or loopback HTTP for isolated fixtures. Redirects are refused.
Run the CLI only with an endpoint you trust to receive authenticated requests.

## Local browser interface

```sh
al bazel run //projects/x_profile_cleaner -- \
  serve --config /absolute/path/policy.lua
```

Open the private link printed on stdout. `serve` binds to an available loopback
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
not fetch a new timeline or rerun the policy.

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
`serve` stays open. A reload restores server status and logs, but unsaved editor
changes and exclusions for an unsubmitted plan are local to that page. After an
interrupted apply, reloading restores its submitted selection and progress, and
**Resume** remains available. **Run plan** replaces that execution record. Resume
state exists only while the same `serve` process remains running.

Policies, plans, and job status are held in memory for this server process.
The private link contains a per-run access key; keep it private. The page removes
the key from its address bar and retains it for that tab's session. HTTP APIs
require it and reject foreign hosts and origins. OAuth credentials remain in Go
and are never sent to the page. Restarting `serve` requires its new access link.

## Commands

Run the Bazel command from the repository root, supplying absolute input and
output paths because Bazel may execute within runfiles:

```sh
al bazel run //projects/x_profile_cleaner -- \
  plan --config /absolute/path/policy.lua --out /absolute/path/plan.json

al bazel run //projects/x_profile_cleaner -- \
  apply --plan /absolute/path/plan.json

al bazel run //projects/x_profile_cleaner -- \
  apply --config /absolute/path/policy.lua

al bazel run //projects/x_profile_cleaner -- \
  apply --config /absolute/path/policy.lua --yes
```

`plan` writes a private, atomically replaced JSON file; omit `--out` to emit
JSON on stdout. The plan contains schema version 1, the authenticated account
ID, a fixed evaluation time, and the exact typed action list. Post deletion
targets the own post ID; an unrepost action targets the source post ID. Repost
removals appear before post deletions, with stable ID ordering in each group.
Execution can interleave the two endpoint groups.

Saved-plan apply validates the entire plan, confirms the authenticated account
matches, and executes the list without rereading the timeline or Lua file.
Providing a saved plan is explicit execution, without an additional prompt.
Inline apply uses the same planner, displays the exact list, and asks for `y`
or `yes`. Declining or reaching EOF performs no mutations. `-y` and `--yes`
skip only that inline confirmation. Config and saved-plan inputs are mutually
exclusive. Ctrl-C interrupts confirmation immediately even while stdin remains
open, exits unsuccessfully, and performs no mutation.

Timestamped progress logs and the confirmation prompt go to stderr. Logs show
account identification, policy evaluation, timeline pages and record counts,
prepared actions, and execution progress. Each mutation reports its action index,
operation, and target ID; confirmed completion counts increase across both
workers. Quota waits report the affected endpoint, wait duration, and resumption
time, followed by a message when the allowance becomes available. A stopped apply
reports confirmed progress; outstanding actions can still have an unknown outcome.
Progress logs omit credentials, headers, response bodies, and post text. API
failures retain HTTP status and up to three problems, with only selected diagnostic
fields truncated to 160 characters each. Known credential values are redacted.
Timeline errors identify the failing page and count already collected, including
HTTP-success responses containing partial errors. Such failures still abort the
complete plan without publishing partial actions.

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

The application owns plan-source validation, saved-plan loading, Lua planning,
and execution. Commands share one application instance, configured after flags
are parsed. The browser adapter retains plans and selects subsets; it uses the
same planner, executor, and endpoint quota state as the CLI.

## Requests and rate limits

API operations use the repository's generated [`x_sdk`](../../internal/x_sdk/).
Discovery requests 100 records per timeline page and caches the pages for the
planning run. It performs no individual post lookups. X returns posts, replies,
and repost activities through the same timeline; a repost-only policy still
receives mixed pages but does not normalize ordinary posts or build a graph.
Post deletion and repost removal each accept one target per request, so apply
uses one sequential worker for each independent endpoint bucket.

The documented OAuth user-context buckets are independent per endpoint:

| Operation              | Endpoint                                | Requests per 15 minutes |
| ---------------------- | --------------------------------------- | ----------------------: |
| Account identification | `GET /2/users/me`                       |                      75 |
| Timeline pages         | `GET /2/users/:id/tweets`               |                     900 |
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

At 100 records per page, 3,200 records need about 32 full-page reads. Deletion is
the bottleneck: 50 per 15 minutes is 200 per hour on average, per bucket, so a large initial
cleanup can take hours; the two mutation buckets can make progress concurrently. API access and billing limits are separate from these
request windows. Apply has no overall policy deadline once the plan is ready.

## Discovery boundary

Complete historical discovery is outside scope. The CLI paginates only the
available user timeline, which X documents as limited to the latest 3,200
posts. It does not import archives or use historical search. A successfully
constructed plan does not establish complete historical or graph coverage.
Partial response errors, conflicting records, incomplete owned records, and
repeated pagination tokens fail planning.

Restricting discovery to a time window can omit a preservation seed, a
connecting quote/reply, or other members of its thread. A directly requested
seed absent from the inventory
fails planning; a seed that is present can still have an incomplete graph if
connecting posts are missing. Thread preservation therefore applies only to
the graph discovered inside the configured window. Age predicates filter
deletion candidates within that inventory; they do not fetch missing context.
Chunked discovery intentionally accepts this limitation and does not establish
complete thread membership or bypass X's documented historical coverage limit.

References: [user timeline](https://docs.x.com/x-api/posts/timelines/introduction),
[post deletion](https://docs.x.com/x-api/posts/delete-post), and
[repost removal](https://docs.x.com/x-api/users/unrepost-post).

## Validation

```sh
al bazel test //projects/x_profile_cleaner:test
```

The browser suite drives the built server with pinned Chromium and a local X
fixture. It verifies policy loading, non-mutating planning, exact exclusions,
confirmation, stale-plan rejection, access checks, cancellation, API failures,
reload behavior, escaped content, and visible quota logs. It produces light,
dark, and narrow-screen screenshots plus `serve-acceptance.json` in Bazel's test
undeclared outputs.

The CLI E2E suite runs the built command against a local HTTP fixture. It checks
paginated collection, branching graph preservation, filtering, own replies to
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
requests, including pagination and mutations, rejects incomplete credentials
before HTTP, and checks that output contains no fixture credential values. It emits
`cleanup-plan.json` in Bazel's test undeclared outputs. No live X credentials
or live deletion are used. Maintained change state lives in `openspec/`.

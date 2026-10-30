## Purpose

Provide account-owned X profile cleanup through Lua selection callbacks and
reviewable action lists, with explicit preservation and deletion precedence.

## ADDED Requirements

### Requirement: Shared OAuth 1.0a authentication

The CLI SHALL authenticate user-context requests using the shared SDK OAuth
1.0a signer and the four credential environment variables used by the Article
uploader. The default runner SHALL reuse its existing injection configuration
without duplicating Vault paths or provisioning credentials. Incomplete
credentials SHALL fail before HTTP requests. Credential values SHALL NOT be
exposed through Lua, plans, command-line arguments, or diagnostics.

Each admitted request SHALL be signed after its quota wait, covering the actual
HTTP method, URL, and query parameters with a fresh nonce and timestamp. The
former OAuth 2.0 bearer environment variable SHALL NOT substitute for a missing
OAuth 1.0a credential field.

#### Scenario: Reuse the configured user credential

- **WHEN** the cleaner runs through its default Bazel target
- **THEN** the Article uploader's existing injector supplies its four credential fields
- **AND** the executable signs reads and exact planned mutations through the shared signer

#### Scenario: Missing credential with an obsolete bearer token

- **WHEN** any required OAuth 1.0a field is missing and a bearer token is present
- **THEN** the CLI reports the missing variable name without credential values
- **AND** no HTTP request is dispatched

#### Scenario: Authentication after a quota wait

- **WHEN** a request waits for an exhausted endpoint allowance
- **THEN** its OAuth timestamp and nonce are generated after the wait
- **AND** pagination parameters remain covered by the signature

### Requirement: Author-scoped inventory

The CLI SHALL collect available timeline pages for the authenticated account,
including own replies to foreign posts and own repost activities. Foreign
content and Articles SHALL be ignored. Historical completeness and Highlights
handling SHALL be documented as outside scope.

Documentation SHALL explain that restricting discovery to a time window can
omit preservation seeds or connecting quote/reply records and therefore cannot
establish complete thread membership. It SHALL distinguish an absent directly
requested seed, which fails planning, from missing connecting records, which
can leave a discovered seed with an incomplete graph. Age predicates select
deletion candidates within the configured discovery window after collection.

#### Scenario: Foreign reply target

- **WHEN** an own reply references a foreign post
- **THEN** the reply remains eligible while the foreign post is ignored

#### Scenario: Repost age

- **WHEN** a recent repost references an old foreign post
- **THEN** age predicates use the repost activity's creation time

### Requirement: Optional Lua discovery window

set_time_window SHALL accept independently optional RFC3339 start and end
timestamps. A missing bound SHALL remain open; no bounds SHALL mean unbounded
discovery. The start SHALL be inclusive and the end exclusive. Timestamps SHALL
be whole seconds on or after 2010-11-06, normalized to UTC; two supplied bounds
SHALL satisfy start before end. Invalid declarations SHALL NOT replace a valid
window. Window changes SHALL be rejected after inventory loading or in callbacks.

Every timeline page SHALL carry the same configured API bounds. All owned
inventory views, Lua queries, graphs, and deletion candidates SHALL remain
inside that window, including forced deletions. The CLI SHALL reject or exclude
out-of-window server records and SHALL NOT fetch an unbounded timeline to fill
missing graph context. A window declaration alone SHALL NOT load the timeline.

#### Scenario: One-sided discovery period

- **WHEN** only a start or only an end is supplied
- **THEN** every timeline request includes that bound and omits the other
- **AND** all queries and planned mutations remain within the supplied bound

#### Scenario: Adjacent chunks

- **WHEN** consecutive runs use one period's end as the next period's start
- **THEN** an activity at that instant belongs only to the second period
- **AND** repost timestamps describe the repost activity rather than its source

#### Scenario: Missing graph context

- **WHEN** a connecting own reply or quote is outside the configured period
- **THEN** it is unavailable to graph traversal and Lua queries
- **AND** no unbounded fallback request retrieves it
- **AND** a directly requested preservation seed outside the period fails planning

### Requirement: Lua policy functions

The CLI SHALL expose documented predicate registration, direct preservation,
and author-only graph queries with string identifiers and fractional age days.
Invalid callback results and callback failures SHALL abort planning. Lua
evaluation SHALL have no access to credentials or live mutation operations.

Inventory kinds SHALL load only when selected or queried. Empty direct arrays
SHALL NOT load inventory. Reply and quote indexes SHALL be built only when
their respective query or preservation traversal requires them. Cross-kind
queries inside callbacks SHALL return complete available data without changing
the iteration of the predicate currently executing. Collection failures SHALL
abort planning even if Lua catches the query error.

#### Scenario: Repost-only policy

- **WHEN** a policy selects only reposts and does not query posts
- **THEN** no post normalization or graph construction occurs
- **AND** mixed timeline data returned by the API may be cached for later queries

#### Scenario: Empty policy

- **WHEN** no selection or data query requires inventory
- **THEN** planning performs account identification without timeline requests

#### Scenario: Shared age predicate

- **WHEN** a policy registers the same 30-day predicate for posts and reposts
- **THEN** own posts, quotes, replies, and repost activities are evaluated

#### Scenario: Direct selection arrays

- **WHEN** a direct deletion or preservation function receives a dense array
  of ID strings
- **THEN** it applies the same operation and callbacks to each identified item
- **AND** an empty array registers no selections
- **AND** an invalid entry fails before any entry in that call is registered

### Requirement: Strict and explainable Lua rules

All runtime functions SHALL reject unsupported arguments. Selection options
SHALL accept an optional nonempty rule name and only callbacks supported by
that operation. Invalid calls SHALL register no partial selection. Each
declaration SHALL have a stable numeric rule identifier even when names repeat.

Generated plans SHALL include informational decisions for loaded activities,
their text and date, final outcome, and matching rule identifiers and names.
Each preservation match SHALL identify one seed that preserved the activity.
The CLI SHALL expose these with --explain. Decisions SHALL NOT add, remove, or
change executable actions, and saved plans without decisions SHALL remain valid.

#### Scenario: Extra deletion filter

- **WHEN** a deletion function receives an unsupported extra callback
- **THEN** planning fails before registering the selection or mutating X

#### Scenario: Reviewing a preserved activity

- **WHEN** an ordinary deletion rule and a preservation rule match an activity
- **THEN** its decision lists both rules and a preservation seed
- **AND** only the explicit actions list controls saved-plan apply

### Requirement: Fixed-time helpers and optional post data

days_ago SHALL accept one finite nonnegative number, return an RFC3339 UTC
timestamp at whole seconds based on the fixed evaluation time, and perform no
inventory read. Unsupported timestamp ranges SHALL fail. Lua records SHALL
expose optional language, attachments, and text entities from the collected
record, retaining absent values as nil. Note-post entities SHALL correspond to
the exposed long text; mutating nested returned tables SHALL NOT alter inventory.

#### Scenario: Relative discovery window

- **WHEN** a policy sets bounds using days_ago before querying inventory
- **THEN** both bounds derive from the same evaluation time on every page

#### Scenario: Rich own post metadata

- **WHEN** an owned post has language, media keys, or text entities
- **THEN** predicates can use those fields without fetching foreign records
- **AND** missing metadata remains nil instead of fabricated values

### Requirement: Preservation context and shared traversal

Inclusion and edge callbacks SHALL receive a fresh second argument with the
canonical seed_id. Edges SHALL identify child or parent traversal direction.
Callback-free seeds in one declaration SHALL share visited nodes. Custom
callbacks SHALL retain independent per-seed visits and filtering. Go traversal
SHALL honor the policy context deadline, including when no Lua callback runs.

#### Scenario: Seed-dependent inclusion

- **WHEN** a predicate selects multiple seeds with an inclusion callback
- **THEN** each traversal supplies its own seed_id and independent visitation

#### Scenario: Overlapping callback-free seeds

- **WHEN** many matching seeds reach the same thread nodes
- **THEN** each node is expanded at most once for that declaration
- **AND** the preserved set equals the union of independent traversals

### Requirement: Reply-connected thread preservation

Preserving a post SHALL recursively preserve its author-owned reply descendants.
Preserving a thread SHALL walk all later quote branches and own reply chains
in both directions, including the own reply ancestors of a reply quoting the
seed. It SHALL NOT follow quoted ancestors. An optional inclusion filter SHALL apply
to all preservation candidates; exclusion SHALL NOT stop traversal. An optional
edge filter SHALL prune only the rejected edge. Non-author nodes SHALL NOT
connect otherwise separate author branches.

Collected edit-history identifiers SHALL resolve to their owned current post
for direct selection, queries, and graph references. Ambiguous or malformed
edit-history identifiers SHALL fail planning before mutation.

#### Scenario: Edited preservation seed and references

- **WHEN** a collected post has an earlier edit ID used by a preservation seed,
  quote, or reply
- **THEN** both IDs resolve to the collected post and its forward branches
  remain connected
- **AND** an earlier ID does not create a second deletion action

#### Scenario: Branching quotes and replies

- **WHEN** two own posts quote a seed and an own reply responds to one quote
- **THEN** preserving the seed thread preserves both quotes and the reply
- **AND** starting at one quote does not preserve the seed or sibling quote

#### Scenario: Filtered intermediate node

- **WHEN** an inclusion filter rejects a reply which has a qualifying quote
- **THEN** the reply is not preserved and the quote remains reachable

#### Scenario: A quoting reply connects its original thread

- **WHEN** an own reply to post1 quotes post2 and the policy preserves post2's thread
- **THEN** the reply and post1 are preserved along with reachable own replies and later quotes
- **AND** an edge filter can reject traversal from the reply to its parent
- **AND** preserving the reply with preserve_post alone does not preserve post1
- **AND** a post quoted by post1 is not included merely through that quote reference

### Requirement: Exact policy selection and deletion precedence

For posts and repost activities the CLI SHALL resolve forced deletion above
preservation above ordinary deletion. Without a deletion request the item
SHALL remain unchanged. Registration order SHALL NOT change precedence.

For each collected activity, the final selection SHALL be exactly:

```text
selected = forced_deletion OR (ordinary_deletion AND NOT preserved)
```

An ordinary request SHALL NOT weaken a forced request. No age rule, fallback
deletion, ancestor deletion, or deletion cascade SHALL be implicit. Callback
failures, non-boolean callback results, invalid identifiers, and invalid
preservation options SHALL fail planning rather than widen the selection.

#### Scenario: Unselected activities

- **WHEN** no deletion rule selects an activity, or ordinary deletion selects
  an activity which is preserved
- **THEN** no action is generated for that activity
- **UNLESS** an explicit forced-deletion request selects it

#### Scenario: Forced deletion of a preserved post

- **WHEN** both preservation and forced deletion select a post
- **THEN** the plan includes deletion of that post only

### Requirement: Non-mutating planning

Planning SHALL NOT mutate the X profile. Inventory collection SHALL use
read-only API operations. Planning SHALL NOT invoke post deletion or repost
removal, including when the policy registers forced deletion. Lua evaluation
SHALL operate only on the collected snapshot and selection requests.

The plan command MAY write its explicitly requested local plan output and
normal command diagnostics. A planning failure SHALL NOT publish a partial
replacement plan file.

#### Scenario: Planning a destructive policy

- **WHEN** plan evaluates a policy selecting every available activity for forced
  deletion
- **THEN** the X API receives no mutation request
- **AND** the result is an action list for later explicit execution

### Requirement: Exact plan construction

The generated plan SHALL contain schema version 1, the authenticated account
ID, one fixed evaluation time, and exactly the final selected action list.
Selected posts SHALL produce delete_post actions targeting their collected
current IDs. Selected repost activities SHALL produce unrepost actions
targeting their source post IDs, never deletion of the source content.
Actions SHALL be unique by kind and target. Generated plans SHALL put repost
removals first, followed by post deletions, with stable target-ID ordering
within each group. A plan SHALL pass complete validation before publication.

#### Scenario: Selected and preserved activities

- **WHEN** a snapshot contains old selected posts, preserved posts, recent
  unselected posts, and selected repost activities
- **THEN** the plan contains exactly the selected post deletions and repost
  removals, with the correct target IDs and forced-deletion indicators
- **AND** it contains no action for any other activity

### Requirement: Exact plan execution

Apply SHALL validate the complete plan and authenticated account before any
mutation. It SHALL execute only the actions listed in that plan, in their
per-endpoint order, preserving dependencies between actions with the same target
ID and using each action's exact kind and target ID. It SHALL NOT add
actions, traverse the graph to expand deletion, reevaluate Lua, or recollect
the timeline. Forced indicators SHALL be descriptive metadata, not a reason
to expand the list. A failed action SHALL cancel outstanding work and stop further dispatch.
Already dispatched requests MAY complete; confirmed completion counts SHALL
include all confirmed successes rather than imply a completed list prefix. Mutation requests SHALL NOT be automatically replayed.

Requests SHALL use the repository-generated X SDK. Timeline discovery SHALL
request up to 100 records per page and cache pages within a planning run,
avoiding per-post requests for data already available from that timeline.
The documented mutation operations accept one target per request; apply SHALL
use one sequential worker per endpoint, allowing the two independent buckets
to progress concurrently. Endpoint-specific quota state SHALL delay a request before
dispatch when its allowance is exhausted. Waiting SHALL be cancellable and
SHALL NOT consume the HTTP request timeout. Unexpected API errors, including
429 responses, SHALL stop execution without automatic replay. An HTTP 404 from
a mutation SHALL be recorded as an already-absent target and skipped. Other
errors and missing or null success indicators SHALL NOT be treated as completion.

#### Scenario: Idempotent explicit apply

- **WHEN** a planned deletion or repost removal returns HTTP 404
- **THEN** apply logs that the target is already absent and continues
- **AND** HTTP 403, HTTP 429, server errors, and transport failures still stop it

#### Scenario: Known exhausted deletion allowance

- **WHEN** a successful deletion reports zero remaining requests and a reset time
- **THEN** the next deletion waits until reset before dispatch
- **AND** operations on another endpoint use their own allowance

#### Scenario: Cancellation while waiting

- **WHEN** cancellation occurs before a quota reset
- **THEN** no waiting request is dispatched and no later action is attempted

#### Scenario: Invalid later action

- **WHEN** a plan contains a valid first action and an invalid later action
- **THEN** apply rejects the complete plan before the first mutation

#### Scenario: Exact saved-plan execution

- **WHEN** a saved plan contains a subset of available post and repost actions
- **THEN** apply issues exactly those mutation requests and no others
- **AND** changing or removing the Lua configuration does not change that list

### Requirement: CLI planning and confirmation

Applying a configuration SHALL prepare and display a plan, ask confirmation,
and execute only after an
affirmative answer unless `-y` or `--yes` is supplied. Declining or reaching EOF
SHALL cause no mutations. Applying a saved plan SHALL require no additional
confirmation. Failed execution SHALL return a nonzero status and report
confirmed completed actions.

#### Scenario: Inline confirmation

- **WHEN** apply is run with a configuration and confirmation is declined
- **THEN** the generated action list is shown and no deletion occurs

#### Scenario: Saved plan

- **WHEN** apply receives a saved plan
- **THEN** no policy or timeline requests are required and only its listed
  post deletions and repost removals execute

#### Scenario: Independent mutation endpoints

- **WHEN** the plan contains unrelated post deletions and repost removals
- **THEN** both endpoint workers can dispatch concurrently within their own quotas
- **AND** each worker preserves its actions' relative plan order

#### Scenario: Cross-endpoint dependency

- **WHEN** post deletion and repost removal share a target identifier
- **THEN** the later planned operation waits for the earlier one to succeed

#### Scenario: Failure with another request in flight

- **WHEN** one endpoint fails while the other has an outstanding request
- **THEN** waiting work is cancelled and no further action is dispatched
- **AND** already dispatched operations may complete despite cancellation

### Requirement: Cancellable confirmation and operator progress

Inline confirmation SHALL observe command cancellation even while stdin remains
open without an answer. An interrupt SHALL promptly terminate with nonzero status
and no mutation. Prompts and structured progress diagnostics SHALL use stderr;
stdout plan JSON SHALL remain independently parseable.

Diagnostics SHALL identify discovery and evaluation phases, fetched page and
record counts, endpoint quota waits and resumption, and mutation dispatch and
confirmed completion counts. Completion counts SHALL remain accurate across
concurrent endpoint workers. Errors SHALL be reported at their handling boundary.
Progress SHALL NOT contain credentials, authorization headers, or post contents.

#### Scenario: Interrupt during confirmation

- **WHEN** an interrupt arrives while confirmation is waiting on open stdin
- **THEN** the command exits promptly without waiting for input or mutating X

#### Scenario: Observe a rate-limited cleanup

- **WHEN** apply waits for an endpoint quota and later completes planned actions
- **THEN** stderr explains the wait, its resumption, and confirmed progress
- **AND** stdout retains its documented output without diagnostic log records

### Requirement: Local browser workflow

serve SHALL bind only to loopback and serve a self-contained page using the site's
maintained theme. It SHALL load an optional CLI configuration file once or use the
maintained example. The page SHALL support a local file picker and policy editing,
planning, item descriptions with exclusion controls, explicit confirmation of the
remaining actions, apply, cancellation, and live logs. Editing a policy SHALL
require replanning before the page offers apply.

The server SHALL require a fresh per-run capability for APIs and reject foreign
Host/Origin requests. The capability SHALL appear only in the initial private
access link and authenticated browser storage/headers, never in progress logs,
Lua, plan data, or request URLs. API inputs SHALL be bounded and strictly decoded;
profile and policy data SHALL never be interpreted as page markup or scripts.

#### Scenario: Policy loaded through either input

- **WHEN** a file is supplied with serve --config or chosen in the page's picker
- **THEN** its Lua text appears in the editor and can be changed before planning
- **AND** the API never accepts an arbitrary server filesystem path

### Requirement: Authoritative browser plans

The server SHALL run at most one planning or apply job at a time, report its
status, and retain the completed plan. Apply SHALL accept only the current plan
identifier and unique valid indexes into that plan's actions. It SHALL execute
exactly that selected subset without rerunning policy or discovering new actions.
Starting planning SHALL invalidate an older plan, and starting apply SHALL
consume its plan before any mutation. The server SHALL retain the submitted
subset and confirmed action indexes for its lifetime. An explicit Resume request
SHALL execute only unfinished actions from that subset, rechecking the account
without reevaluating policy. Confirmed completions SHALL NOT be replayed.
Browser selection of no actions SHALL perform no mutation.

#### Scenario: Resume an interrupted selection

- **WHEN** execution fails or is cancelled with unfinished actions
- **THEN** the UI offers Resume, including after a browser reload
- **AND** resumption uses the original selection and excludes confirmed completions
- **AND** stale or concurrent resume requests are rejected by job identity
- **AND** unknown outcomes remain unfinished until a resumed request confirms success
  or reports the target absent
- **AND** a new plan replaces resume state, and stopping the server loses it

#### Scenario: Excluding a planned removal

- **WHEN** the user excludes an item and confirms apply
- **THEN** only the remaining original plan actions reach X
- **AND** the excluded post or repost remains unchanged

#### Scenario: Stale or fabricated request

- **WHEN** an apply request references an old plan, an invalid index, duplicates,
  or fabricated action fields
- **THEN** the complete request is rejected before any mutation

#### Scenario: Filtering a selection

- **WHEN** the user filters the removal list and selects or excludes matching items
- **THEN** only matching rows change selection and hidden selections are counted
- **AND** excluded rows explicitly say Excluded while their contents stay readable
- **AND** clearing the filter does not reset the selection

#### Scenario: Compact review

- **WHEN** a plan is displayed
- **THEN** post previews are truncated, with an inline control to reveal full text
- **AND** the policy editor remains visible with Load file in its header
- **AND** Policy and Logs fill the available desktop height in the left column
  beside a full-height Review panel
- **AND** the editor, logs, and removal list scroll independently while the review
  header, selection count, and controls remain visible
- **AND** rule names share the post-ID metadata row, wrapping on narrow screens
- **AND** no separate planned-count status row duplicates the selection count
- **AND** review remains usable in light and dark modes and on narrow screens

#### Scenario: Diagnosing incomplete collection

- **WHEN** X returns problem details while collecting a timeline page
- **THEN** planning fails with HTTP status, page number, collected count, and bounded
  problem fields rather than a generic partial-response message
- **AND** credentials and unrelated response fields are omitted
- **AND** no partial plan or mutation is produced

### Requirement: Browser operation visibility

The page SHALL display planning and execution state and cursor-based bounded
progress logs from the same logger as stderr. It SHALL remain responsive during
quota waits and offer cancellation. The server SHALL cancel active work when it
shuts down. Closing or reloading a tab SHALL NOT silently replay an operation.
Informational decisions MAY expose repost source IDs to describe reviewed items,
but only server-owned actions SHALL authorize execution.

#### Scenario: Cancel planning or apply

- **WHEN** the user cancels the current operation
- **THEN** its context is cancelled and status and logs describe the result
- **AND** an interrupted apply can resume its unfinished selection in that server session
- **AND** already dispatched mutations retain their documented outcome uncertainty

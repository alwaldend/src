## Purpose

Keep hostile or accidentally expensive policies and archives from exhausting the
shared web service while preserving failure without a partial executable plan.

## ADDED Requirements

### Requirement: Bounded untrusted inputs and work

The service SHALL enforce configurable limits before and during parsing and
evaluation for policy bytes, request bodies, archive bytes and entries, records,
rule registrations, explicit identifiers, callbacks, graph visits, and output.
ZIP directory limits SHALL apply to actual header count, complete directory span
and variable record lengths before standard-library entry allocation. Parsing
SHALL use the same validated metadata and honor cancellation during reads.
All bounds SHALL apply to the owner as well. Syntax validation SHALL NOT be treated
as proof of bounded execution. Limit violations SHALL fail the operation without
publishing a partial snapshot or executable plan or dispatching a mutation.

#### Scenario: Small source produces excessive allocation

- **WHEN** a syntactically valid policy allocates or traverses beyond its allowance
- **THEN** the operation fails within its resource limits and publishes no executable partial plan

#### Scenario: Compressed input exceeds its decoded allowance

- **WHEN** an archive is within its object-size bound but exceeds entry or decoded-byte limits
- **THEN** parsing stops before exhausting shared storage or backend memory

#### Scenario: Footer understates the ZIP directory

- **WHEN** a footer understates count or size, a record escapes its span, or unchecked footer/ZIP64 metadata could redirect the parser
- **THEN** the worker rejects the archive before allocating an oversized entry list or publishing a partial plan
- **AND** cancellation stops directory reads and a later valid job can use the released slot

### Requirement: Worker failure isolation

Untrusted parsing and evaluation SHALL execute under separate hard memory, CPU,
task, and elapsed-time controls. Workers SHALL receive no X credentials or general
network or host-file access. A worker crash or OOM SHALL NOT terminate the web
backend or grant its output authority without validation. Cancellation SHALL end
the active computation and release its slot after bounded cleanup.

#### Scenario: Worker exceeds its memory limit

- **WHEN** a policy worker is killed for memory exhaustion
- **THEN** login, status, and unrelated bounded jobs remain available
- **AND** the failed operation reports resource exhaustion without mutating X

### Requirement: Credential-free archive script protocol

The cleaner SHALL expose a separately launchable worker mode over private local
gRPC without requiring X credentials or initializing an X client. Evaluation
SHALL accept only bounded Lua script source, a configured source name and provider-specific archive key,
an expected account, an operation ID, and a fixed evaluation time. The backend client SHALL reject
malformed, oversized, incorrectly bound or internally inconsistent results.
The generic worker SHALL validate declared operation payloads and SHALL NOT execute X requests.
The cleanup helper SHALL verify its own result membership against the archive; the
backend SHALL NOT claim independent membership proof without reading that archive. Evaluation
SHALL use no X requests and SHALL return no partial plan on failure.
Every cleanup-library action SHALL match exactly one selected decision with the same
operation, target, and forced indicator. Every selected cleanup decision SHALL have its
matching action; unchanged and preserved decisions SHALL have no action.
The RPC SHALL use a bounded unary submission without payload streaming. Fetching,
parsing and evaluation SHALL share one worker deadline. Submission SHALL acknowledge
a job ID without waiting for computation; later account-bound lookups SHALL expose
running or terminal status and one complete typed protobuf operation plan only on success. Status
SHALL retain source/key bindings. Accepted computation SHALL survive submission
disconnection. Explicit cancellation SHALL stop it.
Transient jobs and retained results SHALL be bounded by count, payload bytes,
and non-renewing terminal TTL. Explicit deletion SHALL release only terminal
jobs. Worker restart SHALL forget transient state without automatic replay;
durable identity and status SHALL remain backend-owned.
The RPC SHALL expose one current unversioned service contract, without legacy
service aliases or version-negotiation fields.
Exported reference arrays SHALL have deterministic reply/quote/repost order.
Policy failures SHALL expose a bounded stage and source line when available,
without source, error-object or archive contents. Cleanup of a delivered terminal
result SHALL be best-effort and SHALL NOT discard its validated plan.
The result SHALL contain ordered typed operations with exactly one declared payload
per operation. The schema SHALL represent post deletion, unrepost, unlike, unfollow,
post creation/read and profile-description updates. Plan methods SHALL convert arguments directly to
protobuf fields without JSON result serialization, arbitrary parameter maps or
opaque payloads. Unknown methods/fields, missing arguments, invalid target IDs,
incorrect scalar types and excessive values SHALL fail without a partial plan.
Cleanup explanations SHALL be a separate optional typed report, unnecessary for
other scripts. The backend SHALL still validate its supported executor set and
authorization before executing any operation. New operation families SHALL extend
the schema alongside their executors, without compatibility aliases.

#### Scenario: Plan operations unrelated to deletion policies

- **WHEN** a Lua script reads multipart likes and following and collects unlike/unfollow, create, read and update operations through plan methods
- **THEN** the worker returns the exact typed protobuf payloads without cleanup declarations or X requests
- **AND** exact string identifiers and explicit profile-description clearing survive transport

#### Scenario: Malformed operation plan

- **WHEN** a script uses an unknown method/field, missing target or incompatible scalar, or a worker returns unknown protobuf plan fields or malformed operation payloads
- **THEN** the entire result is rejected before publication or execution

#### Scenario: Dataset failure caught by script

- **WHEN** a script catches a missing, corrupt, malformed or over-limit dataset read with `pcall`
- **THEN** the evaluation still fails without publishing a partial plan

#### Scenario: Missing and empty collection metadata

- **WHEN** a policy receives records with missing, empty, or populated metadata through the CLI or worker
- **THEN** collections are tables with the supplied entries, missing counters remain nil, supplied zero counters remain zero, and modifying a returned table does not change later reads

#### Scenario: Invalid archive source or payload

- **WHEN** a source/key is unsafe, account mismatches, the ZIP has unsafe entries, malformed literal wrappers, corrupt data, conflicting records or exceeded bounds
- **THEN** the worker rejects it without a partial plan and releases its slot

#### Scenario: Evaluate an archive without authorization

- **WHEN** a worker starts without X credentials and receives a valid source reference and policy
- **THEN** it fetches the archive and returns the exact plan over gRPC using the supplied evaluation time
- **AND** no archive payload passes through the backend and no query acquires missing data from X

#### Scenario: Worker returns an unbound or contradictory result

- **WHEN** a response names another operation, account, source/key or time, or contradicts its executable decisions
- **THEN** the backend client rejects the entire response before publishing a plan

#### Scenario: Policy timeout or cancellation

- **WHEN** the evaluation deadline expires or cancellation is explicitly requested
- **THEN** the worker stops cooperative evaluation, returns no partial plan, and releases its evaluation slot

#### Scenario: Submitter disconnects after acceptance

- **WHEN** a valid source reference has been accepted and its submitting connection closes
- **THEN** computation continues under its worker deadline and another connection can fetch status/result by the same account-bound ID

#### Scenario: Retained result capacity or expiry

- **WHEN** retained jobs/results reach their allowance or terminal TTL elapses
- **THEN** admission fails closed without evicting unexpired results, and expired/deleted jobs release their capacity without polling extending retention

#### Scenario: Worker explanation contradicts execution

- **WHEN** an action lacks a selected decision, a selected decision lacks its action, or their operation, target, or forced indicator disagrees
- **THEN** the cleanup backend rejects the entire response, including actions labeled unchanged or preserved

### Requirement: Typed Lua editor API

The Lua editor API SHALL declare plan, operation and payload types derived from the
protobuf schema, including mutually exclusive operation variants. Account, tweet,
like and following editor types SHALL derive from canonical archive protobuf
schemas. The importer SHALL validate and normalize X JSON into those messages;
the importer SHALL construct those messages directly from literal standard X keys
without Go JSON record models, protobuf JSON-name aliases or case-folded fallback.
Named `archive.account()`, `archive.tweets()`, `archive.likes()` and
`archive.following()` accessors SHALL validate their records and return unwrapped
records with exact string IDs. Optional null fields SHALL become absent Lua fields.
Typed records SHALL use canonical protobuf field names, timestamp objects and enum
names. Optional counters SHALL distinguish absence from zero. Archive integers
outside Lua's exact number range SHALL become exact decimal strings. The runtime
projection and editor declarations SHALL derive from the same descriptors without
a parallel Go-reflection record schema or compatibility aliases.
Explicit `archive.read_raw(name)` SHALL return original JSON-value arrays, including
the null sentinel, without claiming record validation. The plan method names and runtime/editor argument shapes SHALL derive from the
operation descriptors. The shipped operation example SHALL use these methods
without constructing or returning a plan object. The generated projection SHALL have a repository-owned update
command and freshness/stability verification. Runtime validation SHALL remain the
execution boundary; annotation generation and Lua lint SHALL NOT be described as
LuaLS diagnostic verification when that checker is unavailable.

#### Scenario: Schema-derived editor types and example

- **WHEN** the generation command runs and the maintained annotated example is evaluated against a synthetic standard export
- **THEN** generated declarations match their maintained owners and the example produces the expected typed operation plan
- **AND** failed generation does not replace existing output

#### Scenario: Typed dataset validation and fresh projections

- **WHEN** named accessors read multipart datasets or a script mutates a returned record
- **THEN** every returned record matches its declared schema, IDs retain their exact string values and subsequent calls return fresh tables
- **AND** malformed known records fail the evaluation without a partial plan even if Lua catches their read error

#### Scenario: Alternate source field spellings are not compatibility inputs

- **WHEN** account/like/following identifiers use canonical-record names instead of X source names, or a required tweet key uses alternate casing
- **THEN** the worker rejects the source without a partial result, including when Lua catches the immediate read error
- **AND** standard X records still produce the exact normalized protobuf operations

### Requirement: Configurable policy evaluator

Generic evaluation SHALL use script.Evaluator with explicit trusted caller
options and fresh per-call interpreter and optional library state. Cleanup bindings
SHALL belong to cleanup rather than a policy compatibility package. Its timeout option SHALL reject
negative values, inherit the parent context when zero and never extend an earlier
parent or complete worker deadline. Cancellation and deadline errors SHALL retain
their causes, and unsuccessful evaluations SHALL publish no partial plan.

#### Scenario: Evaluation options preserve caller limits

- **WHEN** an evaluator timeout or earlier parent deadline expires, or the caller cancels evaluation
- **THEN** computation stops without a partial plan and exposes the corresponding cancellation or deadline cause
- **AND** reusing the evaluator does not retain Lua globals from earlier calls

### Requirement: Bounded fair job admission

The service SHALL enforce one active account operation, bounded account/global
waiting queues, bounded global worker count, and request-rate admission independent
of monetary budgets. Owner exemption SHALL NOT bypass these controls. Excess work
SHALL receive an explicit capacity response without unbounded goroutine or queue growth.

#### Scenario: Repeated cached evaluations attempt to monopolize the VM

- **WHEN** one account submits more work than its active and queued allowance
- **THEN** excess submissions are rejected while other accounts retain admission opportunities

### Requirement: One typed operation plan through the backend

The planner SHALL produce the protobuf Plan directly. The backend, saved files,
browser selection/resume and executor SHALL retain that payload without converting
to a parallel actions/decisions model. A trusted account/time document SHALL wrap
the payload using strict protobuf JSON. The old version/actions file schema,
policy package aliases and test-only archive.Read wrapper SHALL NOT be retained.
All executor operations SHALL be validated before any X request, and selection
or resume SHALL NOT modify the original retained protobuf plan.
Encoded JSON documents SHALL have an independent 128 MiB limit, with final
formatting checked before saving and bounded file reads before decoding. The
decoded protobuf payload SHALL retain its existing size limit; a successfully
saved document SHALL remain loadable without reevaluation.

#### Scenario: Saved typed plan and obsolete inputs

- **WHEN** a valid protobuf document is saved and applied
- **THEN** only its listed supported operations execute for its matching account
- **AND** an obsolete schema, unknown payload, ambiguous operation or unsupported later executor operation fails before any X request

#### Scenario: Large saved report and escaped text

- **WHEN** a valid plan has a large cleanup report or text expanded by JSON escaping
- **THEN** its successfully saved document reloads under the encoded JSON limit
- **AND** apply executes exactly its selected operations without issuing X reads

### Requirement: Per-evaluation Lua operation collection

The worker SHALL inject plan methods derived from the protobuf Operation oneof.
Each method SHALL validate exact arguments and append a native typed operation in
call order. Required scalar strings or target identifiers SHALL be positional;
remaining fields SHALL be options, while presence-sensitive profile updates SHALL
use a table. Runtime bindings and generated editor signatures SHALL share the
same descriptor rule. No returned-plan decoder or compatibility fallback SHALL
remain. Scripts SHALL return no values; a bare return SHALL be permitted.

The collector SHALL own fresh state for each evaluation, enforce aggregate
operation/byte limits before mutation and copy argument values. Invalid API calls
SHALL fail the entire evaluation even if caught by Lua. Execution, conversion,
archive, cancellation or helper failures SHALL publish no partial plan. Native
cleanup output SHALL enter the collector directly and retain report provenance
and exact action/decision validation without a Lua round trip.

#### Scenario: Generic operation methods and isolated arguments

- **WHEN** a script calls all supported methods over an export and later changes its option tables
- **THEN** the completed protobuf plan contains exactly the operations in call order with their original argument values
- **AND** no X request occurs during evaluation

#### Scenario: Caught failure and removed return contract

- **WHEN** a script catches an invalid plan call or returns a value after collecting operations
- **THEN** the job fails without publishing earlier operations
- **AND** a subsequent evaluation starts with an empty collector and can succeed

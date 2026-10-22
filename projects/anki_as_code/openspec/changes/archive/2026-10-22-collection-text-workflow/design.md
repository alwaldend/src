## Context

This change adds reconciliation behind the merged CLI, public App, config store,
resource/state stores, exporter and collection lifecycle. Use Anki's typed gRPC
services for all collection access. Anki owns storage compatibility, resource
validation, normalization, migrations, generated identities and cards. Keep one
Go module, use the pinned upstream schemas, and add no dependencies.

## Package boundaries

App composes workflows through collection.Store and prepared connections.
internal/cli owns its private application interface; cmd/anki_as_code injects the
concrete App. Config filesystem operations belong to config.Store and its layout
implementation; protobuf TOML/JSON conversion belongs to codec. Exporter owns
native-to-config projection and direct file output, independently of reconciliation.
Resource stores provide generated clients and fresh native snapshots. State stores
own only namespaced persisted bindings through Dump, Import, Update and Replace.
internal/plan owns private resolution models, moves, ownership checks, comparison
and prepared operations. internal/reconcile executes the protobuf plan without
resolving configs or discovering resources.

## Identity and configuration

Use identity-free protobuf-backed TOML in the fixed hierarchy. Paths address
config resources, while native IDs and member bindings live in collection-backed state. Note GUIDs
remain native metadata. A new file creates an object; an explicit file/directory
or member move retains its identity and review progress. New note types and their
notes can be created in one apply. Never generate identities into or rewrite source
configs. Export does not initialize or write state; callers import existing objects
before editing. The selected state key isolates bindings and leaves other keys,
unknown native resources, scheduling presets, history and media unmanaged.

## Pure planning

The merged planning library exposes plan.Build(configs, snapshot, stored). All three values come
from callers. The function performs no I/O, starts no backend, owns no stores, and
changes no inputs. Clone declarations before resolving path references and native
note-type values before removing comparison hints. Moves and proposed bindings
are private planning logic; state.Store has no preview API. Return a protobuf
CollectionPlan with a report, prepared operations, desired bindings and expected
prior state. Runtime models stay private to internal/plan.

Compare managed resources directly instead of constructing an artificial applied
collection. Report resource additions, edits, deletions and explicit moves in a
deterministic path order. Compare known normal-card destinations. Do not simulate
Anki's validation, normalization, template processing or card generation. Mark
possible generated-card changes deferred. Normalized source values can continue
to differ from native values until callers update them.

Check config references and ownership that Anki cannot see: declared paths,
move collisions, named-field mapping, state member address ambiguity, and impacts
on unmanaged resources. Leave name syntax, note-type structural rules, required
members and template validation to Anki. Protect resident filtered cards and their
original normal decks; both referenced decks count as occupied for kind changes.
Only tracked omissions are deleted. Reject edits affecting untracked descendants,
notes or filtered-card identities before native mutation.

## Native application

Reconciler borrows resource and state stores for one prepared collection and loads
the expected stored state before writes. Follow prepared name releases, deck
operations and type migrations through Anki, then batch prepared note updates.
Anki generates GUIDs and card identities. Apply regular-card defaults and explicit
ordinal destinations before removing tracked decks. Filtered cards stay unmanaged.
Reapply presentation after structural migration where Anki rewrites templates,
and record actual native member identities. Persist changed bindings once through
state.Store.Replace. Use native operation-change flags to report mutation without
re-enumerating native identities after applying. Notes are indexed by native ID;
GUIDs stay in native messages rather than reconciliation models. Resolve and
compare each resource group in focused planning helpers.
Preserve deck targets belonging to untracked note-type templates before writes.

## Collection lifecycle and publication

The merged collection.Store constructs and owns worker/collection backends,
reserves source identity through cleanup, and lends a prepared connection.
Reconciliation starts or stops no backend. Existing databases are opened in place;
committed API writes can remain after a later failure. Packages use native restore
and native package export; successful changed scopes publish only their own source
archive. False, failed and cancelled scopes skip package publication. Do not add
manual archive parsing, SQL, database copies or another publication facade.

rebuild.Builder copies a base package into a private directory on the destination filesystem,
uses normal apply there, then atomically renames the completed package. It leaves
base and destination intact on failure and never retargets an open collection.
Plan performs no mutations or database copying. Normal native updates maintain
caches; the application does not force Check Database or parse HTML cache fields.

## Delivery and verification

PR #128 targets master and adds native execution and application integration on
top of the merged #168 planning library. Planning source, schemas and extraction
artifacts remain unchanged. PR #167 supplies direct backend ownership; its
verified completion records remain preserved. Existing merged config, codec,
state, export and CLI behavior stays authoritative.

Use native E2E for convergence, note-type structure/presentation, batch creation,
explicit moves, namespaces, untracked and filtered-card protection, history,
scheduling, presets, media, no-op archives and failure publication. Write behavioral
failure cases before their implementation and emit repeatable native/TOML/JSON
artifacts. Validate pure planning after backend shutdown and unchanged inputs.
Delivery receipts own exact-candidate quality, lint, builds, review and publication.

Member moves rename a field or template within one note-type file. They cannot transfer native identity to another note type. Resource moves handle whole-file path changes separately.

## Execution contract

internal/plan owns pure Build(configs, snapshot, state), reference resolution, ownership checks, member matching, temporary rename selection and comparison. It returns a protobuf CollectionPlan containing prepared operations, desired bindings, expected prior state and a change report. Its private runtime models do not cross package boundaries.

internal/reconcile owns native execution only. Apply accepts the prepared plan, checks its expected state namespace before writes, preserves native metadata, executes the prepared phases, binds Anki-returned IDs and persists resulting state. It does not load configs, discover a collection snapshot, resolve ownership or compare desired resources. Callers provide fresh snapshots and exclusive offline access to the same collection and namespace.

Anki still validates native semantics and owns card generation and structural migrations. Executor metadata reads and generated-card assignment reads execute the supplied intent; no business simulation is added. App composes snapshot/state capture, Build and Apply in a single collection scope. The command prints the plan report; the Go planning API returns the complete protobuf plan.

## Planner review corrections

Desired objects are indexed by configuration path and existing snapshots by native ID.
New objects have zero IDs; prepared note/type/deck references use paths until Anki
returns actual IDs. CollectionPlan separates expected_state and proposed_bindings.
Compilation computes resource differences and operations together, without parsing
report strings. Member matching uses stable IDs or original stored names, including
before explicit moves, and rejects ambiguous current mappings. Converged plans carry
no native operations; assignments are emitted only for destination changes or relevant
generation. Execution resolves paths directly from actual bindings without temporary-ID maps.

## Filtered cloze templates

All cloze card ordinals use native template zero. Filtered-card ownership checks
require that template's prepared original ordinal to remain zero, rejecting
replacement/removal while permitting explicit identity-preserving template moves.
Native fixtures cover resident cards at ordinals zero and one.

## Empty state representations

State-change comparison treats an absent resources table as an empty binding set. The captured expected state remains unchanged so execution can still detect stale state using the original representation. Planning has no backend lifecycle precondition; the fixture stops its backend to prove I/O independence, while App can keep one native scope through capture, planning and execution.

The planning diagram groups setup under Preparation, with App resource/state call edges inside a nested WithConnection region. Collection.Store owns connection lifecycle, while filesystem config loading stays outside WithConnection. Pure plan construction remains the prominent top flow with its output outside the application group.

## Execution data ownership and application routing

Reconciler creates a private applier that borrows resource clients directly. Its
plan stays read-only: only proposed bindings and changed note-type request values
are cloned. Native metadata is merged into those owned requests; note bodies,
reports and captured expected state are reused without copying. Failure fixtures
stop type creation, note creation and card assignment after earlier native writes,
verifying that neither supplied plans nor stored bindings change on error.

App shares config loading, resource/state capture and pure Build through withPlan.
Each operation retains one collection scope and uses the same captured clients for
execution. Config.Store owns configuration-root normalization. rebuild.Builder owns
source checks, temporary archive copying, atomic destination publication and cleanup;
App supplies the ordinary apply callback and home-path expansion. No additional
backend, collection, public handler or serialization contract is introduced.

## Replacing tracked resources

Untracked-name protection is the sole native collision guard during deck and note-type binding. Tracked names may be vacated by deletion or rename; desired-name uniqueness and the existing temporary-name release operations make replacements safe without an additional all-native-name scan.

## Native planning data

Captured resources retain their upstream protobuf types with ownership paths indexed separately. Desired configs resolve once into the existing PlannedDeck, PlannedNoteType, PlannedNote and PlannedCardAssignments messages. Ownership checks and compilation share resolved path/ID indexes, without rebuilding editable configs from native notes, cards or note types. Returned operations own their payloads independently of caller inputs. The pure Build API and executable CollectionPlan contract remain unchanged.

## Case-insensitive hierarchy ordering

Deck payloads use folded-name ordering, matching parent validation. A declared parent is a strict prefix of its folded descendants and is always prepared first, including when their raw names differ in case. Anki retains ownership of final name normalization.

## Shared bindings and move indexes

Preset aliases remain in the owned state clone using their existing protobuf type, eliminating conversion to and from a second ID map. Independent move endpoints allow direct key renames with destination checks against existing typed resource maps: a destination cannot match the source prefix or another move endpoint. Existing collision, idempotency and input-ownership guarantees remain unchanged.

## Root notes and exported paths

An empty creation-deck path delegates deck selection to Anki, as specified by PlannedNote. Repeated application of the same configs converges using the stored root path. Export derives directories from actual regular-card destinations and retains bound filenames; adopting a different exported path requires an explicit move. This is the same path-identity rule as any other resource relocation, without a new planner restriction on native creation.

## Direct planning projections

Desired decks use one path index; parent-before-child ordering belongs to emitted deck operations. Filtered-deck IDs are projected from native decks and desired kind overrides only when card assignments need them, without a second kind map. Moved lookup data uses Resources directly; expected State remains a separate owned snapshot. Typed protobuf cloning preserves detached payload ownership without handwritten type assertions.

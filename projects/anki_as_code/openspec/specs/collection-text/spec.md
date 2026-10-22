# collection-text Specification

## Purpose

Reconcile offline Anki collections through path-addressed configurations and persistent collection state.

## Requirements

### Requirement: State-driven lifecycle and moves

Reconciliation SHALL resolve explicit moves idempotently before object operations, create unbound declarations, and delete only tracked objects absent from configuration. It SHALL reject conflicting bindings, unsafe moves, and destructive changes affecting untracked/filtered resources before native mutation. Field and template moves SHALL rename members within the same note type, without transferring native identity between note types. Native APIs SHALL retain scheduling, history, presets, media, and metadata outside the managed surface.

#### Scenario: Rename a deck directory

- **WHEN** a moved entry accompanies a directory rename
- **THEN** the deck, descendant bindings, notes, and retained cards preserve native identities and source files remain unchanged

#### Scenario: Import an existing object

- **WHEN** import binds a config path to an existing native ID
- **THEN** its content/review state is unchanged and later reconciliation updates that same object

### Requirement: Retained cards survive deck removal

Plan and preflight SHALL reject deleting a deck referenced by retained normal cards
without a declared surviving destination. Omitted card assignments SHALL NOT
request deletion of those cards or their note. Rejection SHALL precede native
writes for both database and archive inputs.

#### Scenario: Retained root note without a destination

- **WHEN** a tracked note moves to root notes with no card assignments and its old deck is removed
- **THEN** plan and apply reject the request without changing native resources or bindings; assigning a surviving normal deck makes the move applicable

### Requirement: Converged database apply is a no-op

Apply SHALL compare editable note-type values independently of runtime hints and
unmanaged metadata, while retaining explicit member identity/order semantics.
A converged database apply SHALL leave source bytes and state unchanged.

#### Scenario: Apply unchanged exported configs

- **WHEN** exported bindings are installed and the unchanged configs are applied repeatedly
- **THEN** plan remains empty and database bytes remain unchanged

### Requirement: Native creation batches retain declaration identities

Reconciliation SHALL use Anki's existing AddNotes operation for new-note batches
and bind each returned native ID to its corresponding declaration. Notes SHALL
resolve existing native IDs through a single precomputed index. Command lifecycle
SHALL remain owned by the shared collection workflow.

#### Scenario: Multiple notes in different decks

- **WHEN** a batch creates distinct notes across declared decks
- **THEN** each path binds its own native note, content and generated-card assignments survive re-export, and subsequent planning converges

### Requirement: Empty deck before changing kind

Plan and apply SHALL reject conversion of an occupied normal or filtered deck to
the other kind before native writes. An empty deck SHALL remain convertible,
without altering cards, scheduling or history already moved to a normal deck.

#### Scenario: Convert a normal deck containing retained cards

- **WHEN** retained notes move to root notes without card assignments and their occupied normal deck becomes filtered
- **THEN** plan and apply reject the conversion without changing native resources or bindings; managed archive packaging follows the collection-store publication policy

### Requirement: Assign cards using the resulting deck kind

After converting an empty filtered deck to a normal deck, Apply SHALL treat cards generated there in the same operation as regular cards. It SHALL honor both inherited deck destinations and explicit ordinal overrides, persist their bindings and converge on the next plan. Existing reviewed cards SHALL retain scheduling and history.

#### Scenario: Convert and create notes together

- **WHEN** an empty filtered deck becomes normal and a new note declares regular-card assignments in that deck
- **THEN** plan reports the conversion and note creation, apply succeeds for databases and packages, and each generated card receives its declared destination

### Requirement: Native execution consumes a prepared plan

Apply SHALL execute the prepared protobuf plan without loading configs, resolving
references, discovering a resource snapshot or recomputing differences. It SHALL
reject changed stored bindings before writes, preserve the input plan and bind
Anki-returned identities before persisting validated state. Provisional identities
SHALL NOT be persisted. The execution result SHALL report whether native operations
or persisted bindings changed, using native operation-change results without
fetching another complete note/card snapshot to verify identities. Apply and build
SHALL leave source configuration files unchanged. Callers SHALL own snapshot
freshness and exclusive offline access. Committed native writes can remain after
a later failure; archive publication SHALL follow its collection owner's policy.

#### Scenario: Execute a detached plan

- **WHEN** a captured plan is serialized and decoded before execution on the same prepared collection
- **THEN** execution performs the prepared work without native inventory reads, retains the original plan and records actual resource identities

#### Scenario: Use a prepared collection for reconciliation

- **WHEN** an application workflow executes a plan through a borrowed prepared connection
- **THEN** execution uses its supplied resource and state services without owning collection files, processes or archive publication

### Requirement: Empty namespace convergence

Execution of an empty plan SHALL preserve an absent state namespace without
reporting or persisting a synthetic binding change. Planning's absent-versus-empty
comparison belongs to the collection-planning capability.

#### Scenario: Execute an empty plan without creating state

- **WHEN** an empty plan is applied to an absent namespace
- **THEN** execution reports no collection change and leaves the namespace absent

### Requirement: Native root-note deck selection

Root-note creation SHALL leave initial deck selection to Anki. Reapplying unchanged configs SHALL preserve its native identity and converge. Export SHALL project actual card destinations without changing bindings; callers adopting changed paths SHALL declare explicit moves to preserve identities.

#### Scenario: Re-export a created root note

- **WHEN** Anki creates a root note with a fallback card and the caller later adopts its exported deck-relative path through an explicit move
- **THEN** the next apply retains the original native note identity under its new binding

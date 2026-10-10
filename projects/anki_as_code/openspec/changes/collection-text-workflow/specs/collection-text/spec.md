## Purpose

Make existing Anki collections editable through version-controlled TOML while
preserving identities, settings, templates, and original archive assets.

## ADDED Requirements

### Requirement: Lossless editable export

The CLI SHALL export all notes into TOML files with named fields containing
note identity, note type, tags, named field content, and every card's identity
and deck. Original HTML fields SHALL round-trip without content loss. Scheduling
and media SHALL remain in the base archive rather than public text.

Notes, decks, root configs, and card appearance SHALL be encoded as complete
TOML objects with the existing library's standard quoting and layout. Every
string value, including leading and trailing newlines, SHALL round-trip exactly.

#### Scenario: HTML quotes and literal delimiters

- **WHEN** a field contains HTML attributes, backslashes, apostrophes, or literal delimiters
- **THEN** the native encoder produces valid TOML without changing the content.

#### Scenario: Deck quoting remains lossless after reconciliation

- **WHEN** deck titles, paths, descriptions, or opaque string settings are exported or rewritten to record an allocated identity
- **THEN** the same native encoding applies and every parsed value is preserved.

#### Scenario: Multiple cards share a note

- **WHEN** a note generates multiple cards
- **THEN** one TOML file exposes shared fields and each card's deck independently.

### Requirement: Safe rebuild from a base archive

The CLI SHALL rebuild modern Zstandard collection packages from the original
archive and a desired-state text export. It SHALL preserve all scheduling columns,
templates, settings, and media except explicitly edited note
content, tags, deck membership, or declared card appearance. It MUST reject duplicate identities, unsupported schemas, malformed
fields, and invalid card ordinals before publishing an output. Apply and rebuilt archives MUST preserve review-log records unchanged.
Review history SHALL be excluded from text and reconciliation plans. Output SHALL be
written atomically and MUST NOT overwrite the base archive.

#### Scenario: Unedited round trip

- **WHEN** exported text is rebuilt without edits
- **THEN** note and card rows and all non-database ZIP entries are identical.

#### Scenario: Content and deck edit

- **WHEN** a field, tags, and card deck are changed in text
- **THEN** the rebuilt archive contains those edits with unchanged scheduling and media.

#### Scenario: Invalid export

- **WHEN** text duplicates a card or uses an unsupported schema
- **THEN** the command fails and leaves existing input and output files untouched.

#### Scenario: Text is independent of archive checksums

- **WHEN** desired text is built against a changed compatible collection archive
- **THEN** build accepts it, preserves that archive's unmanaged scheduling and history, and export omits archive hashes. Legacy base_sha256 metadata is ignored.

### Requirement: Plan and apply desired state

The CLI SHALL expose a read-only plan of note/card/deck additions, changes,
and deletions. Apply SHALL reconcile an offline collection database transactionally
to the text state, retaining existing card scheduling and initializing new cards.
Repeated apply SHALL converge to no changes. Missing note files SHALL delete those
notes and their cards. Deck configuration and note-type appearance SHALL be declared as editable text; collection configuration SHALL remain unmanaged.

#### Scenario: Add and remove cards

- **WHEN** a note file is removed and another valid note/card declaration is added
- **THEN** plan reports those changes and apply produces that desired state.

#### Scenario: Idempotent apply

- **WHEN** a reconciled collection is planned or applied again
- **THEN** it has no remaining note/card/deck/setting differences.

#### Scenario: Input changes during archive compression

- **WHEN** another writer changes the target archive while apply prepares its output
- **THEN** the final input hash check rejects replacement and preserves that writer's input.

#### Scenario: Deck titles exchange or reuse a removed title

- **WHEN** desired decks swap titles or reuse a deleted deck's title
- **THEN** reconciliation succeeds with Anki's unique name index and keeps retained deck identities.

#### Scenario: A card keeps its title override after deck titles swap

- **WHEN** deck titles swap and a card still explicitly targets its original title
- **THEN** plan reports the move to that title's resolved deck ID and reconciliation updates the managed card deck ID, preserving its scheduling. Cards resident in filtered decks remain outside managed state.

#### Scenario: A card follows a renamed deck identity

- **WHEN** a deck is renamed and the card's desired title resolves to the same deck ID
- **THEN** plan reports no card move and reconciliation leaves the card's modification and scheduling state unchanged.

### Requirement: New note identities

The CLI SHALL assign fresh note IDs, card IDs, and GUIDs to existing TOML note
files through generate-id. It SHALL accept file and directory arguments,
recursively discover notes, deduplicate overlapping paths, and preserve
filenames, content, tags, note types, ordinals, and deck overrides. It SHALL
reserve existing identities in the containing export and validate the batch
before changing files. It SHALL reject missing paths, malformed notes, and
symlinks. The CLI SHALL NOT provide new-note or the one-off organize-lessons command.

#### Scenario: Copied notes and directories

- **WHEN** generate-id receives copied note files and overlapping directories
- **THEN** each selected note receives unique identities once and a rebuilt
  collection retains unselected cards' original scheduling.

#### Scenario: Invalid batch

- **WHEN** a selected note is malformed
- **THEN** no selected file is changed.

### Requirement: Git-safe empty decks and normal card homes

Missing per-deck notes directories SHALL represent empty decks, including in
clean Git checkouts. Users SHALL create note files in these directories before assigning identities. Desired
card destinations SHALL be normal decks. Cards resident in filtered decks SHALL be excluded from exported declarations and card reconciliation; their complete database rows SHALL remain unchanged. Plan and apply SHALL reject declarations conflicting with their identities or note/template ordinals, and removal of their referenced notes or decks.

#### Scenario: Empty deck in a clean checkout

- **WHEN** a deck has a deck.anki.toml file but no tracked notes directory
- **THEN** plan and build accept the empty deck and a copied note with generated identities can become its first note.

#### Scenario: Filtered destination declaration

- **WHEN** desired text assigns a card's home deck to a filtered deck
- **THEN** validation rejects it before replacing any collection.

#### Scenario: Occupied filtered deck becomes normal

- **WHEN** a filtered deck containing cards is declared as a normal deck
- **THEN** plan and reconciliation reject the conversion without changing the input.

#### Scenario: Allocate new sibling positions

- **WHEN** new notes with multiple cards are reconciled
- **THEN** each note receives one position, the database nextPos cursor advances
  atomically, and runtime allocation state is omitted from declarative settings.

#### Scenario: Apply an archive with existing permissions

- **WHEN** an offline archive is reconciled in place
- **THEN** its file permission bits are retained.

### Requirement: Project documentation and landing page

The project SHALL expose accurate documentation and a registered landing page
through the main site's shared Hugo layout. Documentation SHALL distinguish
archive rebuilds from in-place reconciliation, describe complete desired-state
text, and use command examples with valid input relationships. The landing
SHALL describe current functionality and link to the project documentation.
The README SHALL list managed and unmanaged entities with reasons for each boundary.

#### Scenario: Browse the project

- **WHEN** the main website is built
- **THEN** the project index, landing, documentation, and existing Go/CLI taxonomy pages link to the Anki as code project.

### Requirement: Large collection support

The CLI SHALL NOT impose fixed total database or note-file size cutoffs.
Database copying and extraction SHALL stream to scratch disk. Format and
schema validation SHALL remain enforced independently of collection size.

#### Scenario: Large collection and field

- **WHEN** a supported SQLite/archive exceeds 512 MiB or a note exceeds 16 MiB
- **THEN** export and reconciliation accept it, preserve its contents, and converge.

### Requirement: Config-driven discovery

The CLI SHALL use collection.anki.toml as the root and resolve its decks_path,
note_types_path, and optional cardless notes_path relative to that file. Path references SHALL use flat decks_path, note_types_path, and notes_path strings. Deck markers SHALL use notes_path. Nested path tables SHALL be rejected as unknown schema fields. It SHALL
recursively find deck.anki.toml markers with explicit Anki title and notes_path,
and discover only \*.note.anki.toml notes recursively. It SHALL accept the root
config path or its containing directory. Filesystem layout SHALL NOT define
Anki titles. Referenced paths MAY be relative or absolute. Unrelated files SHALL
be ignored; symlinks and duplicate ownership SHALL be rejected.

#### Scenario: Rearranged export

- **WHEN** deck markers and notes move to arbitrary referenced directories
- **THEN** plan remains empty and rebuild preserves the same Anki titles and identities.

#### Scenario: Missing root reference

- **WHEN** a configured deck or note type path is missing
- **THEN** loading fails before publishing output.

### Requirement: Editable card appearance

The note_types_path root reference SHALL point to a directory of recursively discovered
\*.notetype.anki.toml files, one per existing note type. Each file SHALL contain
root ID/name/CSS and [[templates]] entries. Collection settings SHALL NOT be
exported, planned, or reconciled; their database records SHALL remain unchanged
except the computed nextPos cursor when allocating new-card positions. The
note type files SHALL expose template front/back
and browser formats, and browser font settings. Content SHALL use lossless
TOML encoding. Plan and reconciliation SHALL include those
edits while preserving unrelated protobuf fields, identities, card scheduling,
and exact unedited template/note type bytes. Note types and template ordinals
SHALL remain structurally fixed; text SHALL explicitly declare existing cards.

#### Scenario: Change card appearance

- **WHEN** CSS and template HTML change in text
- **THEN** build applies those appearance changes, retains unknown fields and
  card scheduling, and a subsequent plan is empty.

#### Scenario: Change a template without changing CSS

- **WHEN** template HTML or browser settings change while CSS remains identical
- **THEN** the parent note type is marked modified and pending sync, its CSS is preserved, and a subsequent plan is empty.

#### Scenario: Separate and movable note-type files

- **WHEN** the note type directory moves or note type files move into nested directories
- **THEN** all referenced note type appearance remains editable, and missing or duplicate declarations fail before mutation.

### Requirement: Upstream Anki protobuf schemas

The CLI SHALL use Go messages generated from integrity-pinned, unmodified
upstream Anki schemas for database protobuf configuration. The repository
SHALL NOT maintain copied schema definitions or hand-written field-number
encoders. Unchanged blobs SHALL remain byte-identical. Edited messages SHALL
retain unrelated known members and unknown field data.

#### Scenario: Preserve noncanonical configuration

- **WHEN** an exported collection contains explicit defaults, reordered fields, optional zero limits, nested unmanaged settings, or unknown fields
- **THEN** an unmodified plan is empty and build preserves original blobs; editing managed values preserves unrelated settings and converges.

### Requirement: Refresh exports and home-relative paths

The CLI SHALL accept current-user home-relative paths for collection operations, identity generation, and configured resource references. Repeated export SHALL refresh managed declarations beneath its output directory, preserve unrelated files, and leave existing declarations unchanged when input loading fails. Direct root config file inputs SHALL be accepted regardless of filename.

#### Scenario: Repeated export into a home-relative directory

- **WHEN** export is repeated against an existing directory using `~/...` inputs
- **THEN** stale declarations are removed, current collection declarations replace local edits, unrelated files survive, and planning against the export converges

#### Scenario: Explicitly named root configuration

- **WHEN** an existing root configuration is supplied as a file with a different name
- **THEN** its references are resolved relative to that file and collection operations remain available

### Requirement: Separate command results from status messages

The CLI SHALL write status messages and errors to stderr and plan JSON to stdout.

#### Scenario: Generate note identities

- **WHEN** identity generation succeeds
- **THEN** note fields remain unchanged and the status message is written to stderr without stdout content

#### Scenario: Missing required command input

- **WHEN** export, plan, apply, or build is missing a required flag
- **THEN** the command fails before collection operations and reports the error through stderr

### Requirement: Storage and domain package boundaries

The model package SHALL contain only data types, with no functions or methods.
The CLI SHALL put validation, planning, and Anki configuration/text encoding in
separate packages, independent of TOML storage, archive resources, SQLite
persistence, and operation orchestration. All appearance SQL SHALL execute through repository methods
within reconciliation's transaction. The domain, text, and archive packages
SHALL NOT depend on SQLite or on collection orchestration.

#### Scenario: Enforced dependency direction

- **WHEN** Bazel builds the collection CLI and its end-to-end suite
- **THEN** explicit package dependencies and visibility enforce the adapter boundaries, and the application exposes no raw SQL handles.

#### Scenario: Filtered cards are unmanaged

- **GIVEN** a note has one card resident in a filtered deck and one regular-deck card
- **WHEN** the collection is exported and its regular-deck card is moved
- **THEN** only the regular-deck card appears in text, and the filtered card row remains unchanged.

#### Scenario: Desired state conflicts with a filtered card

- **GIVEN** a card resident in a filtered deck
- **WHEN** desired text claims its identity or note/template ordinal, or removes a note/deck it references
- **THEN** both plan and apply reject the change without modifying the archive.

### Requirement: Protobuf serialization contracts

The first-party protobuf schema SHALL own editable TOML and plan JSON field
names and types. The CLI SHALL load TOML into generic values, reject unknown
fields and invalid values through ProtoJSON, and parse a generated message
from those values without a separate schema validator. Export SHALL use
ProtoJSON with protobuf field names and convert integer values to TOML numbers
without floating-point conversion. Optional-field omission SHALL follow
ProtoJSON presence semantics. Plan JSON SHALL use ProtoJSON directly, including
its string representation of 64-bit integers. Shared Go
models SHALL NOT duplicate protobuf-defined serialized fields or declare
independent TOML or JSON serialization tags. Configurations, cards, appearance,
and plans SHALL use generated message types directly. Runtime wrappers MAY
combine message pointers with nonserialized metadata.
Runtime-only state SHALL remain outside the serialization schema.

#### Scenario: Exact integer identity through schema parsing

- **GIVEN** a TOML note with a card identity greater than 2^53
- **WHEN** the CLI parses the note through its protobuf schema and exports it
- **THEN** the integer identity is retained exactly, without conversion through floating-point values.

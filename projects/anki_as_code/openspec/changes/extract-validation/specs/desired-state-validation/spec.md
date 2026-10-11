## Purpose

Validate resource identities and mutually exclusive deck kinds independently of collection storage and reconciliation.

## ADDED Requirements

### Requirement: Unique positive note identities

Desired-state validation SHALL reject nonpositive or duplicate note IDs independently of GUIDs and card identities.

#### Scenario: Distinct resources share a note identity

- **WHEN** two desired notes share an ID but have distinct GUIDs and card IDs
- **THEN** validation rejects the desired state.

#### Scenario: Nonpositive note identity

- **WHEN** a desired note has a zero or negative ID
- **THEN** validation rejects the desired state.

### Requirement: Exclusive deck-kind declarations

Deck encoding SHALL reject a normal declaration combined with a present filtered declaration, including an empty filtered payload. A filtered-only empty payload SHALL remain valid.

#### Scenario: Both kinds are declared

- **WHEN** a normal deck configuration and an explicitly empty filtered payload are supplied together
- **THEN** encoding rejects the conflicting declaration.

#### Scenario: Empty filtered configuration round trip

- **WHEN** a filtered-only empty payload is encoded and decoded
- **THEN** the filtered declaration remains present and normal remains absent.

### Requirement: Anki-owned derived field caches

The CLI SHALL defer sort-field and first-field checksum rebuilding to Anki and
SHALL document the required Check Database step after importing changed fields
or sort-field settings. It SHALL NOT interpret field HTML for derived caches.

#### Scenario: Existing note fields change

- **WHEN** the collection writer updates an existing note
- **THEN** its stored sfld and csum remain unchanged until Anki rebuilds them

#### Scenario: A new note is inserted

- **WHEN** the collection writer inserts a new note
- **THEN** its sfld is empty and csum is zero until Anki rebuilds them

### Requirement: Readable filtered deck configuration

Deck declarations SHALL expose filtered configuration as upstream protobuf
fields serialized to readable TOML. They SHALL require exactly one deck kind.
Filtered configuration edits SHALL preserve resident cards and temporary card
scheduling, and existing unknown protobuf data during reconciliation.

#### Scenario: Edit filtered settings

- **WHEN** rescheduling, a search term, limit, or ordering changes in TOML
- **THEN** reconciliation updates the filtered deck configuration and converges
- **AND** resident cards remain unchanged

#### Scenario: Opaque legacy configuration

- **WHEN** a declaration contains filtered_base64
- **THEN** loading rejects that unsupported field and requires re-export

### Requirement: Reuse upstream regular-deck configuration

Regular deck declarations SHALL use the pinned upstream Deck.Normal message
rather than a custom configuration projection. Known settings SHALL be readable
TOML and unknown fields SHALL be retained from existing archive configuration.
Scheduling presets referenced by config_id SHALL remain unmanaged.

#### Scenario: Edit optional regular settings

- **WHEN** a daily limit or desired retention override changes in TOML
- **THEN** reconciliation updates the upstream regular-deck configuration
- **AND** the scheduling preset remains unchanged

### Requirement: Editable note-type definitions

NoteTypeSettings SHALL reuse upstream note-type configuration and ordered Field
messages; TemplateSettings SHALL reuse upstream template configuration. Desired
state SHALL support creation, renaming, field/template structure and known
configuration edits, and removal of types without remaining desired notes.
Every note SHALL match its desired type's named fields; card declarations SHALL
remain explicit. Existing note/GUID/type identities and card ownership SHALL
remain stable. Existing card ordinals MAY change explicitly without resetting
scheduling. Structural changes SHALL mark the Anki schema modified.

#### Scenario: Create and restructure note types

- **WHEN** text creates a type or edits existing fields, templates, and settings
- **THEN** validation accepts a consistent complete desired state
- **AND** reconciliation preserves existing scheduling, history, and matched unknown configuration
- **AND** a repeated plan has no changes

#### Scenario: Invalid desired structure

- **WHEN** fields are missing/duplicated, template ordinals are noncontiguous, or a referenced type is absent
- **THEN** validation rejects the desired state before mutation

### Requirement: Explicit card declarations without content inference

Validation SHALL NOT parse HTML or cloze content to infer or freeze card sets.
Normal cards SHALL reference a declared template ordinal; cloze cards SHALL
use explicit nonnegative ordinals. Filtered-card template position and note-type
kind SHALL remain protected by the repository adapter.

#### Scenario: Literal cloze syntax in a normal note

- **WHEN** a normal note's field contains cloze-like text
- **THEN** validation accepts it without deriving card ordinals from content

### Requirement: Recursive unknown configuration preservation

Configuration preservation SHALL traverse upstream message descriptors and
retain opaque data in corresponding nested messages. Removing a nested message
or list item carrying unknown data SHALL fail explicitly. Unchanged wire blobs
SHALL remain verbatim.

#### Scenario: Remove unknown-bearing filtered search term

- **WHEN** a filtered declaration removes a search term carrying unknown fields
- **THEN** validation rejects the edit instead of silently truncating preservation

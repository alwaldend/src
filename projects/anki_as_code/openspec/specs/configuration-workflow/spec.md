# configuration-workflow Specification

## Purpose

Read and write path-addressed schema-backed Anki configurations independently of collection processing.

## Requirements

### Requirement: Path-addressed configuration

The collection SHALL declare version 1 in collection.anki.toml and use fixed note_types/, notes/, decks/, and nested child decks/ directories. Each deck directory SHALL require deck.anki.toml and derive its title from directory names. Parsed resource wrappers SHALL retain collection-relative source paths. Unknown fields, unsupported versions, unsuitable filesystem entries, and malformed syntax/types SHALL fail; native relationship validation SHALL remain outside internal/config. Historical aliases and removed layouts SHALL NOT be supported.

#### Scenario: Filesystem ownership

- **WHEN** notes, deck markers, and note types are loaded from the fixed hierarchy
- **THEN** their parsed wrappers retain unique relative paths without modifying sources or starting Anki

### Requirement: Identity-free resource configurations

Deck/note-type IDs and note GUIDs SHALL NOT be editable fields. Note-type/deck references SHALL use paths, preset references SHALL be readable names, and sort fields SHALL use names. Template order SHALL follow the array. Upstream native identity/provenance and derived metadata SHALL remain outside editable semantics: collection snapshot projection SHALL omit them, while pure Export SHALL preserve supplied protobuf values. Card ordinals and semantic field-role tags SHALL retain their meaningful behavior. config generate-guid SHALL be removed.

#### Scenario: New configuration

- **WHEN** a user writes a new note or deck file without an ID or GUID
- **THEN** it is inspectable as a path-addressed declaration without an identity-generation step

### Requirement: Pure configuration operations

config.Store.Load and Dump SHALL leave sources unchanged, preserve exact values and optional presence, and perform only syntax/schema/filesystem checks. Dump SHALL default to JSON and support TOML output. config.Store.Export SHALL directly write supplied path/config objects without normalization, identity allocation, validation, cleanup, or directory publication. ProtoJSON SHALL be converted directly into TOML using existing codecs. TOML datetimes SHALL require quoted string content.

#### Scenario: Direct export

- **WHEN** a caller exports a supplied configuration set
- **THEN** each resource is written at its declared relative path, the supplied object is unchanged, and unrelated files are retained

### Requirement: Explicit moves

The root config SHALL support moved entries with from/to paths. Config parsing SHALL preserve these entries without applying state transformations.

#### Scenario: Inspect a move

- **WHEN** config dump loads a directory-prefix or individual-file move
- **THEN** it prints the declared move without changing any source or requiring Anki

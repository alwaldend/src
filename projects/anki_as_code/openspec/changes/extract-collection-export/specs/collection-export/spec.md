## Purpose

Export prepared Anki collections as editable configurations without state output.

## ADDED Requirements

### Requirement: Borrow a prepared connection

Exporter SHALL borrow existing resource and state stores for the same prepared collection and own no native lifecycle. Export SHALL take only context and output directory. Export SHALL read native snapshots through the shared resource store without depending on reconciliation models.

#### Scenario: Reuse the caller's collection

- **WHEN** export completes or fails
- **THEN** the caller can continue using and close the same collection and backend

### Requirement: Read-only identity-free export

Export SHALL write exact note fields, tags, normal/filtered deck configs, and editable note-type fields/templates/CSS into the fixed configuration layout. It SHALL NOT generate, serialize, write or install state. It SHALL leave native objects, scheduling presets and review history unchanged. Native identity/provenance and derived fields SHALL be omitted from configs. Resident filtered cards SHALL remain unmanaged.

#### Scenario: Export a reviewed collection

- **WHEN** a collection containing notes, native presentation settings and review history is exported
- **THEN** config loading preserves editable values, no state file is generated and stored bindings remain unchanged

### Requirement: Stable readable paths

Export SHALL reserve tracked note/type filenames before allocating case-insensitive collision suffixes for new native objects. Initial note names SHALL derive from bounded safe sort-field text. Deck hierarchy SHALL map literal Anki title components into nested decks directories. Existing preset aliases SHALL be retained without editing presets.

#### Scenario: A native type reuses a bound filename

- **WHEN** a bound type is renamed and an unbound type uses its former name
- **THEN** the bound filename is preserved and the new type receives an unused filename referenced by its notes

#### Scenario: Derive a new note filename from rich content

- **WHEN** an unbound note's sort field contains inline markup, line breaks, entities or audio/media references
- **THEN** export uses Anki's HtmlToTextLine result, then applies bounded Unicode-preserving filesystem sanitization without modifying the field
- **AND** bound filenames remain independent of native text conversion

### Requirement: Write files directly

Export SHALL delegate config-file output to config.Export. It SHALL overwrite
declared configs in place and leave other files and directories untouched,
including stale declarations and caller-owned state files. It SHALL NOT stage,
clean, copy, replace or roll back the output tree. Filesystem errors SHALL be
returned and may leave partial output. Callers SHALL own directory cleanup,
synchronization and separate state operations.

#### Scenario: Export into an existing directory

- **WHEN** the destination contains old configs, state files and unrelated files/directories
- **THEN** export writes the requested paths and leaves all other entries unchanged

#### Scenario: A config write fails after earlier files were written

- **WHEN** a later config file cannot be written after earlier successful writes
- **THEN** export returns the filesystem error and retains the earlier written configs without rollback

## Why

Editable Anki configurations should identify objects by collection-relative paths. Native IDs and GUIDs belong in collection-backed state, so creation needs no identity generation and explicit moves preserve existing objects.

## What Changes

- Remove editable deck/note-type IDs and note GUIDs; remove config generate-guid.
- Retain source paths in parsed configuration wrappers and use readable path/preset/field references.
- Derive field/template positions from array order and exclude native identity/provenance metadata from exports.
- Support moved declarations in collection.anki.toml. Keep format version 1 without historical aliases.
- Keep parsing, dumping, and direct object export independent of Anki; collection-backed state, import, and reconciliation belong to #128.
- Split the configuration facade, filesystem layout, and protobuf/TOML/JSON codecs into focused Go packages within the existing module.

## Capabilities

### New Capabilities

- configuration-workflow: Load and serialize path-addressed configuration sets without native identity fields.

### Modified Capabilities

None.

## Impact

Use existing protobuf, TOML, Go, and Anki dependencies. Keep #151 and #128 independently against master; #128 consumes the configuration prerequisite. No host or scheduling-preset changes.

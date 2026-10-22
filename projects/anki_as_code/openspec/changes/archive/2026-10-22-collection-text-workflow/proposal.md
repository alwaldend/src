## Why

Editable Anki configurations should identify objects by collection-relative paths. Native IDs and GUIDs belong in collection-backed state, so creation needs no identity generation and explicit moves preserve existing objects.

## What Changes

- Remove editable deck/note-type IDs and note GUIDs; remove config generate-guid.
- Retain source paths in parsed configuration wrappers and use readable path/preset/field references.
- Derive field/template positions from array order and exclude native identity/provenance metadata from exports.
- Support moved declarations in collection.anki.toml. Keep format version 1 without historical aliases.
- Reuse the merged state operations through Anki configuration APIs, export only configs without changing the collection, and reconcile path identities through the existing native backend.
- Derive initial note filenames from sort-field text and retain recorded names on re-export.
- Remove the archive abstraction and call native package APIs in the collection
  workflow; reconciliation borrows connections without owning processes.

## Capabilities

### New Capabilities

- collection-text: Reconcile path-addressed configurations with persistent state and explicit moves.

### Modified Capabilities

None.

## Impact

Use existing protobuf, TOML, Go, and Anki dependencies. Build on the merged configuration prerequisite. No host or scheduling-preset changes.

## Follow-up package boundaries

Use focused implementation capabilities: app routing/backend composition, prepared-connection export/text publication, config files, shared codec,
reconciliation, and persisted state. Keep one
Go module. Consolidate reconciliation-only models and helper packages behind one
small API; preserve native ownership guards and direct database behavior.

## Merged state and codec prerequisite

[PR #158](https://github.com/alwaldend/src/pull/158) provides state.Store and
internal/codec on master. This change is rebased onto that merged prerequisite;
shared state and codec implementation changes are no longer part of its diff.

## Merged collection-workflow prerequisite

PR #163 provides native collection lifecycle and state CLI archive integration on master. Reuse WithCollection and construct stores from its prepared connection. The remaining feature diff adds reconciliation on top of merged export and pure planning.

## Public Go application boundary

Keep the concrete App in pkg/ankiascode, with path-based Config, Col and ColState
operations. internal/cli owns its private application interface and depends on
protobuf messages without importing App or native backends. The command entry
point constructs App and injects it into cli.New. Remove handlers, the separate
API package and collection reference objects. Reuse collection.Store's prepared
connection directly for scoped services and lifecycle.

Expose text-only config export for serialized CollectionConfigSet input, reusing the existing configuration store and codec without new dependencies.

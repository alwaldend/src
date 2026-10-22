## Why

The merged exporter and collection workflow can provide a useful export CLI independently of reconciliation. Extract that integration from #128 into a small master-based PR.

## What Changes

- Implement the existing export command for offline databases and collection packages using the merged native workflow and Exporter.
- Write TOML configuration files only; preserve source collections, stored bindings and caller-owned state files.
- Adapt existing native E2E cases for CLI export, direct writes, paths, error handling and selected state namespaces.
- Document the working command while plan, apply and build remain placeholders.
- Put command workflows behind the private CLI application interface and keep state store construction and native lifecycle inside App, with isolated command verification.

## Capabilities

### New Capabilities

- collection-export-cli: Configuration-only export through the managed collection CLI.

### Modified Capabilities

None.

## Impact

Only projects/anki_as_code changes. Reuse the merged exporter, resource/state stores, config APIs and backend without implementation changes or new dependencies. Keep #128's overlap until this independent prerequisite merges.

## Public Go application boundary

Keep the concrete App in pkg/ankiascode, with path-based Config, Col and ColState
operations. internal/cli owns its private application interface and depends on
protobuf messages without importing App or native backends. The command entry
point constructs App and injects it into cli.New. Remove handlers, the separate
API package and collection reference objects. Reuse collection.Store's prepared
Collection directly for scoped services and lifecycle.

Expose text-only config export for serialized CollectionConfigSet input, reusing the existing configuration store and codec without new dependencies.

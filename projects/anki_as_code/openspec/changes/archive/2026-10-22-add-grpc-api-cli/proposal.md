## Why

The CLI exposes collection workflows but cannot invoke arbitrary upstream Anki RPCs. A descriptor-driven command tree provides the whole API without maintaining another set of operation wrappers.

## What Changes

- Add `api <service> <method>` commands assembled at runtime from compiled Go descriptors, with JSON/TOML requests and responses.
- Select either `--collection` with the bundled stdio worker or `--backend` with a TCP/Unix gRPC connection; the options are exclusive.
- Add `--write-back` and `--no-write-back`, enabled by default with the last occurrence taking precedence.
- Reuse linked native hosting through the same application executable with server serve; retain the shared application and caller-owned source APIs.
- Delegate archive persistence to the collection backend regardless of request outcomes; read operations may repackage unchanged collection content.
- Use one `WithCollection` entry point with `CollectionOptions` containing the source path and write-back policy.

## Capabilities

### New Capabilities

- `grpc-cli`: Complete descriptor-driven API commands, explicit backend selection and deterministic bundling.

### Modified Capabilities

None.

## Impact

CLI, App construction, collection scopes, packaging, documentation and native E2E in this project. Paired changes in anki_grpc own connection lifecycle and write-back; infra/src records complete bindings within the existing Anki pin. No external dependency or live deployment changes.

## Purpose

Provide a portable WebAssembly command for text configuration workflows without
bundling either Anki server implementation.

## ADDED Requirements

### Requirement: Server-free WASI configuration

The repository SHALL provide `--//projects/anki_as_code:wasm=true` to build the public
`//projects/anki_as_code` command for `wasip1/wasm` with both bundled engines disabled.

#### Scenario: Build the WebAssembly command

- **WHEN** a user builds the public target with this configuration
- **THEN** the result is a WASI Preview 1 WebAssembly module
- **AND** neither the cgo adapter nor the embedded Rust WASM adapter is a dependency

#### Scenario: Mixed native and WebAssembly build

- **WHEN** the invocation builds the WASI CLI alongside another native target
- **THEN** only the selected CLI and its dependencies use `wasip1/wasm`
- **AND** sibling targets retain the invocation's normal platform and engine selections

#### Scenario: Server-engine flags are enabled

- **WHEN** the WASI CLI is selected while both `anki_grpc` engine flags are enabled
- **THEN** the selected CLI still excludes both bundled adapters
- **AND** other targets retain their requested engine flags

### Requirement: File-only command execution

The WASI command SHALL execute help and text configuration export/dump using
runtime-provided standard streams and filesystem access without requiring its
own executable path or a local server.

#### Scenario: Round-trip configuration

- **WHEN** a user exports a valid serialized configuration set and dumps its TOML tree
- **THEN** the resulting configuration set retains its declared content

#### Scenario: Invalid configuration input

- **WHEN** a user exports malformed serialized input
- **THEN** the command fails before publishing configuration files

#### Scenario: Request local server hosting

- **WHEN** a user invokes `server serve` in the WASI command
- **THEN** the command reports that no local engines are compiled

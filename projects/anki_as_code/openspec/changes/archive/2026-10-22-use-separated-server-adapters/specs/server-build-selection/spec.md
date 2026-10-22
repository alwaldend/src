## Purpose

Allow application builds to select native or WASM worker hosting while retaining
the established command behavior and supervised child-process boundary.

## ADDED Requirements

### Requirement: Explicit server build selection

The application SHALL provide native and cgo-disabled WASM builds through one
command package and public Bazel target. The independent `//projects/anki_grpc:cgo` and `//projects/anki_grpc:wasm` booleans SHALL include their engines explicitly. Cgo defaults to enabled and WASM defaults to disabled. Both generic subprocess and in-process selection SHALL prefer cgo and otherwise use WASM. Disabling both SHALL retain external-only access. Backend choice SHALL retain the host platform.
Guest building and embedding SHALL belong to the server adapter; application
commands SHALL invoke its argument-vector API without handling guest bytes.
Both SHALL launch their own executable in server mode for offline collection
work and SHALL retain existing process shutdown and collection semantics.

#### Scenario: Native build migration

- **WHEN** an existing native command opens a collection after the adapter package move
- **THEN** it uses the native adapter and retains its existing output and persistence behavior

#### Scenario: WASM worker

- **WHEN** a WASM application command opens an offline collection
- **THEN** it starts the same executable as a child and communicates with the WASM backend over gRPC stdio without native linkage when cgo is disabled

#### Scenario: Configured application build

- **WHEN** the same public application target is built with native-only and WASM-only engine settings
- **THEN** each executable contains only its enabled adapters, preserves the command name, and passes the existing CLI and collection workflows

#### Scenario: Composed optional engines

- **WHEN** both cgo and WASM engines are enabled
- **THEN** generic subprocess and in-process selection use cgo
- **AND** both qualified subprocess schemes are registered on the scoped factory and select their named engine for same-executable workers

#### Scenario: Independent worker compatibility

- **WHEN** a qualified subprocess address identifies an independently supplied executable such as standalone `anki_grpc serve`
- **THEN** the application preserves the encoded argument vector and leaves that executable's engine selection to its caller

#### Scenario: Engine selector reference

- **WHEN** the CLI reference is generated from the shared command tree
- **THEN** the server help includes the same root `--server-engine` flag as the installed application

#### Scenario: Explicit native worker in a combined build

- **WHEN** both engines are compiled and an operation selects `subprocess+cgo://`
- **THEN** the worker executes the native CLI even when WASM compilation is unavailable
- **AND** `subprocess+wasm://` still executes WASM and reports its own compilation failures without falling back

#### Scenario: External-only application

- **WHEN** both cgo and WASM engines are disabled
- **THEN** explicit external clients remain available without native linkage, while local server hosting and in-process schemes fail clearly

#### Scenario: Explicit in-process WASM selection

- **WHEN** a WASM-enabled application receives `inprocess+wasm://`
- **THEN** it hosts a fresh WASM instance in the application process over in-memory gRPC streams
- **AND** generic `inprocess://` prefers native when available and otherwise selects WASM, while qualified schemes retain explicit engine selection and the default remains the supervised worker

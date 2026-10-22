## Purpose

Provide scoped ownership of native Anki backends and collections while allowing callers to select the underlying transport independently.

## ADDED Requirements

### Requirement: Shared backend scope

The private withBackend helper SHALL start a supplied unstarted Backend, lend it to a callback and call Stop after every startup attempt, including callback failure, cancellation and panic unwinding. It SHALL reject a missing callback or already-cancelled context before startup. Managers SHALL supply the owned backend and cleanup timeout. Cleanup SHALL have the supplied timeout independently of caller cancellation while retaining context values. Both managers SHALL delegate their common lifecycle handling to this private function and preserve existing timeout policies.

#### Scenario: Compose scopes through existing managers

- **WHEN** a caller uses a transport or collection manager
- **THEN** both managers use the same private scope implementation without adding a public lifecycle function

### Requirement: Scoped transport ownership

The stdin/stdout manager SHALL create and start a fresh backend for each callback scope, lend it only while prepared, and stop it before returning on success, failure, cancellation or panic unwinding. Constructor options SHALL be copied. Independent scopes SHALL not share a worker. Invalid callbacks and already-cancelled contexts SHALL fail before startup.

#### Scenario: Callback finishes or fails

- **WHEN** a callback completes normally or returns an error
- **THEN** its worker is reaped before scope return, retained access is rejected, and the callback error remains discoverable

#### Scenario: Callback cancellation or panic

- **WHEN** a callback cancels its context or panics
- **THEN** bounded cleanup runs independently of cancellation and the cancellation cause or original panic remains observable

### Requirement: Transport-independent collection composition

A collection manager SHALL accept the transport manager interface and lend a prepared collection backend for a supplied path. It SHALL close that collection before the transport scope returns. Database access, package restoration, consuming exports and readiness SHALL retain the existing native backend contracts. Missing manager, empty path and invalid callback inputs SHALL fail without acquisition.

#### Scenario: Persist and reopen a collection

- **WHEN** a scope creates native content, closes, and another scope opens the database or its native archive
- **THEN** the native content remains available and completed temporary restores are removed

#### Scenario: Delegate through a supplied manager

- **WHEN** the collection manager uses a supplied transport manager
- **THEN** it obtains its backend through that manager and closes the collection while the transport remains available

### Requirement: Preserve scope and cleanup errors

Managers SHALL preserve underlying startup, callback, cancellation and cleanup error identities, including both callback and cleanup failures. Manager cancellation checks SHALL preserve both the standard context error and any custom cancellation cause. Cleanup SHALL use bounded contexts that retain context values without inheriting caller cancellation. Native uncertainty SHALL retain the existing restore-retention contract.

#### Scenario: Worker exits during a failed callback

- **WHEN** a callback fails and its worker also exits unexpectedly
- **THEN** both the callback cause and the reaped process-exit error remain discoverable

### Requirement: Keep application policy outside managers

Managers SHALL add no idle-worker cache, source-path reservation registry, archive publication policy or anki_as_code resource/state/export dependency. Existing low-level backend names and borrowed ownership SHALL remain supported.

#### Scenario: Read a collection package

- **WHEN** a package scope reads content without explicitly exporting through Anki
- **THEN** the supplied package remains byte-identical and the manager performs only native lifecycle operations

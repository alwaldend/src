# Backend access Specification

## Purpose

Provide explicit backend connection lifetime and consistent writable collection lifetime for offline applications.

## Requirements

### Requirement: Writable archive lifetime

Collection owners SHALL leave source archives unchanged by default. The explicit `WriteBack` constructor option SHALL enable publication with media on shutdown, independently of prior request outcomes. Disabling publication SHALL NOT prevent RPC mutations or change direct database persistence.

#### Scenario: Default collection scope

- **WHEN** a caller omits `WriteBack` and mutates a restored archive
- **THEN** shutdown leaves the source archive unchanged

#### Scenario: A request fails after earlier changes

- **WHEN** write-back is explicitly enabled, a collection contains committed changes and a later RPC fails
- **THEN** shutdown still publishes those changes when write-back is enabled

### Requirement: Recoverable publication failure

Publication failures SHALL be returned and retain restored collection files for recovery; operation errors SHALL remain discoverable alongside cleanup failures.

#### Scenario: Archive destination cannot be published

- **WHEN** shutdown cannot export the collection to its source
- **THEN** it reports the error and retains recoverable restored data

### Requirement: Borrow remote server lifetime

A network backend SHALL connect to an explicit TCP or Unix endpoint and close only its client connection during shutdown.

#### Scenario: Client disconnects

- **WHEN** a client stops its network backend
- **THEN** the remote process and collection remain available to another client

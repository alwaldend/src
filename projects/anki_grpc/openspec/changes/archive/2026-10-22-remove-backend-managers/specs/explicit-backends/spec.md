## Purpose

Provide reusable native process and collection owners with explicit lifecycle contracts while leaving application callback scopes to their consumers.

## ADDED Requirements

### Requirement: Explicit native ownership

The Go backend library SHALL expose process and collection owners through Start, Connection and Stop. The collection owner SHALL borrow a running process backend and SHALL NOT start or stop that process. Callers SHALL close the collection before stopping its process. The library SHALL provide no transport or collection manager API.

#### Scenario: Construct and close native owners

- **WHEN** a caller constructs a process and a collection backend
- **THEN** the caller starts them in dependency order, obtains the prepared connection and stops them in reverse order while retaining the existing error and restore contracts

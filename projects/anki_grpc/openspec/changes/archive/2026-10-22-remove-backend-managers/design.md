## Context

Managers create fresh owners and wrap callbacks; they do not pool backends. Existing BackendSTD and BackendCollection already implement explicit lifecycle and error reporting.

## Goals / Non-Goals

Remove the additional manager API while retaining reusable native owners. Worker protocol, native collection behavior and transport options are outside this change.

## Decisions

Delete manager declarations, constructors and their private helper. The application constructs and brackets the existing backend owners. Existing direct lifecycle E2E continue verifying startup failure, shutdown deadlines, crashes, collection preparation and restore retention. Scope callback/error coverage moves to the application E2E before implementation.

## Risks / Trade-offs

Removing a source API breaks callers using managers. A repository-wide Go call-site search identifies only the migrated store and tests; update those in the same PR and document explicit construction. No dependency is added.

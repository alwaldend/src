## Purpose

Provide collection-scoped native Anki connections over a caller-owned prepared
backend, with explicit readiness and storage-resource cleanup.

## ADDED Requirements

### Requirement: Borrow prepared backend lifetime

The collection backend SHALL implement the existing backend lifecycle interface
and accept its own options object. It SHALL require an already-started backend
and SHALL NOT start or stop it. Callers SHALL close the collection before stopping
the underlying backend. Each instance SHALL permit one Start attempt, reject
connections until prepared and after Stop, and make Stop idempotent.

#### Scenario: Close a collection and reuse its backend

- **WHEN** a prepared collection is stopped
- **THEN** its retained clients reject calls while the borrowed backend remains
  available to open another collection

### Requirement: Direct databases and native package restore

Database inputs SHALL open directly through Anki. Package inputs SHALL restore
through Anki's native package service. An omitted restore directory SHALL create
an owned temporary directory. An explicit destination SHALL require an empty
directory and SHALL remain on disk after Stop. Database/media paths SHALL follow
Anki's normal filesystem layout. No manual archive codec SHALL be introduced.

#### Scenario: Restore a package temporarily

- **WHEN** a package is opened without a restore-directory override
- **THEN** native methods access the restored collection and a clean Stop removes
  its temporary directory without modifying the source archive

#### Scenario: Retain an explicit restore

- **WHEN** a caller provides an empty restore directory
- **THEN** the restored database and media remain there after Stop

### Requirement: Respect native consuming operations and failures

The borrowed connection SHALL continue to use generated native service clients.
Completed native collection close or full-package export SHALL invalidate
collection access without causing Stop to close the collection twice, including
a package export that returns a native failure. Stop SHALL
wait for all active unary callers within its context while retaining ordinary-call
concurrency. Any unconfirmed unary response SHALL invalidate access and keep
cleanup uncertain, even when later calls succeed. Failed startup SHALL expose
no prepared connection and SHALL leave explicit Stop and backend lifetime to the caller. When native
completion is uncertain, cleanup SHALL retain potentially active restore files,
report the failure and retained restore location for cleanup after backend shutdown.
The public backend methods SHALL remain Start, Connection and Stop; callers
requiring programmatic cleanup SHALL choose an explicit destination through options.

#### Scenario: Export through a generated client

- **WHEN** a generated client successfully exports the native collection package
- **THEN** subsequent collection access fails, Stop does not close twice, the
  package can be reopened and the underlying backend remains running

#### Scenario: Fail without taking over the backend

- **WHEN** collection startup fails or Stop is cancelled
- **THEN** the collection reports its error and the caller retains responsibility
  for stopping its prepared backend and any required post-shutdown file cleanup

#### Scenario: Clean up after a native export error

- **WHEN** native package export consumes the collection and returns a completed
  native error
- **THEN** Stop removes the owned temporary restore without closing twice, and
  the borrowed backend remains available to open another collection

#### Scenario: Bound shutdown while an ordinary call is active

- **WHEN** a media download is active and its client is still waiting when
  collection Stop reaches its deadline waiting for active callers
- **THEN** Stop reports the deadline and retains the active restore, which the
  caller can safely remove after stopping and reaping its backend

#### Scenario: Retain files after an unconfirmed ordinary response

- **WHEN** an ordinary native RPC returns cancellation without a native completion
  response and later native work can still be queued or running
- **THEN** collection access fails, Stop reports uncertain cleanup and retains the
  restore until the caller stops and reaps its borrowed backend

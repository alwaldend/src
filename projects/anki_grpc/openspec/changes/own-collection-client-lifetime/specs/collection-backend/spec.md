## MODIFIED Requirements

### Requirement: Borrow prepared backend lifetime

A started Client SHALL provide OpenCollection with its own CollectionOptions
value and CloseCollection without a separate collection-client object.
OpenCollection SHALL require a ready transport and no existing collection phase.
CloseCollection SHALL publish and close storage while retaining that transport.
Each phase SHALL permit one opening attempt and reject collection access until
prepared and after close begins. A completed clean close SHALL be idempotent and
permit the same client to open another collection. Terminal Stop SHALL close any
remaining phase before stopping the transport.

#### Scenario: Close a collection and reuse its backend

- **WHEN** a prepared collection is closed through CloseCollection
- **THEN** retained collection connections reject calls while the same client remains
  available to open another collection
- **AND** Connection supplies raw transport access between collection phases

### Requirement: Direct databases and native package restore

Database inputs SHALL open directly through Anki. Package inputs SHALL restore
through Anki's native package service. An omitted restore directory SHALL create
an owned temporary directory. An explicit destination SHALL require an empty
directory and SHALL remain on disk after CloseCollection or Stop. Database/media
paths SHALL follow Anki's normal filesystem layout. No manual archive codec
SHALL be introduced.

#### Scenario: Restore a package temporarily

- **WHEN** a package is opened without a restore-directory override or enabled WriteBack
- **THEN** native methods access the restored collection and clean collection cleanup
  removes its temporary directory without modifying the source archive

#### Scenario: Retain an explicit restore

- **WHEN** a caller provides an empty restore directory
- **THEN** the restored database and media remain there after CloseCollection or Stop

### Requirement: Respect native consuming operations and failures

Connection SHALL continue to supply generated native service clients. Completed
native collection close or full-package export SHALL invalidate collection access
without causing cleanup to close the collection twice, including a package export
that returns a native failure. CloseCollection SHALL wait for active unary callers
within its context while retaining ordinary-call concurrency. Any unconfirmed
unary response SHALL invalidate access and keep cleanup uncertain, even when
later calls succeed. Failed opening SHALL expose no prepared connection and SHALL
leave explicit CloseCollection or Stop to the caller. When native completion is
uncertain, cleanup SHALL retain potentially active restore files, report failure
and the retained location for cleanup after transport shutdown. Callers requiring
programmatic file cleanup SHALL choose an explicit destination through options.

#### Scenario: Export through a generated client

- **WHEN** a generated client successfully exports the native collection package
- **THEN** subsequent collection access fails, CloseCollection does not close twice,
  the package can be reopened and the underlying transport remains running

#### Scenario: Fail without taking over the backend

- **WHEN** collection opening fails or CloseCollection is cancelled
- **THEN** the client reports its error and the caller retains responsibility for
  terminal Stop and any required file cleanup after Done

#### Scenario: Clean up after a native export error

- **WHEN** native package export consumes the collection and returns a completed
  native error
- **THEN** CloseCollection removes the owned temporary restore without closing twice,
  and the client remains available to open another collection

#### Scenario: Bound shutdown while an ordinary call is active

- **WHEN** a media download is active and its client is still waiting when
  CloseCollection reaches its deadline waiting for active callers
- **THEN** CloseCollection reports the deadline and retains the active restore,
  which the caller can safely remove after Stop and final Done completion

#### Scenario: Retain files after an unconfirmed ordinary response

- **WHEN** an ordinary native RPC returns cancellation without a native completion
  response and later native work can still be queued or running
- **THEN** collection access fails, CloseCollection reports uncertain cleanup and
  retains the restore until terminal transport completion

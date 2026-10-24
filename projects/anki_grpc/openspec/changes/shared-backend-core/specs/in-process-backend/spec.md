## MODIFIED Requirements

### Requirement: Native lifetime

The backend SHALL permit one start attempt and reject calls once stopping begins. Start SHALL run initialization synchronously and check its context before and after initialization. If cancellation or shutdown occurs during initialization, Start SHALL finish initialization and release any unpublished backend before returning. Command cancellation SHALL stop waiting without interrupting native work or freeing resources still in use. Stop SHALL release native state after initialization and active calls complete; its deadline SHALL bound waiting without pretending unfinished work was terminated. Repeated Stop calls SHALL be safe, and Done SHALL observe final cleanup.

#### Scenario: Cancellation during native work

- **WHEN** a native media request remains active after its caller cancels and Stop times out
- **THEN** retained clients reject calls, native resources remain valid, and later Stop completes after the native work finishes

#### Scenario: Shutdown overlaps native initialization

- **WHEN** native initialization is still running and Stop reaches its deadline
- **THEN** Stop returns the deadline error, Done waits for startup cleanup, and startup does not expose a connection after shutdown

#### Scenario: Cancellation during native initialization

- **WHEN** the startup context is cancelled while native initialization remains active
- **THEN** Start waits for initialization, releases the unpublished backend and returns a context error; concurrent Stop can time out while Done remains open until startup and cleanup finish

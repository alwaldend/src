## Purpose

Establish verified X identities and temporary cleanup permission without storing
credentials that could recreate access after logout, expiry, or backend restart.

## ADDED Requirements

### Requirement: OAuth 2.0 identity and separate cleanup permission

The service SHALL use authorization code flow with PKCE S256, single-use bound
state, verified scopes, and a fixed callback origin. Login SHALL establish a stable
X user ID without granting cleanup write operations. Cleanup authorization SHALL
obtain read/write user permission and verify that its X identity matches the login.
Both flows SHALL request offline access for on-demand renewal during active use.
Application client credentials SHALL come from deployment environment injection.

#### Scenario: User switches X accounts during cleanup authorization

- **WHEN** the cleanup token identifies a different X account than the logged-in account
- **THEN** the service rejects attachment of that grant and performs no cleanup operation

#### Scenario: Replayed or mismatched callback

- **WHEN** an authorization callback has expired, reused, or mismatched state
- **THEN** the service rejects it before creating a session or attaching credentials

### Requirement: Memory-only user credentials

X access tokens, refresh tokens, PKCE verifiers, and equivalent reusable secrets
SHALL exist only in backend memory. They SHALL NOT enter the database, cookies,
browser storage, worker environment, logs, backups, swap, or core dumps. Persistent
application sessions SHALL NOT be sufficient to recreate X authorization.

#### Scenario: Backend restarts with persistent account data

- **WHEN** a user returns after the backend restarts
- **THEN** their valid application session can access retained account data
- **AND** fresh X requests require new authorization

### Requirement: Safe refresh within the activity window

The service SHALL renew access tokens on demand using its memory-held refresh
token, serialize renewal per grant, and adopt rotated credentials atomically.
Refresh failure SHALL stop X dispatch and require reauthorization. Refresh SHALL
NOT change the grant's verified account binding or extend the inactivity deadline.

#### Scenario: Concurrent requests encounter access-token expiry

- **WHEN** multiple admitted operations need renewal for the same grant
- **THEN** renewal is serialized and all later use refers to the accepted current token pair
- **AND** an uncertain failed exchange is not blindly replayed

### Requirement: Twenty-four-hour inactivity expiry

The service SHALL discard X authorization after 24 hours without intentional
authenticated user activity. Polling, worker progress, and token refresh SHALL NOT
reset that deadline. Expired grants SHALL NOT be revived by a subsequent visit.
Disconnect, logout, session expiry, and backend shutdown SHALL also discard local
authorization, cancel active X work, and stop new dispatch. Revocation SHALL be
attempted when possible without retaining credentials indefinitely on failure.

#### Scenario: User returns before inactivity expiry

- **WHEN** the access token has expired but the user returns within the activity window
- **THEN** the service can refresh the memory-held grant without requesting authorization again

#### Scenario: Unattended job and polling pass the deadline

- **WHEN** 24 hours pass without intentional activity while a job or browser poll continues
- **THEN** the grant is discarded and no further X request is dispatched
- **AND** retained snapshots remain available for local evaluation

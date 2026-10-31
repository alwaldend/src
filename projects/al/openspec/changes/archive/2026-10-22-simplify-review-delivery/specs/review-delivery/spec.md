## Purpose

Make validated pull-request review delivery reliable and concise while preserving exact-candidate and mutation consistency checks.

## ADDED Requirements

### Requirement: Guarded reply and resolution

The CLI SHALL offer a single review-address operation using a published validation receipt, a thread ID, and a reasoned reply file. It SHALL derive current expectations, post at most one reply per durable reply receipt, and resolve only the verified resulting thread. A retry SHALL resume an existing valid reply receipt rather than posting another reply.

#### Scenario: Address a finding

- **WHEN** the published validated head and unresolved thread match the supplied candidate receipt
- **THEN** the operation posts the reasoned reply, verifies it, resolves the unchanged thread, and returns the resulting inventory

#### Scenario: Candidate or thread changes

- **WHEN** the published head differs or the target thread changes during processing
- **THEN** the operation refuses mutation or stops further mutation with diagnostic state

### Requirement: Pre-mutation read recovery

The CLI SHALL preserve the original unexpired reply receipt when a transient provider read or advancing review epoch prevents observing a coherent inventory before any resolution mutation. It SHALL retain single-use consumption on attempted mutations, semantic mismatches, and expiry, without extending authority.

#### Scenario: Advancing epoch before resolution

- **WHEN** read-only inventory collection observes an advancing GitHub epoch before resolution
- **THEN** the operation reports that resolution was not attempted and leaves the original unexpired receipt available for a guarded retry

#### Scenario: Unknown mutation result

- **WHEN** a resolution mutation was attempted but its result cannot be verified
- **THEN** the receipt remains consumed and the tool requires fresh inspection

### Requirement: Bounded complete inventory

Review inspection SHALL include every returned thread and comment while fetching the first comment page with the thread page. It SHALL fetch additional comment pages when necessary and retain request, byte, and node limits, identity checks, and duplicate detection.

#### Scenario: Many ordinary threads

- **WHEN** every returned thread has one page of comments
- **THEN** inspection obtains their comments without one separate details request per thread

#### Scenario: Comment pagination

- **WHEN** a thread has more comments than its embedded page
- **THEN** inspection follows its cursor and returns all comments or refuses incomplete state

## MODIFIED Requirements

### Requirement: Reporting and failure handling

The publisher SHALL report the identifier it obtained and SHALL fail loudly
rather than report success when a request was rejected. Failures MUST identify
which operation failed and MUST NOT be retried automatically.

For every non-2xx API response, including HTTP 503, the error SHALL preserve and
display the original operation, status, body, and only these response headers:
`X-Rate-Limit-Limit`, `X-Rate-Limit-Remaining`, `X-Rate-Limit-Reset`,
`X-User-Limit-24hour-Limit`, `X-User-Limit-24hour-Remaining`,
`X-User-Limit-24hour-Reset`, and `Retry-After`. It SHALL preserve every value of
those headers and SHALL NOT copy authorization, cookies, or other headers into
the error. A rejected response whose body cannot be completely read SHALL
retain its status and headers alongside the body-read error cause.

Valid Unix reset timestamps SHALL also be displayed as human-readable UTC.
A combined rate-limit retry not-before boundary SHALL be reported only when
known exhausted windows have valid, future reset times at response receipt.
A window SHALL count as exhausted only when its unambiguous remaining value
is zero. The boundary SHALL use the latest reset among those windows, extended
by a valid later `Retry-After` value interpreted as decimal seconds from receipt
or an HTTP date. Missing, malformed, ambiguous, or already-elapsed reset data
for an exhausted window SHALL prevent a combined boundary. A positive remaining
count or `Retry-After` alone MUST NOT imply upload eligibility. The diagnostic
SHALL distinguish a not-before boundary from guaranteed service recovery or
request success, and SHALL preserve raw values when timing is unknown.

#### Scenario: Report the created identifier

- **WHEN** a draft is created
- **THEN** the publisher reports the draft's article identifier

#### Scenario: The API rejects a request

- **WHEN** an API request is rejected
- **THEN** the publisher reports the failing operation and the API's error
- **AND** it does not report the operation as successful or retry automatically

#### Scenario: Service unavailable with remaining allowance

- **WHEN** a 503 response includes remaining allowance and reset headers
- **THEN** the error retains those headers and displays readable UTC resets
- **AND** it does not claim that the resets guarantee service recovery

#### Scenario: Both rate windows are exhausted

- **WHEN** both remaining counts are zero and both resets are valid future times
- **THEN** the rate-limit boundary is no earlier than the later reset
- **AND** valid later Retry-After advice extends that boundary

#### Scenario: Timing evidence is incomplete or invalid

- **WHEN** an exhausted window has a missing, malformed, ambiguous, or elapsed reset
- **THEN** the diagnostic reports that a current retry time is unknown
- **AND** it keeps the original header values without inventing an exact time

#### Scenario: A rejected response body is truncated

- **WHEN** a non-2xx response terminates before its declared body length
- **THEN** the error retains the operation, status, partial body, and allowlisted headers
- **AND** callers can still inspect the body-read error cause

#### Scenario: Validation runs without network access

- **WHEN** the project's checks validate conversion
- **THEN** they run without contacting the X API
- **AND** no test requires live credentials

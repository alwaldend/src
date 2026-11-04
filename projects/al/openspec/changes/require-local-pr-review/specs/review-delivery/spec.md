## ADDED Requirements

### Requirement: Local review before each pull-request publication

The delivery procedure SHALL require a local code review using the
`review-agent` workflow before each pull-request creation or update. The review
SHALL cover the complete merge-base-to-prepared-candidate diff, including
documentation and configuration changes. The caller SHALL evaluate all findings,
fix in-scope defects, and review the corrected candidate before publication.
Passing checks or a remote review SHALL NOT substitute for local review.

#### Scenario: First publication

- **WHEN** a prepared candidate has passed its required checks
- **THEN** the caller completes local review before creating its pull request

#### Scenario: Candidate changes during delivery

- **WHEN** a review fix, rebase, or other edit changes the reviewed candidate diff
- **THEN** the caller reviews the updated candidate before publishing it

### Requirement: Remote feedback handling after each publication

The delivery procedure SHALL require remote feedback inspection after each
pull-request publication. The caller SHALL fix valid in-scope feedback, locally
review and validate the resulting candidate, publish it, and reply with the
reasoned change before resolving its thread. Inapplicable feedback SHALL receive
a reasoned reply before resolution. Reviews running for the final head SHALL be
observed through a supported mechanism; unavailable results or unsupported
resolution SHALL be reported as incomplete rather than successful.

#### Scenario: Remote finding requires a fix

- **WHEN** remote feedback identifies a valid in-scope defect
- **THEN** the caller corrects, checks, locally reviews, and publishes the new
  candidate before replying and resolving the finding with guarded operations
- **AND** the caller inspects feedback again for the updated publication

#### Scenario: Review cannot be verified

- **WHEN** the available provider cannot observe review completion or resolve
  a feedback thread
- **THEN** the caller reports the remaining review state and capability limit
  without claiming full review completion

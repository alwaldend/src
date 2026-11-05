## Purpose

Allow an interrupted source publication to synchronize its exact pull-request
projection without losing later local work or overwriting independent edits.

## ADDED Requirements

### Requirement: Exact post-push head transition

Publication SHALL retain the observed pre-push PR expectation. An update SHALL
accept only its original head or the exact pushed candidate, with all other
identity and editable metadata unchanged. A failed PR update after a successful
push SHALL retain the known candidate, remote head, PR URL, and pending status.

#### Scenario: PR head lags behind a completed push

- **WHEN** Git reports the candidate pushed but the unchanged PR still reports its prior head
- **THEN** the tool synchronizes the intended projection without treating that exact head transition as a human edit

#### Scenario: Independent edit during publication

- **WHEN** PR identity, branch topology, title, body, author, draft status, or head changes outside the expected transition
- **THEN** the tool refuses the metadata mutation and retains the completed push evidence

### Requirement: Receipt-bound PR recovery

The CLI SHALL provide explicit PR-only recovery from the original trusted
preparation receipt. It SHALL require the same repository, transport, endpoints,
refs, recorded PR identity, and exact already-pushed candidate. It SHALL verify
the prepared commit and scope. It SHALL reject unexpected remote changes,
replacement PRs, and metadata outside the original or desired projection.

#### Scenario: Local work advanced after a partial push

- **WHEN** the original candidate is still the remote feature head and its PR projection is interrupted, while local commits or working changes have advanced
- **THEN** explicit recovery synchronizes only that PR, preserving local HEAD, refs, index, working files, and the original receipt

#### Scenario: Recovery is already complete

- **WHEN** the PR already matches the receipt's candidate and desired projection
- **THEN** recovery succeeds without another metadata mutation

#### Scenario: Unknown or changed publication state

- **WHEN** the remote head is not the exact prepared candidate, the receipt binding differs, the PR is replaced, or human-edited metadata is observed
- **THEN** recovery refuses without pushing, recreating the PR, or changing local source

### Requirement: Recovery reports its limited scope

Recovery SHALL distinguish synchronized metadata from complete publication
verification. It SHALL retain known commit and PR identities if an attempted
metadata mutation cannot be verified. It SHALL NOT claim checks, remote review,
or validation of later local source.

#### Scenario: Provider read fails after recovery mutation

- **WHEN** the metadata update is attempted but its final read fails
- **THEN** the tool returns pending recovery evidence and requires inspection before any retry

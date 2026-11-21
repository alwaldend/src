## Repository scan verification

The pinned TruffleHog 3.99.7 history and checkout tests passed with credential
verification disabled. The checkout had 72 materialized Git LFS files.

The reviewed text matches have inline source annotations. Both scans
reported zero unmatched findings.

A disposable Git LFS fixture produced a finding in its tracked materialized
payload. The scan excluded an untracked fixture file and rejected the checkout
when the tracked payload was replaced with its pointer.

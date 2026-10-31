## 1. Delivery behavior

- [x] 1.1 Extend the existing integration fixture before implementation to exercise combined publication, failed checks, and review-read failure with verification-only recovery.
- [x] 1.2 Add combined validation/publication and partial publication reporting; verify the delivery integration suite passes without weakening refusal cases.
- [x] 1.3 Document combined execution and pending verification; verify the CLI help and a representative structured result.

## 2. Delivery

- [x] 2.1 Prepare the reviewed candidate for guarded delivery after formatting, owner-local OpenSpec validation, package tests, repository quality, and semantic lint. Exact-commit validation and publication are recorded in the delivery receipt.

## Evidence

The new isolated Git cases passed before the full selected suite (32 tests).
They cover passing combined execution, failed checks, retained publication
identity after a review epoch failure, and verification-only recovery.
Affected-package semantic lint and strict OpenSpec validation passed. CLI help
exposes the explicit publication flags. The final delivery receipt under
`out/agent-ergonomics/` binds publication checks to the exact commit.

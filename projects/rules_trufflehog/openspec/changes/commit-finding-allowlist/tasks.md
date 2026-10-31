## 1. Contract and fixture

- [x] 1.1 Define and write subprocess E2E failure cases before wrapper implementation; retain an execution report artifact.
- [x] 1.2 Establish the standalone module and repository wiring using existing pins; verify dependency reconciliation and module discovery.

## 2. Implementation

- [x] 2.1 Implement strict commit-and-line matching, scan-only execution, preserved original output, and error handling; pass the E2E command fixture.
- [x] 2.2 Expose and document the public Bazel test rule; build and test the project and its rule fixture.
- [x] 2.3 Migrate all existing path exclusions into exact finding entries and update both repository scan targets and the skill; verify the real scan reports ten approved and zero unmatched findings.

## 3. Delivery

- [x] 3.1 Run independent module checks and the complete workspace matrix; diagnose failures and complete mandatory quality and affected-target lint checks against the prepared candidate before updating the existing PR.

## 4. Upstream metadata types

- [x] 4.1 Replace the handwritten Git metadata projection with TruffleHog's generated protobuf type at the scanner's pinned version; reconcile Go and Bazel metadata.
- [x] 4.2 Validate the existing E2E failure cases, independent module, and real scan; complete exact-candidate publication checks and update the PR.

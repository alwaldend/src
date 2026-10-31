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

## 5. Scanner ownership

- [x] 5.1 Define acceptance with scanner omitted from repository consumers and the module-owned CLI included in standalone build checks; confirm the current API fails that contract.
- [x] 5.2 Move the existing scanner release lock and toolchain into the module, supply it by default, retain explicit overrides, and make repository CLI aliases forward to the module.
- [x] 5.3 Validate standalone builds/tests, default and overridden scanner selection, generated module state, and a fresh repository scan before publishing the checked candidate.

## 6. Generic registration and root cleanup

- [x] 6.1 Register all scanner toolchains generated from the module-owned lock and remove platform-specific registration aliases and parent registration.
- [x] 6.2 Remove the obsolete root third-party package and update documentation for the remaining repository CLI alias.
- [x] 6.3 Validate default toolchain resolution from both workspaces, existing override behavior, repository documentation discovery, and a fresh scan before updating the PR.

## 7. Scan state-machine documentation

- [x] 7.1 Add one canonical Mermaid state machine and place its rendered diagram first in the README and project landing; register the landing through the existing external-module registry.
- [x] 7.2 Render the shared site light/dark SVGs and a dark WebP, inspect the diagram, and capture an image of the actual repository allowlist.
- [x] 7.3 Validate standalone documentation, main-site output and image references, and required publication gates before updating the PR.
- [x] 7.4 Give TruffleHog a nested state machine with start and end markers, simplify the wording, distinguish scanner and wrapper outcomes, and verify both documentation surfaces.
- [x] 7.5 Replace paragraph-sized state labels with short actions and a linked documentation key; inspect the dark render and both site pages before publication.
- [x] 7.6 Model scanner execution and stdout processing concurrently, join only after EOF and process exit, and verify repeated findings and both completion orders before publication.

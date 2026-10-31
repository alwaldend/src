## 1. Dependency update

- [x] 1.1 Verify the latest stable publisher release, update the existing archive pin, and confirm the pinned executable reports version 3.99.7.
- [x] 1.2 Run the existing repository scan with the updated binary, retain original local findings, and record completion and counts separately from the test exit status.
- [x] 1.3 Add exact path exclusions for the four user-confirmed false positives and verify a completed scan with those exclusions.
- [x] 1.4 Disable local test-result reuse for both history targets and verify a fresh scan of the committed candidate.

## 2. Validation

- [x] 2.1 Reconcile generated module state and verify repository quality and affected-target lint before delivery.

## 3. Shared Go build-action sandboxing

- [x] 3.1 Add the verified execution-metadata setting to the shared Bazel configuration and document its scope; verify the default action graph removes `no-sandbox` from Go build actions and retains it on the local test runner.
- [x] 3.2 Exercise actual Go compilation with the default configuration, run the environment-policy test, and validate the final candidate through repository quality, affected-target lint, skill checks, and a fresh secret scan before updating the existing PR.

## 4. Current synthetic fixtures

- [x] 4.1 Add native same-line ignore comments with reasons to current XO-origin and delivery-runner synthetic credential fixtures, preserving historical commit-and-line exclusions; verify the pinned scanner honors the comments and can inspect ignored lines explicitly.

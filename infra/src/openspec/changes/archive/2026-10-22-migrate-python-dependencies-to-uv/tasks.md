## 1. Migrate tooling

- [x] 1.1 Pin the uv toolchain and migrate local lock targets; verify generated module freshness and target analysis.
- [x] 1.2 Patch firmware's active locking rule to uv while preserving upstream pins and lazy hub loading; verify patch application and target analysis.
- [x] 1.3 Update workflow documentation and Python skill guidance; verify skill validation and offline eval configuration.
- [x] 1.4 Migrate the remaining K3s bootstrap installer and Autoscroll install hint; verify role packaging, offline command rendering and installation idempotency in a disposable local environment.

## 2. Validate and deliver

- [x] 2.1 Regenerate all four local locks through Bazel-managed uv, inspect pin changes and refresh Gazelle's manifest through its owning workflow.
- [x] 2.2 Run explicit lock comparisons, representative Python consumers and repository quality checks; record repeatable artifacts and any external blockers.
- [x] 2.3 Record acceptance evidence and the guarded delivery approach; verify OpenSpec structural validation. Publication state is owned by Git and delivery receipts.

## 3. Simplify the locking API

- [x] 3.1 Use upstream uv lock directly, remove optional comparison tests and unused shared helpers, and eliminate cross-component internal loads and their documentation patch.
- [x] 3.2 Update workflow documentation and skill guidance, reconcile generated module metadata, and validate lock output and affected consumers before publishing the revised candidate.

## 4. Adopt native uv locks

- [x] 4.1 Upgrade rules_python and its Gazelle plugin to the current stable release, inspect their native lock APIs, and update the selected Python runtime to the latest stable runtime bundled with that rules release, without a custom archive.
- [x] 4.2 Replace local requirements manifests and locks with component-owned pyproject.toml and uv.lock files, wire Bazel and Gazelle directly to native locks, and migrate the firmware integration.
- [x] 4.3 Update documentation and Python skill guidance, inspect package-version changes, validate native locking and representative consumers, and publish the exact validated candidate.

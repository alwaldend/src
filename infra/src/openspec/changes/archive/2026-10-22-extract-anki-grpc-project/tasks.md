## 1. Extract source and ownership

- [x] 1.1 Copy projects/anki_grpc and its required native build definitions from the exact #146 snapshot; verify source and Cargo lock equality and exclude Go consumer and diagram files.
- [x] 1.2 Update project verification documentation for the separate Go consumer; verify every listed build target exists in this PR.

## 2. Verify and publish

- [x] 2.1 Build the independent native library, server, schemas, docs, and Clippy; inspect the upstream closure and existing Go E2E artifact for the unchanged server implementation.
- [x] 2.2 Format and validate OpenSpec and select the exact-candidate repository quality, semantic lint, and build gates. Delivery receipts own publication against master and the prerequisite relationship in #146.

## 3. Configurable transports and native ownership

- [x] 3.1 Write CLI and TCP/Unix native persistence/reconnect E2E cases before implementation and verify they fail against the stdio-only binary.
- [x] 3.2 Add CLI transport options, help, validation, TCP and Unix listeners with startup initialization and a shared backend; verify E2E artifacts and default Go stdio compatibility.
- [x] 3.3 Clarify that Anki's source/build definitions belong to third_party and the gRPC project depends on its native library; verify native closure and unchanged dependency pins.
- [x] 3.4 Update both PRs' maintained work and validation plans, preserving their master bases; delivery receipts own exact-candidate checks and publication.
- [x] 3.5 Address the landing-page review by registering the project, adding owned landing content, and verifying its built index, landing, docs, and taxonomy pages.

## 4. Direct Cargo workspace import

- [x] 4.1 Inspect upstream manifests, generated helper coverage, and the original approved package inventory; prepare the exact workspace-only dependency proposal.
- [x] 4.2 Obtain approval for the workspace-only locked additions before declaring or acquiring them.
- [x] 4.3 Import the unmodified pinned Cargo workspace, use its standard dependency helpers, and remove crate_dependencies and existing_rule reflection.
- [x] 4.4 Verify native/server builds, transport E2E artifacts, unchanged Go stdio behavior, and dependency boundaries; update docs and both PR snapshots.
- [x] 4.5 Format and select the exact-candidate delivery gates for the direct workspace import; delivery receipts own check execution and publication.
- [x] 4.6 Wire workspace-enabled native TLS to the existing pinned OpenSSL build, verify its C consumer, and mark anki_grpc finished in its README and landing page.

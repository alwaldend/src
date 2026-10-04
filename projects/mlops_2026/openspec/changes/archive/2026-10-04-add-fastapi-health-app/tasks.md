## 1. Define behavior verification before implementation

- [x] 1.1 Write the E2E failure coverage and harness first, covering startup,
      readiness, health content, OpenAPI, docs URLs, independent routing, port
      conflict, shutdown, missing data, incorrect runfile resolution, training failure, and invalid model output; verify it fails meaningfully against the placeholder.

## 2. Package the course data

- [x] 2.1 Add the user-approved source repository under
      third_party/com_github_yandex_practicum_mlops_freetrack with immutable commit,
      integrity, narrowly visible CSV targets, and provenance; verify archive
      resolution and that data.csv and dataset.csv match the pinned upstream bytes.

- [x] 2.2 Add upstream Python requirements to the shared manifest, regenerate
      locks and Gazelle metadata, and verify dependency and manifest tests.
- [x] 2.3 Package the upstream training script as a Python binary, add a
      sandboxed model build rule, and verify that the declared model.pkl loads
      and provides a valid recommendation without source-tree writes.

## 3. Implement the local application

- [x] 3.1 Add the agreed cmd/recommender, internal/recommender/application, internal/recommender/http, and
      internal/recommender/http/v1 packages with mandatory annotations and fresh app/router
      factories; verify factory isolation through the HTTP E2E scenario.
- [x] 3.2 Implement GET /health with a typed response model and preserve default
      OpenAPI/documentation endpoints; verify health JSON, schema contract, docs
      responses, and absent prediction routes through E2E HTTP checks.
- [x] 3.3 Add Uvicorn startup with factory mode and configurable host/port using
      existing locked dependencies; generate and inspect BUILD wiring, then verify
      the documented Bazel executable serves HTTP on an explicitly selected port.

- [x] 3.4 Add both CSV targets to the server runfiles and use existing runfiles
      support for path resolution; verify both can be opened from the actual app
      runfiles outside the repository directory and record checksums in E2E output.

## 4. Documentation and final acceptance

- [x] 4.1 Update the README and OpenSpec context to describe the implemented
      service and exact run command; verify documentation builds and preserve the
      course link and absence of a release deployment requirement.
- [x] 4.2 Run the complete project E2E check, retain its JSON result artifact,
      and verify clean child-process shutdown; run project semantic lint and select the required exact-candidate
      repository quality checks for receipt-bound delivery.
- [x] 4.3 Verify implementation against the change artifacts and prepare accepted
      deltas for archive and receipt-bound publication. Final PR-head verification
      and exact-candidate validation remain owned by repo-delivery receipts.

## 5. Recommendation refinement

- [x] 5.1 Define HTTP failure coverage first: missing query, invalid n, known and
      unknown tracks, missing model startup failure, and response/schema shape.
- [x] 5.2 Add namespaced shared response models and a per-app model lifecycle,
      calling upstream load_model with a packaged path via a compatibility patch.
- [x] 5.3 Serve /api/v1/recommend with validated query count and typed response.
- [x] 5.4 Add and inspect the README Mermaid architecture render; update docs
      and Python skill import requirements and offline skill coverage.
- [x] 5.5 Validate refined HTTP/model E2E, dependency contracts, documentation,
      skill packaging, and semantic lint; select exact-candidate publication gates.

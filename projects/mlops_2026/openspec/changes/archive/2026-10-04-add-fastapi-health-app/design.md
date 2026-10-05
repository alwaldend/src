## Context and decisions

The project began as an MLOps course placeholder. The user requested a FastAPI
recommendation app served by Uvicorn, a buildable upstream trainer and model,
packaged CSVs, and namespaced imports. The existing shared Python toolchain and
pinned Mermaid and ast-grep tools satisfy the request without new tooling.

## Layout

- cmd/recommender/main.py parses host and port and starts Uvicorn in factory mode.
- internal/models/recommender.py owns shared Pydantic response models.
- internal/recommender/application assembles the app and resolves runfiles.
- internal/recommender/service owns a model reference for one application.
- internal/recommender/http owns health and versioned recommendation transport.
- internal/recommender/training owns the build action wrapper and rule.
- test/e2e exercises actual processes, HTTP, and the trained model.

Future apps use sibling internal namespaces and separate cmd entries. No
language-first package layer or unnecessary initializer files are required.

## Application lifecycle and HTTP contract

create_app constructs a fresh FastAPI, fresh routers, and a fresh recommendation
service. Its lifespan calls upstream load_model once, using the model path
resolved through runfiles, and releases the reference on shutdown. Imports are
inert; no app, router, model, or service instances exist at module scope.

GET /health returns status ok without computing recommendations. Default
/openapi.json, /docs, and /redoc remain enabled. GET /api/v1/recommend takes a
required track_name and n defaulting to 5, bounded to 1–20. A sync handler lets
FastAPI run recommendation work in its threadpool. Shared response types expose
requested_track and songs with track_name, artists, album_name. Unknown tracks
retain upstream empty-result behavior; invalid queries receive 422. Prediction
and deployment remain outside scope.

## Upstream data, training, and loader compatibility

Pin a08aadd41a411d2c8b5e897dc9d423ce97d1fa42 through an integrity-verified archive
under third_party/com_github_yandex_practicum_mlops_freetrack. CSV targets and
training code are visible only to the project. Both CSVs and model.pkl are app
runfiles, resolved independently of cwd. Specify the main source repository explicitly
for Rlocation so manifest resolution does not rely on caller-path inference. Upstream load_model opens a relative
filename, so the user selected a minimal optional-path patch that preserves its
model.pkl default instead of changing the process working directory.

The training script remains unchanged. A sandboxed ctx.actions.run action stages
its declared CSV inputs in a temporary action directory, invokes the upstream
binary, verifies its output, and moves it to declared model.pkl. Clear inherited
runfiles environment variables so the child discovers its own runtime tree.
Fix numerical thread counts at one. No network or source-tree writes are needed.
The shared manifest, lock, and Gazelle mappings use their existing generators.

The full data joins to 6,788 tracks, with a 0.34 GiB similarity matrix and an
approximately 359 MiB serialized artifact. The explicit model and isolated model
verification targets are manual, but the server consumes the model, so building
or testing it triggers training when the artifact is uncached. No upstream
LICENSE file was found; provenance docs do not imply a first-party license.

## Verification and maintained documentation

The HTTP harness and failure coverage were written before implementation and
failed meaningfully against the placeholder. The recommendation refinements
were similarly covered before their implementation. Exercise health, OpenAPI,
docs, known/unknown tracks, count validation, independent routing, conflicting
ports, missing model startup failure, graceful shutdown, and CSV checksums from
an external cwd. The model consumer loads the full artifact and recommends a
real track. Tests retain repeatable JSON artifacts in Bazel undeclared outputs.

The README consumes the maintained Mermaid light/dark SVG targets. The raster
preview of the same source was inspected: both pipelines, container titles,
load_model connection, and labels are readable without edge/title collisions.

Namespaced import policy is owned by repo-python and enforced with the pinned
ast-grep syntax rule in the repository quality suite. Ruff ICN003 only supports
explicitly named modules, so it cannot enforce a universal ban. Syntax fixtures
reject normal, aliased, multiline, relative, and star imports while accepting
module imports, comments, and strings. Existing imports were migrated while
preserving Autoscroll exports. Shared changes and offline skill coverage are
tracked by the companion infra/src record.

## Risks and acceptance

The available-port probe has a bind race; the harness detects startup failure
and reaps children. Docs UIs retain FastAPI's external CDN assets; tests verify
served HTML/schema references, not external asset availability. A missing model
must fail startup rather than silently creating a health-only app. The model is
loaded from the pinned build output, never an arbitrary client-supplied pickle.

Exact-candidate quality, semantic lint, project/spec/archive checks, shared
consistency checks, and affected consumers are mandatory delivery gates. Source
acceptance precedes archive; publication and final PR-head verification belong
to trusted repo-delivery receipts in ignored out/mlops-python.

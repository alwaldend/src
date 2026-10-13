## 1. Contract and behavioral cases

- [x] 1.1 Record supported platform, approved dependency closure, build ownership, and source immutability; verify the OpenSpec artifacts and dependency inventory.
- [x] 1.2 Add the synthetic collection E2E before the API implementation; cover persistence, generated cards, native errors, cancelled contexts, unknown commands, and closed handles, producing a database and JSON report.

## 2. Cargo dependency and API

- [x] 2.1 Replace the custom Cargo action with crate_universe, standard upstream Rust targets, and a rust_static_library bridge; verify the dependency versions and absence of networked compilation actions.
- [x] 2.2 Add the consumer-owned native bridge and documented Go API; verify the synthetic E2E and upstream-schema consumers.
- [x] 2.3 Document supported prerequisites, cache/build limitations, API boundaries, and deferred PR #128 integration; verify README matches actual behavior.

## 3. Delivery

- [x] 3.1 Format and validate the exact candidate with repository quality, affected semantic lint, and native/API E2E; inspect the repeatable database artifact.
      Publication and final validation are tracked by the repo-delivery receipts for this branch, independently of this implementation checklist.

## 4. Simplification review

- [x] 4.1 Review the complete PR and simplify redundant linking, unused dependencies/configuration, generated repository exposure, compile inputs, and native response ownership while preserving the upstream dependency versions and Go API.
- [x] 4.2 Validate the revised native/Go consumer, inspect its database artifact and action graph, complete quality and semantic lint gates, and update PR #146.

## 5. Dependency ownership

- [x] 5.1 Make the third-party Anki library independent of project source and build rules; move the bridge targets into the project and the generic dependency-provider rule into shared Rust build support.
- [x] 5.2 Verify the third-party dependency graph contains no project inputs, preserve the approved external dependency inventory, and validate the independent library and Go E2E before updating PR #146.

## 6. Single Cargo lock

- [x] 6.1 Remove the separate Cargo.Bazel.lock and its references; generate Bazel metadata from the pinned Cargo.lock and document the normal module-cache behavior.
- [x] 6.2 Verify metadata regeneration, unchanged Cargo dependency pins, independent native build, and Go E2E; inspect the module lock delta before delivery checks.
      Publication and mandatory exact-candidate checks remain tracked by delivery receipts.

## 7. Grouped upstream schemas

- [x] 7.1 Expose one protos target containing all upstream Anki schema files discovered by glob, while preserving the existing generated Go package boundaries.
- [x] 7.2 Build the complete schema target, verify its descriptor includes all upstream files, and check existing Go consumers and the native collection E2E before delivery.

## 8. Stdio gRPC worker

- [x] 8.1 Resolve the requested grpc-rust framework and obtain exact dependency approval; preserve existing locked dependencies.
- [x] 8.2 Extend the collection E2E before implementation with worker startup failure, crash reporting, process reaping, shutdown, and stderr separation cases.
- [x] 8.3 Replace the C ABI with a small Rust gRPC server and a Go process client; preserve upstream protobuf semantics.
- [x] 8.4 Add and render the considered-approaches Mermaid diagram; update API and build documentation.
- [x] 8.5 Verify the worker, Go E2E artifact, Clippy, schemas, and documentation; prepare the isolated PR #146 delivery scope. Exact-candidate checks and publication remain owned by delivery receipts.

## 9. gRPC project ownership

- [x] 9.1 Move the server and command generator into projects/anki_grpc with project documentation, executable, Clippy, and standard release metadata targets; keep third_party limited to external definitions.
- [x] 9.2 Point internal/ankiapi at the project server, update documentation and replace the detailed comparison with four independent general-approach graphs with separate card-shaped notes on the right, and preserve all dependency versions.
- [x] 9.3 Verify the independent upstream graph, Go E2E artifact, Clippy, docs/rendering, and delivery scope. Exact-candidate checks and publication remain owned by delivery receipts.
- [x] 9.4 Title the diagram containers Approach 1 through Approach 4 from top to bottom, simplify the translation-layer labels to Application and ABI/gRPC/Other, end every approach in Storage with SQLite (with Protobuf) on its second line, put Python and Rust on the second line of their Application nodes, and inspect the rendered dark WebP.

## 10. Server prerequisite extraction

- [x] 10.1 Extract projects/anki_grpc and required native build support into independently validated PR #147 against master; verify unchanged source, lock, server executable, and consumer E2E.
- [x] 10.2 Update the Go API's documentation and contracts to identify #147 as the server prerequisite while retaining the approach diagram in #146; delivery receipts own exact-candidate checks and publication.
- [x] 10.3 After #147 merges, rebase #146 onto master to remove identical shared server/native files; verify the remaining Go API, diagram, E2E, and delivery scope.

- [x] 10.4 Consume the prerequisite server CLI and wait for gRPC readiness instead of calling Initialize; verify startup deadlines, reaping, persistence, and default stdio with the existing E2E.

- [x] 10.5 Consume the directly imported Anki workspace from #147, remove the superseded helper snapshot, and verify the existing Go stdio E2E against the new build.

- [x] 10.6 Review the merged native build for simplification, record the supported generator limitations and remaining path coupling, and update consumer docs for the merged server. Delivery receipts own exact-candidate rebase validation.

## 11. Generated clients and shared project ownership

- [x] 11.1 Migrate the existing collection E2E to caller-constructed generated service clients before implementation; cover original service routing, overlapping collection methods, caller trailer options, native error status, persistence, cancellation, closed connections, and process failures.
- [x] 11.2 Move the reusable connection owner to projects/anki_grpc/pkg/backend, expose grpc.ClientConnInterface, and replace the generic Call API with one shared unary interceptor.
- [x] 11.3 Generate clients from the existing upstream service schemas and derive original service routes through Anki's descriptor helpers in Rust, retaining backend-prefixed paths and approved dependency versions.
- [x] 11.4 Update owner documentation and validate typed stdio E2E artifacts, existing TCP/Unix compatibility, Clippy, upstream-schema consumers, quality, semantic lint, and OpenSpec before updating PR #146.

## 12. Explicit backend lifecycle review

- [x] 12.1 Extend the E2E before implementation for constructor configuration, interface use, repeated startup, stop before start, invalid binary/arguments, cancellation after readiness, startup deadline cleanup, and OS-observed child reaping.
- [x] 12.2 Introduce Backend and BackendSTD, make New record the binary path and arguments, move execution into Start, and simplify startup into clear resource-owning steps; preserve typed clients and shared error handling.
- [x] 12.3 Update lifecycle documentation and validate native persistence artifacts, lifecycle failure reports, TCP/Unix compatibility, Clippy, quality, semantic lint, and OpenSpec.
      Delivery receipts own exact-candidate validation, publication, and the guarded reply to the review thread.

## 13. Lifecycle context and connection feedback

- [x] 13.1 Extend E2E cases before implementation for an error-returning Connection, unavailable handles during startup/shutdown, a worker surviving startup-context cancellation, and Stop deadline/cancellation cleanup with OS-observed reaping.
- [x] 13.2 Return connection availability errors, confine Start context to startup, add Stop context and forced cleanup on expiry, and update docs and examples while preserving typed service clients.
- [x] 13.3 Validate native persistence and lifecycle artifacts, TCP/Unix compatibility, Clippy, quality, semantic lint, and OpenSpec. Delivery receipts own publication and guarded replies to all three review threads.

## 14. Process command simplification

- [x] 14.1 Replace the cancellation-detached CommandContext with exec.Command, preserving startup context checks and explicit Stop ownership; validate existing lifecycle and persistence E2E artifacts before delivery.
- [x] 14.2 Fix the crash E2E's lifecycle race by constructing its generated client before killing the worker; verify retained-client errors and reaping repeatedly, then review the final diff.

## 15. Command-owned process cancellation

- [x] 15.1 Extend the existing native shutdown E2E before implementation for explicit cancellation and default grace expiry, retaining availability checks, caller context causes, repeated Stop results, and OS-observed reaping.
- [x] 15.2 Use CommandContext with a backend-owned lifetime context, cancel it on shutdown expiry, and remove manual process killing and ineffective WaitDelay configuration; preserve startup-only contexts and graceful EOF shutdown.
- [x] 15.3 Update lifecycle documentation, verify E2E artifacts, and review the complete final PR. Delivery receipts own required validation, publication, and final verification.

## 16. Constructor options and interceptor readability

- [x] 16.1 Migrate native E2E callers before implementation to BackendSTDOptions; verify copied arguments, configured shutdown timeout, invalid timeout rejection, and native failure details while retaining lifecycle coverage.
- [x] 16.2 Accept only BackendSTDOptions in New, expose documented binary/argument/timeout fields, preserve the five-second default, and flatten the shared interceptor without changing error causes or trailer handling.
- [x] 16.3 Update owner docs and contracts, inspect E2E artifacts, and review the full updated PR. Delivery receipts own required validation, publication, guarded replies, and final verification.
- [x] 16.4 Add a top-of-README architecture diagram with a standalone gRPC client, separate anki_grpc server (Rust binary) and anki_grpc library (Go) groups with single-line titles, gRPC server -> Anki library -> external Anki storage, and a Backend node in the Go library linked to the server group by a Manage edge; render and inspect the published assets.
- [x] 16.5 Replace the standalone client with Actor, add Actor -> Get gRPC connection -> Backend in the Go library while retaining the gRPC request edge, and inspect the rendered diagram and documentation assets.

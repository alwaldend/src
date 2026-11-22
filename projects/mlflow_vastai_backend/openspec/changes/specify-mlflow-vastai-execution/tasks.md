## 1. Job services

- [x] 1.1 Define the protobuf services and verify the Python client and Go clients can use the generated contracts.
- [ ] 1.2 Implement server and worker Stateless machines with SQLite persistence. Verify transition recovery and concurrent requests.
- [ ] 1.3 Add read-only rental discovery and the SSH worker client. Verify worker selection and confirm that ambiguous submission retries do not start a second workload.
- [x] 1.4 Implement Stateless CLI dispatch for `server serve` and `worker serve` with stub application entry points.
- [x] 1.5 Export the configured command machine as DOT through `state-machine graph` and render a themed SVG with Bazel-managed Graphviz.
- [x] 1.6 Embed the generated CLI graph in a dedicated state machines documentation page and verify the built page and asset.

## 2. MLflow integration

- [x] 2.1 Register the Python backend and implement `SubmittedRun`. Verify against a local gRPC fixture that a test project records metrics and an artifact under the returned run ID.
- [ ] 2.2 Verify completion, cancellation, SSH disconnect, server restart, and loss of a worker using a repeatable end-to-end fixture and preserved results.

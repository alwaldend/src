## 1. Job services

- [x] 1.1 Define the protobuf services and verify the Python client and Go clients can use the generated contracts.
- [x] 1.2 Implement the job server's Stateless machine with SQLite persistence. Verify transition recovery and concurrent requests.
- [x] 1.3 Add read-only rental discovery and invoke the worker client through SSH. Verify worker selection and confirm that ambiguous submission retries do not start a second workload.
- [x] 1.4 Compose CLI, application, and worker components with private machines, named methods, and stub service entry points.
- [x] 1.5 Export the composed machines as DOT through `state-machine graph` and render a themed SVG with legible machine titles.
- [x] 1.6 Embed the generated CLI graph in a dedicated state machines documentation page and verify the built page and asset.
- [x] 1.7 Implement the rental-local worker and command client. Verify artifact retrieval, SQLite persistence, replay, cancellation, and listener restart.
- [x] 1.8 Give the application a worker child that owns the listener and job machines. Verify concurrent job state, listener restart, and the combined graph.
- [x] 1.9 Run workloads under the listener and verify graceful shutdown, crash recovery, and the combined graph.
- [x] 1.10 Remove staged files after failed or cancelled workloads, retain logs and job records, and verify cleanup of background child processes.

## 2. MLflow integration

- [x] 2.1 Register the Python backend and implement `SubmittedRun`. Verify against a local gRPC fixture that a test project records metrics and an artifact under the returned run ID.
- [x] 2.2 Verify completion, cancellation, SSH disconnect, server restart, and loss of a worker using a repeatable end-to-end fixture and preserved results.

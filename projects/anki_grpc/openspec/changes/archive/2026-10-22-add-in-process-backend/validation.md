## Feature scope

The candidate adds optional direct C ABI execution on the Client API merged through PR #186. Core and application libraries remain pure Go. Executables own worker entrypoints and optional native imports; registered constructors create fresh clients and preserve exact argv.

## Current simplification

Client exposes Start, Connection and Stop. Backend ownership methods, registration flags and address classifiers are unnecessary: each client implementation owns its cleanup behavior. Path selects managed collection lifetime and Address selects its client, including caller-reserved TCP/Unix servers with shared paths. Borrowed Connection remains exclusive. AsyncClient.Done remains necessary to retain source reservations until native work has actually finished.

Network registration adapts NewGRPC(address) directly; the one-field options type and forwarding constructor are removed. Collection and application scope callbacks return only errors. WriteBack controls archive publication independently of callback results. Domain change results and the fixture's explicit export predicate remain meaningful and are preserved.

## Behavioral evidence

The existing TestAPINetworkWorkflow was extended before implementation. Its TCP and Unix cases failed at the ownership gate, then passed after removal. It writes a deck through managed API calls, verifies collection close and server reuse, reopens the database and compares the persisted result. Repeatable artifacts include api-tcp-managed.json, api-unix-managed.json and the corresponding databases.

The five affected E2E suites cover native/network calls, scoped factories and defaults, persistence, cancellation, source reservations, archive publication after errors and borrowed connections. Existing direct-native tests reopen persisted notes/cards through Anki and observe final cleanup. Immutable SQLite notes/cards observations supplement native readback; host SQLite cannot inspect Anki's unicase-backed tables.

## Delivery evidence

Both owner-local changes are validated strictly. Pure library and engine-free command builds, semantic lint, native lint, repository quality and the affected suites form the delivery gates. Exact candidates, remote leases and publication state belong to Git and repo-delivery receipts under ignored out/client-scope-simplification; raw logs and generated artifacts are kept there or in Bazel test outputs.

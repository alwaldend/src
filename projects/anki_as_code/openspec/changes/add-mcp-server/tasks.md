## 1. Acceptance and dependency setup

- [x] 1.1 Describe failure cases and write real stdio MCP acceptance tests before implementation; verify coverage includes discovery, invalid input, workflow/native mutations, cancellation, cleanup and repeatable output evidence.
- [x] 1.2 Pin the requested Go SDK and update owning Go/Bazel metadata; verify dependency resolution and a consumer build.

## 2. MCP implementation

- [x] 2.1 Implement descriptor-derived schemas and individual native RPC registration; verify full-catalog discovery and representative nested, enum, map, oneof, bytes and 64-bit payload behavior through MCP.
- [x] 2.2 Implement the ten workflow tools and stdio command using the existing App; verify config/state round trips, plan/apply/build, errors and process cleanup in E2E tests.
- [x] 2.3 Parallelize application and CLI suites in isolated Bazel shards and remove redundant large-fixture archive passes; verify all cases, original payload sizes, round-trip assertions and per-shard timings.
- [x] 2.4 Measure the large WASM cases independently and compare normal and optimized Rust compilation; document the optimized verification command and preserve measured timings.

## 3. Validation and delivery

- [x] 3.1 Document MCP setup, tool names and persistence limitations; verify generated command help and representative usage.
- [x] 3.2 Run affected E2E tests, semantic lint, repository quality and OpenSpec validation; preserve sanitized repeatable evidence in validation.md.
- [x] 3.3 Review the final implementation against the capability contract and verify no requested behavior remains incomplete. Git and delivery receipts own subsequent commit, synchronization, push and pull-request state.

## 4. Pull request feedback

- [x] 4.1 Use the SDK stdio transport directly and extract the workflow registration declaration; verify installed MCP discovery and workflow behavior.
- [x] 4.2 Document why ProtoJSON schemas require descriptor conversion instead of the SDK's Go-struct inference; retain complete catalog and payload verification.
- [x] 4.3 Make test sharding opt-in through Bazel's supported command-line setting and document the parallel E2E command.
- [x] 4.4 Write explicit-session failure scenarios before implementation, then expose col_open/col_close handles and pass their connections to existing methods. Verify reuse, source exclusion, deferred publication, native-error continuity, invalid handles, independent sessions and disconnect cleanup with repeatable evidence.
- [x] 4.5 Review the final changes against feedback and document acceptance evidence and validation commands. Exact-candidate check results, publication and review-thread resolution belong to delivery receipts and logs.

- [x] 4.6 Align the explicit session controls with the existing col\_\* namespace as col_open/col_close and name the application operation ColOpen; verify discovery and the session workflow across all CLI engines.

- [x] 4.7 Verify the initial MCP hosting route through the installed CLI E2E path. Task 4.12 replaces its dependency direction.
- [x] 4.8 Simplify collection lifetime results and tool dispatch without dropping path access, explicit handles or cleanup guarantees; verify native, WASM and combined MCP workflows.
- [x] 4.9 Replace the callback-to-session adapter with direct Store.Open ownership. Verify failed startup releases its source, cancellation and repeated close finish cleanup, archives publish on close, and source reservations outlive asynchronous native cleanup through lifecycle and MCP E2E coverage; inspect the close-wait cancellation path.
- [x] 4.10 Return the existing ankiclient.Client contract from ColOpen and retain one client per MCP handle. Delegate native ownership to the linked anki_grpc own-collection-client-lifetime change; verify source reservations, cancellation, publication and all CLI engines.
- [x] 4.11 Use one factory-created client with private collection state, remove the reservation-client wrapper, and place WriteBack in ColOpenOptions. Verify source exclusion through asynchronous completion, immediate reuse after Stop, publication and MCP behavior.

- [x] 4.12 Inject a separate MCP interface into CLI and construct its server with App at the executable entry point. Remove App-to-MCP dependencies and verify installed native, WASM and combined MCP workflows and generated CLI help.

- [x] 4.13 Omit implicit forwarded native aliases from MCP discovery while retaining all declared RPCs and canonical names. Verify 269 native tools and 281 total tools, complete schemas and native/workflow behavior across CLI engines.
- [x] 4.14 Keep each workflow tool's declaration and handler together and remove duplicate dispatch bookkeeping. Preserve envelope validation, source resolution and error handling; verify the existing MCP E2E workflows.

- [x] 4.15 Replace pairwise protobuf oneof exclusions with linear alternatives, preserving absent/null and single-member behavior. Exercise advertised schemas through an independent SDK consumer before implementation, then verify complete catalog discovery and native workflows.
- [x] 4.16 Remove MCP handle and registry shutdown flags already covered by Client.Stop/Connection and the SDK's drained-handler guarantee. Preserve serialized calls, cancellation, archive publication and disconnect cleanup through existing E2E scenarios.
- [x] 4.17 Add MCP to the simplified CLI architecture and document MCP discovery, explicit handles, calls and cleanup in a dedicated diagram. Render and inspect maintained dark WebP outputs and publish the diagram sources through the existing docs build.
- [x] 4.18 Collapse backend variants in the simplified CLI diagram into one Server reached through gRPC. Keep connectivity details in the existing backend-strategies diagram, align the overview prose, and inspect the updated dark WebP.

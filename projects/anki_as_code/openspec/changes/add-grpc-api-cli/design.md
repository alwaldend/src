## Context

The CLI uses Cobra and a reusable proto.Message codec. Master provides linked self-hosting and a shared application interface supporting managed paths and borrowed connections. Generated Go descriptors omit prose comments, but retain schemas. See proposal.md for motivation.

## Goals / Non-Goals

Use one descriptor-based executor and the existing protobuf runtime. Preserve upstream payloads and backend errors. Remote schema negotiation and new dependencies are outside scope.

## Decisions

- Generate one import-only Go file directly into internal/rpc. The project-local Go generator in tools/generate_anki_schemas reads the pinned native FileDescriptorSet in a declared Bazel run_binary action; the file is a build output, with no separate schema library or checked-in generated registration. Link all 25 existing Go schema packages through the shared third-party catalog. Iterate their compiled descriptors, sort services/methods, and construct Cobra commands at startup. Resolve message constructors through the Go registry and invoke through grpc.ClientConnInterface; runtime discovery needs no schema files or server.
- Preserve Backend service names. For missing backend methods, use descriptors of their matching collection service, following the existing upstream forwarding convention. Declared methods take precedence.
- Render schema-based help with request/response names and field types. Do not promise source comments stripped by protoc-gen-go.
- Keep --collection and --backend exclusive. Offline scope owns the worker and collection; connected scope owns only its gRPC connection. Explicit collection open/close calls require connected access; package export remains supported through the collection owner's consuming-call tracking.
- Use a proper boolean pair with default true and parse-order precedence. Pass the resulting value as WriteBack. Application collection scopes explicitly enable write-back by default, while the backend library remains opt-in. Validate and decode the request before backend startup and print successful responses after cleanup completes.
- Reuse master's linked native hosting and explicit current-executable worker arguments (server serve). No executable embedding, extraction, compression, cache, runfiles or PATH discovery is introduced. Keep server orchestration in App and native composition in main. The executable ignores SIGPIPE during native hosting, as the Rust CLI runtime does, so client disconnection cannot terminate the linked server. Extend the shared ankiascodeapi.App contract with one serialized APIRequest operation. App resolves the method, validates the request and manages native lifetime; CLI owns command discovery, input acquisition and output streams. Reuse CollectionSource for path and borrowed connections, with a separate explicit backend address for app-owned remote connections; retain caller-owned CollectionSource behavior for existing operations.
- Start the request timeout before reading input. Race blocking input acquisition against cancellation while preserving caller ownership of streams; a canceled arbitrary reader may finish when its owner closes it, and process exit releases CLI stdin. Decode and validate the complete request before opening native storage.
- Expose one WithCollection(ctx, CollectionOptions, callback) entry point. CollectionOptions contains Path and WriteBack; omitted WriteBack disables archive publication. Application workflows explicitly enable publication and API calls pass the selected CLI setting. Collection.Store exposes only offline collection scopes; App owns remote client startup and cleanup directly. BackendCollection owns publication independently of callback change reports. Adjust existing callers/tests and current documentation to match.

## Risks / Trade-offs

- Missing linked schemas: audit the pinned proto tree and exercise tags/i18n in native CLI E2E.
- Explicit lifecycle calls can consume offline storage: reject direct open/close/import in managed mode and retain backend tracking for exports.
- Source archives can be repackaged by read commands: preservation means unchanged collection content, as agreed.
- Dynamic registry failures: validate descriptors and constructors during command assembly and test complete help without a server.

## Delivery

Consume backend, complete third-party bindings, linked hosting and shared application contracts from master. The CLI owns descriptor registration. Resolve mechanical rebase conflicts using current package names, worker arguments and borrowed-source lifetime; retain explicit write-back options and semantic native-content assertions. No backend packaging changes remain in this PR.

## Boundary integration decision

Verdict: revise the inherited callback API. It kept command-side managed-lifetime policy and native preparation outside App, conflicting with merged package boundaries. One APIRequest containing route, payload, encodings and source selection moves orchestration to App without per-method wrappers. Compiled descriptors remain the authority in internal/rpc; CollectionSource remains the authority for path/borrowed selection. The extra request metadata type replaces a callback contract and avoids growing positional parameters. Existing installed/native E2E and a caller-owned native API workflow verify both access and failure paths.

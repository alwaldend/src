## 1. Existing preservation corrections

- [x] 1.1 Resolve public paths in filesystem order and refuse required default-deck
      removal during planning; retain the existing behavioral coverage.

## 2. Caller-owned boundaries

- [x] 2.1 Adapt behavioral coverage before implementation for export alongside
      a source, ordinary overwrite behavior, raw clients and prepared handles.
- [x] 2.2 Remove protected-source plumbing and export source/output guards.
- [x] 2.3 Own single-use collection storage and cleanup in the application,
      remove generic client collection phases and migrate API/fixture consumers.

## 3. Validation

- [x] 3.1 Validate affected native, WASM, MCP and collection-scope workflows,
      inspect artifacts, and complete strict OpenSpec and delivery gates.

The user superseded the earlier source-alias refusal requirement with caller-owned
export destinations. Prior candidate evidence remains in Git and delivery receipts.

Validation completed across 45 native, WASM, CLI/MCP and repository targets,
with strict owner-local OpenSpec checks and semantic lint. The WASM-only CLI
suite uses two independent shards after its complete unsharded run exceeded
five minutes; all cases and each shard's original timeout remain unchanged.

## 4. Shared client capabilities

- [x] 4.1 Return ankiclient.Client from ColOpen and Store.Open, remove the
      separate Collection interface, and use StartableClient where startup is owned.
- [x] 4.2 Verify the shared interface migration and retained collection cleanup
      behavior through the affected existing workflows and delivery checks.

## 5. Narrow lifecycle integration

- [x] 5.1 Inject constructor functions and observe optional transport Done
      through a private collection capability without changing cleanup policy.
- [x] 5.2 Remove the ServerServe application forwarding method; keep the serve
      argument prefix in CLI and verify existing hosting and collection workflows.

## 6. Handwritten public symbols

- [x] 6.1 Audit all handwritten exports against callers and required framework
      or serialization contracts, including declarations outside Go interfaces.
- [x] 6.2 Replace the redundant collection.CollectionOptions alias with
      ankiascodeapi.ColOpenOptions, remove the test-only filesystem.FS.ExpandPath
      helper, hide resource.Store.Search and Presets fields used only in resource,
      and replace codec.Format with string parameters and untyped JSON/TOML constants.
- [x] 6.3 Verify affected existing workflows and complete documentation,
      review and delivery checks for the visibility and alias reductions.

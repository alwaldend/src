# Unified client validation

The reusable client owns its transport and optional private storage state.
Existing native E2E cases are migrated to OpenCollection and CloseCollection on
that same client. They cover unprepared access, successful operations after
startup-context cancellation, invalid sources, stop before start, repeated stop,
archive write-back, retained connections and uncertain native work.

The preimplementation TestUnifiedClientCollections case opens and closes two
databases on one factory-created client, checks that old collection connections
fail, uses raw native progress between collections and observes final Done after
Stop. It retains unified-client.json as repeatable evidence. A concurrent-close
case checks that terminal Stop honors its collection deadline while another close
is awaiting active native work. Existing database,
archive and per-scenario JSON outputs remain available for inspection.

Anki as Code Store and native/WASM/combined MCP E2E exercise ColOpen source
reservations, publication, disabled write-back, cancellation and disconnect.
The pending native-call fixture checks that reservations survive cleanup errors
until native completion. Publication, failure and panic cases remain in scope.

The unified-client baseline passed all 17 affected test targets with optimized
builds, four shards and eight concurrent test jobs. The concurrent-close fixture
failed before the deadline-aware join fix and passed afterward. Native and WASM
application suites, all three CLI engines, collection workflows and transport
fixtures passed together. Generated CLI reference, WASI binary and both updated
diagrams also build. Exact candidate identities, mandatory quality/lint checks
and publication state belong to delivery receipts under ignored
out/anki-mcp-support.

The private-state simplification passed the existing real-worker lifecycle,
consuming-call, publication and concurrent-close cases. The focused invocation
passed six targets covering ankiclient, the collection Store, CLI commands and
native/WASM/combined CLI workflows, and built the WASI executable. It used optimized
builds, four shards and eight concurrent test jobs and finished in 78.326 seconds;
five tests executed and one passing result was reused. The MCP catalog contains
269 declared native RPCs, ten workflows and two session controls: 281 tools.

Evidence is retained in out/anki-mcp-support/simplified-behavior.log and
simplified-behavior.bep.json. The lifecycle artifacts confirm transport reuse,
bounded concurrent close, archive publication and retained restores on failure.
Strict validation of both OpenSpec changes, all 26 repository quality tests and
semantic lint of the 22 affected packages pass. Delivery receipts own exact
candidate validation and publication state.

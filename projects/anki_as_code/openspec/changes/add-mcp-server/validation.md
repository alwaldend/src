# MCP validation

The discovery contract is 281 tools at the pinned Anki revision: 269 declared
native RPCs, ten workflows and two collection-session controls. It omits 226
implicit forwarded aliases while retaining declared Backend methods and existing
canonical names. Focused native, WASM and combined-engine MCP checks verify the
reduced catalog and workflows. Earlier 507-tool observations below are baseline
evidence. Publication identity and exact-candidate validation belong to the
delivery receipt rather than this coverage report.

The complete WASM application and WASM-only CLI workflows use two isolated
Bazel shards. Other targets retain default sharding. Earlier timing observations
below used explicit shard counts and do not describe current target defaults.
The SDK's `StdioTransport`
now owns protocol streams directly, replacing the custom output wrapper.

Explicit session E2E covers the user-selected open/close handle model. Several
native and workflow calls reuse one worker and prepared connection. Another
process observes the original archive until close, then observes the published
change. Duplicate opens and mixed handle/path controls fail. Known native
errors leave the handle usable; closed handles are rejected. Separate handles
own independent collections. Physical MCP disconnect closes remaining handles,
publishes enabled archives and stops their workers. Disabling write_back on
open leaves the archive unchanged. `mcp-worker-lifetimes.log` records process
starts and graceful exits as repeatable public fixture evidence.

App.ColOpen returns the shared ankiclient.Client capability, already prepared.
MCP borrows Connection for existing App methods and calls Stop to publish and
release it. The application's private scope owns native storage and cleanup;
its factory-created transport implements StartableClient. Store retains source
exclusion through private completion tracking. Client exposes only Connection
and Stop, while StartableClient embeds Client and adds Start. Explicit handles
publish on close or disconnect without an application-specific public interface.

CLI groups protocol commands behind its private mcpServer interface, receiving
the concrete mcp.New(app) value. Its only method is Serve(context.Context) error;
no MCP interface is exported or added to ankiascodeapi. Task 4.22 and exact-candidate
delivery receipts own validation of this consumer-owned contract. The existing
CLI dispatch and native, WASM and combined-engine CLI/MCP workflows passed.
Retained tool, state and plan artifacts confirm discovery and convergence.

## Historical function-hosting evidence

The superseded hosting-function follow-up passed existing native, WASM and combined-engine
CLI/MCP workflows and focused command dispatch checks. It injected the concrete
MCP server's Serve method into CLI beside App. Task 4.22 replaces that function
dependency with a private CLI interface; these earlier results do not validate
the private-interface candidate.

## Behavioral evidence

The installed executable is driven over real stdio by the Go SDK client, with
an empty PATH. Discovery also runs without a compiled Anki engine and with an
unavailable temporary directory. An independent SDK consumer resolves every
advertised input and output schema.

The E2E workflow exercises all ten workflow tools, native deck creation and
lookup, recursive search nodes, media bytes, string-encoded native identities,
state maps, archive building, and disabled archive write-back. Invalid enums,
oneofs, integer overflow, unknown payload fields and unknown envelope keys fail
before opening an empty collection. A native execution failure leaves the MCP
session usable. Applying a desired change produces a subsequent empty plan;
the original deck identity and the build's base archive remain unchanged.

Cancellation reaches an active native RPC, and a physical stdio disconnect
allows the server to cancel the request and exit successfully. Normal shutdown
leaves no restored collection scratch. The client distinguishes graceful
session close from physical transport loss in the disconnect fixture.

Repeatable outputs are `mcp-tools.json`, `mcp-state.json`, and `mcp-plan.json`
under each CLI test's `test.outputs` directory. They contain public fixture
data. Reproduce the dedicated acceptance cases with:

```sh
bazel_agent bazel test --repo_contents_cache= \
  //projects/anki_as_code/cmd/anki_as_code:cli_workflow_test \
  --test_filter=TestMCP --test_output=errors
```

## Repository checks

- The native, WASM, and combined-engine CLI suites pass, including MCP cases.
- The existing native application workflows and the other project tests pass.
- Native and WASM application suites pass with eight isolated shards each,
  including the original 17 MiB field and 513 MiB media fixtures.
- Semantic lint passes for the affected CLI, MCP, RPC, documentation, OpenSpec,
  and Go dependency packages.
- All 26 root repository quality checks pass.
- The standalone WASI CLI builds successfully.
- Strict OpenSpec validation and Bzlmod dependency resolution pass.
- Generated CLI help contains `mcp serve`. The updated dependency diagram is
  rendered through its existing Bazel rule and visually inspected.

The fastbuild project and root-quality invocation passed all 42 test targets in
124.856 seconds. It executed 27 targets and reused 15 passing results. The
application suites were cached from the preceding parallel run; this duration
is not a cold-suite benchmark. Observed fastbuild shard durations are:

| Suite               | Shards | Longest shard |
| ------------------- | -----: | ------------: |
| Native application  |      8 |      25.986 s |
| WASM application    |      8 |     257.957 s |
| Native CLI          |      4 |      11.956 s |
| Combined-engine CLI |      4 |      43.824 s |
| WASM CLI            |      4 |     118.058 s |

The application timings were measured with the host's four-worker limit; CLI
timings were measured with eight workers. Timing evidence is retained in the
task-local Build Event Protocol log. The full-size application cases passed
under their original timeout categories without extending deadlines.

An isolated run measured the individual large WASM cases, followed by the same
cases with `-c opt`. Both runs used eight workers and fresh shard-local
compilation caches. The Rust action parameters confirm `opt-level=0` for the
default fastbuild guest. Optimization reduced the media case by 6.4 times:

| WASM case     | Fastbuild | Optimized |
| ------------- | --------: | --------: |
| 17 MiB field  |   25.30 s |    9.08 s |
| 513 MiB media |  111.01 s |   17.37 s |

These are complete fixture workflows, including native archive creation,
rebuild, restoration and assertions, rather than single archive reads. The
optimized invocation took 32.678 seconds including build/cache work; its build
used existing action-cache results and is not a clean-build measurement.
The owning README now selects `-c opt` for full-project verification.

## Validation environment

The initial native build encountered `anki-native-repository-cache-symlinks`:
a cached Git-based Rust dependency contained absolute symlinks into another
worktree's output base. Disabling the repository contents cache for validation
materialized correct dependency sources without changing pins or shared host
configuration. Logs remain under ignored `out/anki-mcp-support/`.

The initial full-project command's 600-second outer deadline and a subsequent
serial run's 900-second test deadline were insufficient. The latter interrupted
archive publication in the 513 MiB fixture, causing its worker termination and
cleanup error. The initial parallel validation used eight App shards and four CLI
shards with the original timeout categories. Each process has its own test
environment and compilation cache.

Large-media fixture setup now exports configs and installs bindings on the
database before creating the archive. Its full-size build/import remains, and
convergence inspection borrows the already-restored connection. The large-field
case asserts convergence from its existing CLI plan output instead of issuing
a duplicate plan. Sharding also exposed an implicit CLI fixture dependency on
another test registering a subprocess factory; the fixture now constructs that
client explicitly.

This host defaults to four Bazel jobs and quarter-CPU allocation. Task-local
validation overrides these to eight jobs and half the host's 16 CPUs; shared
configuration is unchanged.

## Session improvements

`anki-e2e-serial-archive-work` is fixture-tested: serial suites and repeated
large-archive publication caused the observed deadlines. Isolated shards and
database-first fixture setup resolve that task-local bottleneck while retaining
the payloads and round-trip assertions. `anki-cli-fixture-registration-order`
is fixture-tested: sharding exposed implicit process-global registration, now
replaced by explicit fixture construction. `anki-native-repository-cache-symlinks`
remains an environment workaround, with the owning dependency pins unchanged.
Future runs can resume from the OpenSpec tasks, this bounded report and the
delivery receipt; no shared skill or host-policy changes are needed.

## Historical unified-client and hosting follow-up

This evidence describes the earlier unified-storage implementation, superseded
by the current application scope and shared client capabilities.

CLI receives separate App and MCP dependencies. The concrete MCP Server receives
App at construction and tools call its application methods. App no longer imports
or delegates to MCP. ColOpen accepts Path, Address and WriteBack in ColOpenOptions
and returns the same factory-created client after starting and opening it.

The linked anki_grpc change removes separate collection and reservation clients.
Store retains source/session policy through a plain record observing Client.Done;
WithCollection uses the same Open/Connection/Stop path. Existing lifecycle and
real stdio MCP E2E remain the acceptance cases for publication, source exclusion,
worker reuse, cancellation and disconnect. Library E2E covers collection reuse
on one client and raw access between collection phases.

Before alias reduction, the unified-client follow-up passed 17 affected test
targets with optimized builds, four shards and eight concurrent test jobs,
including native/WASM application suites and all three MCP CLI engines.
Those baseline catalog artifacts contain 507 tools; each engine records
four worker starts and four stops across explicit-session scenarios. Unified-client
artifacts verify transport reuse, publication-failure retention and final native
completion. The concurrent-close fixture failed before its deadline fix and passes
afterward. Generated command reference, WASI binary and updated diagrams build.
Mandatory quality/lint checks and final candidate publication remain recorded in
delivery receipts under ignored out/anki-mcp-support.

## Catalog and declaration simplification

The user approved omitting forwarded RPC aliases and colocating workflow tool
declarations with their handlers. The acceptance catalog has 269 declared native
RPCs and 281 total tools. Declared Backend methods, canonical names, descriptor
schemas, envelope checks, source resolution and workflow behavior remain covered.
The focused simplification run passed all six test targets: native, WASM and
combined-engine CLI workflows, CLI commands, collection Store and ankiclient
lifecycle cases. It also built the WASI executable. The optimized invocation used
four shards and eight concurrent test jobs and finished in 78.326 seconds,
executing five tests and reusing one passing result. This is a combined validation
duration, not a cold-build benchmark.

The current catalog is 269 declared native RPCs plus ten workflows and two session
controls, totaling 281 tools. Tool declarations now contain their handlers;
existing E2E workflows verify schema discovery, source selection and execution.
The log and Build Event Protocol evidence are retained in
out/anki-mcp-support/simplified-behavior.log and simplified-behavior.bep.json.
Before linear oneof constraints, all three catalog artifacts contain the same 281 tools and are 2,381,046 bytes,
down from 4,406,769 bytes with forwarded aliases. Each engine records four worker
starts and four stops. Strict validation of both OpenSpec changes, all 26 repository
quality tests and semantic lint of the 22 affected packages pass. Exact candidate
validation and publication state belong to delivery receipts.

## Linear schemas and MCP diagrams

The independent SDK consumer now exercises nine advertised-schema cases for
absent/null oneof members, one selected member and conflicting selections. This
extension passed before and after replacing pairwise exclusions with linear
alternatives. Native ProtoJSON validation remains unchanged.

All MCP cases passed across native, WASM and combined CLI targets using four
shards and eight concurrent test jobs. The invocation also built WASI and rendered
both diagram sources, finishing in 14.548 seconds. This focused run selects
TestMCP cases; it does not rerun unrelated CLI workflows. Every engine emitted the
same 281-tool catalog and recorded four worker starts and four stops. Removing
MCP's redundant lifecycle flags preserves explicit close and disconnect cleanup.

The catalog is now 1,927,707 bytes with indentation (previously 2,381,046), or
791,379 compact bytes (previously 910,475). Linear oneof constraints reduce the
compact catalog by 13.1 percent without removing tools or changing payloads.
Evidence is retained in out/anki-mcp-support/mcp-diagram-behavior.log,
mcp-diagram-behavior.bep.json and mcp-diagram-artifacts.json.

The simplified CLI diagram includes CLI-to-MCP hosting, stdio agent access and
MCP-to-App tool dispatch. The MCP sequence shows discovery without a worker,
explicit handle reuse and publication on close/disconnect. Both use maintained
site SVG targets and dark WebP previews. Exact candidate checks, rendered-output
inspection and publication state belong to delivery receipts and task evidence.

The CLI overview now collapses all Anki backend variants into one Server reached
through gRPC. The existing backend-selection diagram owns deployment and transport
details. The overview SVG and dark WebP were rebuilt and visually inspected; the
README describes the same abstraction. This follow-up changes documentation only,
so the preceding runtime behavior evidence remains applicable. Final-candidate
diagram, quality, lint and specification checks belong to delivery receipts.

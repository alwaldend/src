## Context

Extract the browser-only implementation from immutable source commit
`b3c926655884ab070d628314b907950707306cc9` onto trunk
`a2c655b20432b78e4feb75e16da1a95b73b3d30c`. PR1 merged as PR195 at
`98559dd11cbf65dc0658feb3bd80be0b468d6862`; PR2 is rebased onto that master
commit without duplicating prerequisite source.

## Goals / Non-Goals

Preserve generated protobuf clients, owned buffers, explicit startup/shutdown,
and caller-owned JavaScript storage. The browser host requires neither HTTP/2 nor JSPI.
Application workflows, durable persistence and native adaptation remain outside
this change. The approved Emscripten provider and build support are composed
from their scoped owners; the Rust target and shared host contract are unchanged.

## Decisions

- Consume PR1's shared direct-client and Rust ABI. Build a WASI reactor with
  one-time initialization using the existing pinned Rust reactor CRT.
- Start joins factory completion, checking its context before and after
  initialization. Active startup registration lets concurrent Stop retain its
  waiting deadline while Done observes eventual cleanup. Canceled startup never
  publishes the backend. Command cancellation keeps its existing behavior.
- Keep production `server.js` and `wasi.js` as the browser runtime boundary.
  JavaScript adapts the caller's Node-compatible Sync API directly, owns only
  guest descriptors/positions/cookies/rights and provider descriptor mappings,
  and attempts descriptor cleanup on initialization failure, close and trap.
  Go owns only timezone callbacks until backend destruction.
- NewClient accepts only a factory. The caller captures fs/root; the factory
  receives the timezone host. createBrowserServer(module, fs, {root, timezone})
  mounts the absolute provider root at guest /. Provider and root identity are
  caller responsibilities; there is no backing store or identity enforcement.
- installGoFilesystem(fs), exported through :go, runs after wasm_exec.js import
  and before Go.run. It preserves runtime output, provider method receivers,
  fixed cwd / and POSIX path resolution. Go uses NewOsFs/BasePathFs; root
  selection does not alter globals. Keep the provider and globals stable.
- Serialize commands per browser backend across the complete JavaScript call.
  Timezone callbacks can still schedule Go callers; direct JavaScript filesystem
  calls do not justify weakening the existing reactor serialization contract. Queued calls remain owned by shared active-command cleanup.
- Extend the existing generated-client Go harness with a storage mode, launched
  in a Chromium worker by a focused adjacent runner. Use the repository's pinned
  Puppeteer/Chromium and Go SDK, with no imports from `anki_as_code`.
- Keep Node command/lifecycle coverage separate from actual Chromium host
  acceptance. Capture hashes, runtime identity, completed cases, host calls and
  failures in `browser-storage.json` even when a browser assertion fails.
- PR1 owns both `internal/backend/BUILD.bazel` and
  `internal/directclient/BUILD.bazel`, including browser visibility grants.
  PR2 excludes those files and has no supplemental prerequisite patch.

## Failure Cases Before Implementation

1. A collection cannot be created when its shared parent path is a file;
   the native error must remain observable and a later valid open must work.
2. Rust writes must appear in the caller-owned JavaScript filesystem. After
   collection close and backend shutdown, a fresh backend using that provider and root
   must reopen the same SQLite file and read the note's ID and Unicode fields.
3. Binary media created by Rust must be byte-identical when read by Go and
   survive backend replacement. A Go-created media file must be visible to Rust
   media inspection after reopen.
4. Invalid note lookup must fail without poisoning the backend. Calls retained
   after shutdown must fail, and Done must observe completed cleanup.
5. A non-UTC named browser timezone must load in Go. Actual Rust operations must
   traverse timezone imports, with known winter/summer offsets checked on the
   supplied host endpoint.
6. Startup, worker error, Go exit or timeout must produce a failing retained
   report; collection failure must not be misreported as a passing Node check.
7. Cancel startup while the JavaScript factory Promise is unresolved. Start
   must remain pending until that Promise settles. Concurrent Stop may stop
   waiting on its own context, but Done must stay open until the factory and
   cleanup finish. After release, Start must report cancellation, expose no
   connection and close the backend exactly once.
8. If a startup assertion fails before factory release, the Go harness must
   release the Promise and join Start, Stop and Done before reporting that
   failure. Pending cleanup callbacks must not outlive the Go runtime or hide
   the original contract error.
9. Concurrent generated media calls must not reenter the same Rust reactor
   while synchronous host callbacks run. Both calls
   must complete, and their distinct bytes must survive backend replacement.

## Risks / Trade-offs

Cancellation during initialization waits for the factory to finish; a stuck
factory delays startup failure and cleanup. Synchronous Rust commands cannot
be interrupted until they return. The caller owns filesystem durability and provider selection. Chromium is
the acceptance browser, not a cross-browser matrix. Independent master-based
validation must pass after the prerequisite merge before publication.

## JavaScript filesystem replacement acceptance

Before replacing the host, write production-host Node failure cases covering:
provider receiver preservation, mounted root, missing parents, exclusive/create
and directory flags, read/write rights, positional and append I/O, unsafe offset
and memory rejection, provider inode identity, and descriptor cleanup on init
failure, normal close and actual guest trap. Reuse the existing generated-client
collection/media workflow over NewOsFs/BasePathFs, including concurrent media,
backend replacement and timezone callbacks.

The Chromium runner accepts an explicit caller-supplied module exporting fs/root.
Its Bazel target supplies an application-owned Emscripten fixture and declares
the adapter and generated module as runtime inputs. Node success does not satisfy
Chromium acceptance. Historical Afero-host evidence remains identified by its
old candidate and does not validate the new filesystem boundary.

## Production-host I/O correction

Preserve completed read/write counts across later recognized provider errors,
including chunk and iovec boundaries. Before any progress, return the errno.
Unknown exceptions remain terminal. For sync-after-write failure, immediately
returning errno would discard progress and invite replay; ignoring it would
hide failed durability. Proceed with a descriptor-local deferred sync error:
stop after the completed bytes and report that error once on the next write,
sync, datasync or close. Close invalidates the guest mapping before one provider
close attempt. If close also fails, its error takes precedence over the retained
sync error because WASI returns one errno. Recognized close errors remain ordinary
filesystem failures; unknown close exceptions remain terminal. Do not aggregate
the two errors into an unmapped exception. A failed provider close leaves release
uncertain, so neither later guest close nor backend shutdown retries it. Only the
fault fixture can reclaim a descriptor it deliberately left open. A later
successful explicit sync establishes the provider's own flush
guarantee; a transferred-byte count alone does not establish durability.

Forward lookup flags into open. Inspect the final component with lstat when
no-follow is requested and add native O_NOFOLLOW when the provider exposes it.
Preserve explicit follow and exclusive creation; do not expand path sandboxing
or introduce link creation. Map concrete symbolic filesystem error codes,
including quota, stale handles and remote I/O, while leaving unknown errors
terminal. Preview 1 has no EREMOTEIO value, so map it to IO.

## Standalone Emscripten MEMFS composition

Use integrity-checked official SDK inputs under
`third_party/com_github_emscripten_core_emsdk`. Public build support under
`tools/emscripten` uses declared repository Python and Node toolchains in
sandboxed network-free actions. The `emscripten_memfs_files` target exports an
ordinary modularized ES module whose async factory creates independent classic
MEMFS with `.FS` and `.ERRNO_CODES`. No upstream source fragments or handwritten
backing filesystem are used. The selected pinned output is JavaScript-only;
its artifact E2E establishes loading without a WebAssembly global. Upstream
deprecates this output mode, so upgrades must reassess it and declared artifacts.

Before artifact implementation, its owner specified failures for missing
exports, shared factory state, missing parents/files, invalid descriptors,
stream positions, rename/truncate, Unicode/binary corruption and undeclared
runtime dependencies. Provider translation and its standalone evidence belong
to `add-emscripten-filesystem-adapter`; the [browser README](../../../pkg/ankigrpc/browser/README.md)
owns consumer usage and limitations. This change owns composition with the
generated Go clients and real Rust reactor in Chromium, without SDK imports.

Rust retains its existing WASI target and buffer ABI. Emscripten supplies shared
storage only; a Rust Emscripten target would require a separate linking/runtime
and filesystem port. The parent resolved that alternative in favor of this
composition. Publication still requires matching Chromium and SDK evidence,
fresh review and exact-candidate delivery gates.

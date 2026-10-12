## Context

PR #146 owns the Go Anki API and considered-approaches diagram. PR #147 owns
projects/anki_grpc and the native third-party build support extracted from
f3869e7700f33246f55d77c5927a8341db553387. PR #128 owns collection reconciliation.

## Goals / Non-Goals

Provide typed upstream service clients and Go ownership of process lifecycle.
Collection business logic belongs to Anki and its CLI consumer. Native build
definitions and the dependency lock remain unchanged.

## Decisions

Backend exposes Start(context.Context) error, Connection() (grpc.ClientConnInterface,
error), and Stop(context.Context) error.
BackendSTD implements it over stdin/stdout. New accepts one BackendSTDOptions
value with Binary, Arguments, and ShutdownTimeout, copies its arguments, and
defaults a zero timeout to five seconds. Start rejects negative timeouts before
creating resources. Construction only records configuration; Start creates the
process and waits for readiness.
The Start context controls startup only; cancellation after success leaves the
server running. The child uses exec.CommandContext with a backend-owned lifetime
context, separate from Start's context. Stop closes stdin, bounds graceful shutdown
by its context and the configured timeout, and uses context.AfterFunc to cancel the
command's context on expiry. os/exec owns forced termination; Go still owns Wait
and pipe cleanup. No manual Process.Kill or WaitDelay configuration is needed.
Environment and diagnostics
inherit the parent process, with stderr reserved for diagnostics. Separate helpers
own pipe creation, child startup, connection creation, and readiness waiting.
The backend owns pipes, one gRPC connection, process wait, and shutdown.
HTTP/2 readiness obeys its deadline and reaps a failed child. The reusable package
lives in projects/anki_grpc/pkg/backend and exposes its grpc.ClientConnInterface.
The name reflects ownership of the server process as well as its connection;
generated service clients belong to the caller.
Callers construct generated upstream clients; Go adds no per-operation wrappers.
One unary interceptor preserves context and gRPC causes, rejects calls after
stop, and decodes native errors from binary trailers. Caller call options,
including trailer capture, remain effective. Cancellation stops waiting without
promising rollback.

Use the existing rules_go go_grpc_v2 compiler alongside go_proto for the exposed
upstream schemas. Their source remains unmodified. Anki implicitly forwards
normal service methods through its backend services, so the Rust descriptor
generator maps original service paths to the official backend indices as well
as retaining backend-prefixed paths. For methods declared in both services,
verify their request/response descriptors agree before adding the original path.
Do not duplicate Anki's dispatch or operation implementations.

Stop signals EOF on stdin while draining stdout, waits for exit, and kills an
unresponsive worker when the Stop context expires or after the configured grace
period, then reaps it. A cancelled shutdown preserves its context cause alongside
process exit errors. Repeated Stop preserves its
result. The default grace is five seconds. Failed operations are never
automatically replayed or sent to a restarted process. The existing synthetic
E2E owns behavioral verification.

Each instance accepts one Start attempt; a failed, running, or stopped instance
cannot start another child. Stop is safe before Start and after failed startup.
Lifecycle configuration and cleanup are serialized. Connection returns an error
before readiness, during startup/shutdown, and after Stop, rather than a nil usable
handle. A readiness flag is set only after Start's handshake succeeds. Stop releases
the lifecycle mutex before waiting so Connection can report shutdown immediately.
Callers construct generated clients after successful Start.
Retained clients reject calls once shutdown begins. A constructor copies its
arguments so caller mutation cannot change the eventual process command.

Consume the projects/anki_grpc Bazel runtime target from master. PR #147
merged at 51a66d3f0fd2b0744bf18ab93fcef28e07e977e9; the rebase removes shared
server/native files and generated metadata from #146's diff. Build ownership
and dependency pins remain with the merged server and third-party package.
The server initializes Anki at startup; the client waits for gRPC readiness.
Server build details belong to its README and third-party package.

Maintain four independent approach graphs, numbered from top to bottom. Each
ends in Storage with SQLite (with Protobuf) on its second line. Python and Rust
are the second lines of their Application labels; the translation layer lists
ABI/gRPC/Other. Standalone card-shaped comments remain to the right.

## Risks / Trade-offs

- Cancellation can leave an Anki mutation running; callers inspect state before
  deciding what to do next.
- Worker failure terminates pending calls; Go reports the error and reaps it.
- Server/native support currently targets Linux x86_64, as documented by #147.

## Current state and next action

PR #147 is merged, and #146 is rebased onto the exact merge commit. The requested
follow-up moves the Go API into the server project and replaces generic Call
with caller-constructed generated clients, requiring small service routing and
schema compiler changes in this PR.
The merged Anki build review found no demonstrated correctness defects. Five
workspace crates and three upstream build scripts are required by the pinned
source. rules_rust deliberately leaves workspace-member targets to callers;
replacing them with a generator or opaque Cargo action would add a new layer.
The bounded review is in out/anki-grpc-project/extraction/build-review-after-merge.md.
The synthetic E2E was migrated before implementation and failed on the missing
connection accessor and generated clients. Typed stdio calls and existing TCP/Unix
paths now pass. The emitted SQLite database independently contains the expected
note fields and generated card. The moved backend, TCP/Unix compatibility, Clippy,
schema consumers, rendered documentation, quality, semantic lint, and strict
OpenSpec checks pass. Delivery receipts own final exact-candidate validation and
publication. The lifecycle review is implemented: Backend exposes Start,
Connection, and Stop; BackendSTD's constructor is side-effect free; startup uses
separate resource-owning steps. E2E cases were added before implementation and
now verify interface use, argument copying, repeated startup, invalid commands,
context cancellation, and OS-observed child reaping. Quality, lint, Clippy,
TCP/Unix compatibility, rendered documentation, and strict OpenSpec checks pass.
Delivery receipts under out/anki-cargo-api/lifecycle own exact-candidate checks,
publication, and the guarded reply to thread PRRT_kwDOOgezIM6rzmO8.

The follow-up context review is implemented. Connection reports unavailable
states; Start's context is detached from the successful process lifetime;
Stop's context bounds graceful shutdown before forced termination and reaping.
The collection E2E, startup/stop transition cases, startup-context isolation,
Stop deadline artifact, OS reaping checks, TCP/Unix compatibility, quality,
lint, Clippy, and strict OpenSpec checks pass. Delivery receipts under
out/anki-cargo-api/review-followup own exact-candidate publication and the three
guarded review replies.

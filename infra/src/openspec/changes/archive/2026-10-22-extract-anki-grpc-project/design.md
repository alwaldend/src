## Context

The extraction starts from PR #146 head f3869e7700f33246f55d77c5927a8341db553387.
Master is 6f02d8079faa80810f0f51cff2d3b56090777a1a. The user additionally requested
a transport-selecting CLI with stdio as default, TCP and Unix sockets, and
backend initialization at server startup rather than through a custom RPC.

## Goals / Non-Goals

The server PR builds independently from master. Go transport, reconciliation,
and the approach diagram retain their existing owners. Anki's library build
remains under third_party; projects/anki_grpc depends on that native target and
contains only our server and its command lookup generator.

## Decisions

Keep approved Cargo pins and import the upstream source workspace directly.
Use standard crate_universe dependency helpers and rules_rust targets. Use the standard library for CLI parsing
and the already locked Tokio network APIs for listeners. Transport support
requires no additional packages beyond its approved set. Backend startup settings, when needed, are CLI parameters rather than
a separate initialization protocol. The CLI exposes --transport stdio|tcp|unix and --listen, with explicit
endpoints required for network transports. Help and invalid arguments finish
before starting the backend or emitting protocol traffic.

Initialize one backend with Anki's default BackendInit options before serving.
All client connections share it; Anki owns collection and command synchronization.
Remove the custom Initialize RPC. The Go client waits for HTTP/2 readiness and
retains its startup deadline and child cleanup guarantees. Stdio serves one
connection; TCP and Unix listeners accept later clients after disconnects.
They use the same opaque codec and generated command mapping. Existing Unix
socket paths are never removed automatically; the caller owns socket cleanup.

Both PRs target master. The server PR is the merge prerequisite. PR #146 keeps
the required shared snapshot until this prerequisite merges, then rebases to
retain its Go API and diagram. Removing the server immediately would make its
Go test runtime inputs unresolved and prevent the required delivery checks.
The prerequisite source remains authoritative throughout the transition.

Write native CLI/network E2E cases before implementing the transports and
startup initialization. Emit persisted synthetic collections and verification
reports, inspect them independently, and run the existing Go stdio E2E against
the same server sources. Native/server builds, Clippy, quality, semantic lint,
and strict OpenSpec validation remain exact-candidate gates.

## Risks / Trade-offs

- Shared files appear temporarily in #146 until this PR merges; exact source
  equality and a rebase remove the overlap without losing consumer progress.
- All network clients share one open collection; this is a local backend server,
  not separate collection sessions per client.
- Anki remains a Linux x86_64 source build with the approved workspace lock and
  Cargo acquisition requirements. The server owns no third-party source patches.

## Direct workspace import

The user requested direct import after reviewing the provider adapter. Consume
the pinned upstream Cargo workspace and its original source layout, with
standard all_crate_deps dependency maps instead of reflecting a generated
external library's providers and attributes. Only native/backend targets are
compiled for the server. Existing transport and startup contracts remain.

The original dependency approval covers 493 external versions in a conservative
backend closure. Upstream's full workspace Cargo.lock contains 89 additional
versions outside that approval, including the Python bridge and development
tools. out/extract-anki-grpc-project/workspace-dependency-proposal.md owns the
reviewable exact inventory. The user approved the exact 89-version proposal on 2026-10-13. Declare and
acquire only that approved set plus previously approved backend/server versions. Use the upstream locked versions and reuse
the approved server dependencies. Unexpected lock changes require inspection.

Workspace feature unification enables native TLS as well as Rustls. Reuse the
existing pinned OpenSSL source build with a standard openssl-sys annotation
for declared headers and static libraries. Expose its already produced static
outputs and verify the existing C consumer; do not require host development
packages or add a vendored OpenSSL dependency. The user also requested the
project's README and landing metadata to use the finished status.

## Current state and next action

The initial extraction built a byte-identical executable; its upstream closure
contained 24,064 labels without project implementation inputs. CLI/network E2E
failed against that stdio-only server, then passed after adding transports.
Updated E2E failed while initialization still required an RPC, then passed
with startup initialization. Default Go stdio E2E passes using readiness
instead of the removed RPC. Independent SQLite reads confirm expected note fields and generated cards in
both network artifacts. Those observations describe the validated transport candidate before the
direct workspace import. Landing registration now
produces the project index, landing, documentation, and Rust/gRPC taxonomy links.
Delivery receipts own final gates and publication.

The directly imported workspace builds with standard dependency helpers and
the existing pinned OpenSSL. Its native closure contains 28,258 labels,
including build definitions, with no project implementation or Python bridge
inputs. The lock audit covers 599 packages with no unapproved versions. TCP,
Unix, and Go stdio E2E pass; independent SQLite reads verify note fields and
generated cards. Both PR snapshots build a byte-identical server. Clippy and
the existing C OpenSSL consumer pass. README and landing metadata use finished.
Delivery receipts own final exact-candidate gates and publication.

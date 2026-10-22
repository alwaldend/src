## Context

App already delegates owned lifecycle to `collection.Store`; its business services borrow gRPC connections. The existing private CLI interface describes the complete application surface and will move to the shared API package.

## Goals / Non-Goals

Expose per-operation source selection while preserving one unbound App and existing CLI behavior. Keep build as a staged-copy workflow over a base archive path. Do not add backend pooling, registries, dependencies, transactions, or automatic publication for borrowed connections.

## Decisions

Add `api/ankiascodeapi` with matching Go package and Bazel target names. Its `App` interface describes all configuration, collection, state and server workflows. Its small `CollectionSource` interface exposes `Path() string` and `Connection() grpc.ClientConnInterface`; lightweight `CollectionPath` and `CollectionConnection` constructors perform no I/O. Exactly one source must be present, and source descriptions are immutable. The package depends on generated protobuf contracts and gRPC, not native backend code. These shared contracts have one owner and keep CLI independent of the concrete App.

App's collection operations accept the shared source. A private helper preserves path-based `collection.Store` behavior and invokes borrowed operations directly. Borrowed calls never construct or stop native owners and never export a collection package. They preserve cancellation and error causes. Existing error-only mutation results stay unchanged; borrowed owners decide when to publish their package, including whether to enable backend write-back.

Path-source export retains existing source/output alias protections. A borrowed connection carries no source pathname, so its owner must choose outputs distinct from the native collection files; config/state output conflict checks remain enforced. No source discovery RPC or hidden filesystem lookup is introduced.

Move generated schema ownership to `api/ankiascodepb`, rename the schema file to `ankiascode.proto`, use Go package `ankiascodepb`, and change the protobuf namespace from `anki_as_code.collection` to `anki_as_code`. Preserve message names, fields, field numbers and state keys. Update current source, BUILD consumers and documentation; historical completed changes retain their original evidence.

Move the complete CLI workflow contract into `ankiascodeapi.App` and remove the private duplicate. Commands consume that public interface and wrap path flags without opening resources. The typed application variable in cmd checks the concrete App against the public contract during compilation. Preserve trunk's explicit worker construction for path sources; a zero-value App can still use borrowed connections.

Reconcile trunk's server workflow boundary by including `ServerRun` and `ServerServe` in the shared interface. CLI owns their command syntax and forwards opaque native arguments; App delegates through the runner injected by cmd. Context controls server startup only, while the injected native adapter owns hosting and process signal behavior. Keep trunk's package documentation and reconcile its ownership guarantees with borrowed collection lifetimes and renamed schema contracts.

## Risks / Trade-offs

- A source with neither or both alternatives is invalid: reject it before native work.
- Borrowed operations can commit database writes before a later error: preserve native behavior and caller-owned recovery.
- Exporting a package consumes its open collection: only the owner publishes, after App operations finish.
- Renaming protobuf fully qualified names changes descriptors and Any type URLs: the user explicitly authorizes this breaking rename; ordinary field-based TOML/JSON and binary encodings retain their schema fields.
- Moving shared declarations adds one API package: it gives their contracts one owner and avoids App's platform-specific runtime dependency in CLI.

## Failure Coverage Before Implementation

Exercise one borrowed collection through export, state installation, plan, apply, and convergence; verify media, native IDs, persisted output, and source archive preservation. Reuse one App across two open connections and revisit the first. Verify failed and cancelled operations leave caller-owned access intact. Reject nil, empty and conflicting sources. Keep existing path lifecycle, cleanup, cancellation and CLI E2E coverage. Save a final native archive and exported state as repeatable test artifacts.

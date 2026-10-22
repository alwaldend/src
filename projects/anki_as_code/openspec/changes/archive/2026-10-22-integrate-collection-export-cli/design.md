## Context

See proposal.md. Master already includes the prepared-connection exporter, native collection workflow and state commands. The remaining #128 integration can operate without reconciliation.

## Goals / Non-Goals

Expose the merged export capability through the existing CLI while keeping native lifecycle and config projection with their current owners. Plan/apply/build remain placeholders. Do not modify library implementations or add dependencies.

## Decisions

Extract the existing CollectionHandler.Export adapter from #128: resolve paths, reject a destination containing the source, and borrow the shared native collection connection. Construct resource/state stores for the selected namespace, invoke Exporter.Export and return an unchanged flag so archive inputs are never repacked. Keep the existing two required command flags and stderr status.

Adapt the existing config-only and filesystem export E2E cases before production changes. Replace reconciliation-only fixture setup with explicit state CLI operations. Reuse master's native fixture helpers and keep the exporter implementation unchanged. Existing library E2E continues to own detailed projection coverage.

## Risks / Trade-offs

Direct filesystem errors may leave partial output; callers own staging and synchronization. Database opening follows the native backend contract. Export does not initialize state; callers import existing identities explicitly. Keep #128's prerequisite overlap until this independent master-based PR merges.

## Physical path containment

Resolve the source through filesystem symlinks and resolve the destination through its nearest existing ancestor before checking containment. Keep this guard in the app export adapter and leave shared native lifecycle/config APIs unchanged. New destination directories remain supported. Native E2E covers aliases on either side and verifies complete source-tree preservation after rejection.

## Application boundary

Commands consume a private flat application interface. The concrete App owns the scoped registry and coordinates config, collection and state workflows through their services. State methods receive the collection path explicitly. The root accepts only its application interface. Cobra owns parsed flag values; collection methods receive the namespace as an explicit string argument. App retains no command configuration, and a fake application can observe every option without a backend. Remove the unused collection-options object. App owns a scoped collection registry built on the existing native lifecycle owners.

Before implementation, cover command routing, unchanged opaque import IDs, batched bindings, context forwarding, state input/output codecs, required flags, malformed syntax, application error causes and output failures. Reuse native E2E for persistence and cleanup rather than mocking gRPC or duplicating native behavior.

## Test ownership

Keep isolated command dispatch, syntax, codecs and streams with the CLI. Move config parsing, hierarchy and direct export behavior to internal/config and call config.Store directly, without a subprocess. Move collection input failures to pkg/ankiascode and exercise App.ColStateDump. Existing native app E2E retains full workflow and artifact coverage.

## Reusable collection scopes

App methods accept paths and use collection.Store.WithCollection. Callers that
need repeated native operations use the injected connection directly, including
locally constructed state stores and exporters. The collection-store capability owns
source reservation, native lifetime, publication and cleanup. There are no
collection reference IDs or separate lookup operations. Plan, apply and build
remain placeholders in this export integration prerequisite.

## CLI grouping

Keep native workflows and state beneath col, with collection as a Cobra alias. Keep filesystem-only config at the root, with conf as its alias. Register state-key on the collection group; retain operation-specific source options and required-flag checks. Do not retain unreleased ungrouped command spellings.

## Flat application boundary

The concrete App remains in pkg/ankiascode and exposes namespaced workflow
methods. Each collection argument is a path string. internal/cli owns a private
application interface using only strings, protobuf messages and standard-library
types. Its command factories and complete command-tree tests do not import App,
collection lifecycle or native backend implementations. cmd/anki_as_code is the
composition root and passes App to cli.New.

App delegates native access to collection.Store from the independent extraction
merged in PR #165; retain only consumer visibility adjustments required by this
integration. Remove the old App registry, reference interface, file selectors,
scope IDs, handlers and accessors. Existing business services remain authoritative.

## Configuration-set export command

Expose config export [FILE] --output DIRECTORY, with omitted FILE or - reading
stdin. Decode a CollectionConfigSet through the shared codec, defaulting to JSON
and supporting --format toml. App.ConfigExport accepts the parsed object and
output directory, then delegates to config.Store.Export without opening Anki.
Retain direct writes, unchanged input objects, unrelated output entries and
caller-owned valid paths/semantic checks. Schema and read errors precede writes;
filesystem failures can leave partial output.

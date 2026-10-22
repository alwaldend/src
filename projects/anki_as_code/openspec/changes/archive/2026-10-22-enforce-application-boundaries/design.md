## Context

The merged linked-server implementation registers server commands and invokes native hosting in main. The rest of the command tree uses a private application interface. All 17 Go packages across the two projects have doc.go files, but command documentation uses a Command prefix and several ownership exclusions remain implicit.

## Goals / Non-Goals

Route every operation through App, document and review package ownership, and preserve native CLI behavior. Do not change native parsing, process isolation, collection publication, reconciliation or dependencies.

## Decisions

- Inject one synchronous native runner through application options. App builds native argument vectors and owns ServerServe delegation; Cobra only selects the operation and passes context, display name and arguments.
- Keep Rust/cgo linkage and process signal adaptation in the executable's injected adapter. App remains reusable without linking Rust; the startup context is checked before hosting, while native serving owns the calling process until return or termination.
- Keep server command factories under internal/cli, following existing config and collection commands. No handler invokes another handler.
- Audit package comments against imports, effects, call sites and lifecycle. Retain the established lower-level owners rather than moving their domain algorithms into App. doc.go describes package responsibilities and boundaries; the Go skill owns implementation/review instructions, while README and diagrams describe composition.

## Risks / Trade-offs

The injected runner is an explicit optional dependency: configuration and collection workflows do not need native hosting in their process. Installed help, invalid arguments, cancellation before hosting and interrupt behavior need direct verification.

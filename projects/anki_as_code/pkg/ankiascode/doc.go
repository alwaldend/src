// Package ankiascode coordinates all Anki as code application operations.
// App implements ankiascodeapi.App and owns Config, Col, ColState, API and Server
// workflows. All application-owned file access uses the construction-time Afero
// filesystem, defaulting to OS storage; custom paths retain their filesystem's namespace.
// Server and backend-library storage contracts remain independently owned. App
// validates workflow paths,
// composes lower-level services and delegates native hosting through an injected
// runner. ServerRun's context controls startup only; CLI owns subcommand argument prefixes.
// Collection scopes own source reservations and enabled archive publication; export,
// state, resource, plan, reconcile and rebuild retain their domain responsibilities.
// This package does not parse CLI flags, manage process-global signals, link
// Rust/cgo or duplicate native transport/storage logic. The executable supplies
// native adapters. Collection operations select an owned path or borrowed
// connection per call; only path sources require explicit backend configuration.
// Path operations select an address per call or delegate defaults to the injected factory.
// Raw API calls also accept explicit server addresses; ankigrpc owns their client scopes.
// An opt-in address-only constructor enables local in-process execution without
// importing native code here. Any client address can accompany a managed API path;
// network servers must share paths and be reserved for the managed collection scope.
// Borrowed connections retain caller-owned lifetime and package publication.
// ColOpen delegates explicit managed sessions to Store.Open;
// callers borrow the returned client's Connection and call Stop to publish and release it.
// Protocol adapters consume App; this package does not import or host MCP.
package ankiascode

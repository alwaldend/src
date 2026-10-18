// Package ankiascode coordinates all Anki as code application operations.
// App implements ankiascodeapi.App and owns Config, Col, ColState and Server
// workflows, validates workflow paths,
// composes lower-level services and delegates native hosting through an injected
// runner. ServerServe delegates to ServerRun; its context controls startup only.
// Collection scopes own source reservations and archive publication; export,
// state, resource, plan, reconcile and rebuild retain their domain responsibilities.
// This package does not parse CLI flags, manage process-global signals, link
// Rust/cgo or duplicate native transport/storage logic. The executable supplies
// native adapters. Collection operations select an owned path or borrowed
// connection per call; only path sources require explicit worker configuration.
// Borrowed connections retain caller-owned lifetime and package publication.
package ankiascode

// Package transport verifies the standalone native CLI and collection persistence
// over TCP and Unix gRPC using Anki's official protobuf clients. Tests own synthetic
// inputs, worker/client cleanup and closed collection/report artifacts.
// This package supplies no production server, application workflow, client
// lifecycle implementation or storage semantics.
package transport

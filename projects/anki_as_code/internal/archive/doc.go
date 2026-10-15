// Package archive manages private Anki collections and guarded file publication.
// It borrows an existing gRPC connection; its caller owns the server lifetime.
// Anki owns package formats and media handling; offline database bytes stay opaque.
package archive

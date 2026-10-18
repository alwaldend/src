package ankiascodeapi

import "google.golang.org/grpc"

// CollectionSource selects one collection for an application operation.
// Exactly one of Path or Connection must be nonempty. Implementations must
// return stable values throughout an operation. Borrowed connection owners
// coordinate exclusive access, package publication and collection cleanup.
type CollectionSource interface {
	// Path returns a source filename, or an empty string for a borrowed connection.
	Path() string
	// Connection returns an already-open collection connection, or nil for a path.
	Connection() grpc.ClientConnInterface
}

// CollectionPath describes an existing source for application-owned preparation,
// enabled archive publication and cleanup. It performs no filesystem work.
func CollectionPath(path string) CollectionSource {
	return collectionSource{path: path}
}

// CollectionConnection borrows an already-open collection without taking ownership.
// Managed collection workflows never open, close or publish this collection.
// Raw API requests retain their requested native side effects. The owner keeps
// the connection available through each call and coordinates native lifecycle and publication.
func CollectionConnection(connection grpc.ClientConnInterface) CollectionSource {
	return collectionSource{connection: connection}
}

// collectionSource retains immutable source alternatives without native resources.
type collectionSource struct {
	// path identifies a collection whose lifetime the application will own.
	path string
	// connection identifies a collection whose lifetime the caller retains.
	connection grpc.ClientConnInterface
}

// Path returns the described path without resolving it.
func (source collectionSource) Path() string { return source.path }

// Connection returns the borrowed connection without querying its readiness.
func (source collectionSource) Connection() grpc.ClientConnInterface { return source.connection }

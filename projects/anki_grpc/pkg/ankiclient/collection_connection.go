package ankiclient

import (
	"context"
	"fmt"

	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// collectionConnection tracks unary-call completion and native consuming calls
// without taking over the borrowed gRPC transport or serializing ordinary calls.
type collectionConnection struct {
	// ClientConnInterface delegates the unchanged borrowed transport interface.
	grpc.ClientConnInterface
	// owner supplies readiness and serialized consuming-call lifetime.
	owner *ClientCollection
}

// Invoke registers native calls with the owner and retains uncertainty after an
// unconfirmed response. Close/export additionally consume collection ownership.
func (connection *collectionConnection) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("call Anki collection %s: %w", method, err)
	}
	consuming := method == collection.BackendCollectionService_CloseCollection_FullMethodName ||
		method == import_export.BackendImportExportService_ExportCollectionPackage_FullMethodName
	if consuming {
		if err := connection.owner.acquire(ctx); err != nil {
			return fmt.Errorf("consume Anki collection: %w", err)
		}
		defer connection.owner.release()
	}
	if _, err := connection.owner.Connection(); err != nil {
		return fmt.Errorf("call Anki collection %s: %w", method, err)
	}
	connection.owner.lifecycle.Lock()
	if !connection.owner.ready || connection.owner.stopped {
		connection.owner.lifecycle.Unlock()
		return fmt.Errorf("call Anki collection %s: collection is not prepared", method)
	}
	connection.owner.active++
	if consuming {
		connection.owner.ready = false
	}
	connection.owner.lifecycle.Unlock()
	var trailers metadata.MD
	err := connection.ClientConnInterface.Invoke(ctx, method, input, output, append(options, grpc.Trailer(&trailers))...)
	completed := err == nil || len(trailers.Get(ErrorTrailer)) == 1
	connection.owner.lifecycle.Lock()
	connection.owner.uncertain = connection.owner.uncertain || !completed
	if connection.owner.uncertain {
		connection.owner.ready = false
	}
	if consuming && completed {
		connection.owner.opened = false
	}
	connection.owner.active--
	if connection.owner.stopped && connection.owner.active == 0 {
		close(connection.owner.drained)
	}
	connection.owner.lifecycle.Unlock()
	if err != nil {
		return fmt.Errorf("call Anki collection %s: %w", method, err)
	}
	return nil
}

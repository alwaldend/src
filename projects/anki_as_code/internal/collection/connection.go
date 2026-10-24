package collection

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	ankicol "github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/import_export"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// scopeConnection observes borrowed calls that can still be using owned restore files.
type scopeConnection struct {
	// ClientConnInterface delegates native RPCs to the selected transport.
	grpc.ClientConnInterface
	// owner admits calls and retains native completion uncertainty.
	owner *scope
}

// Invoke registers unary work until its response, including native consuming calls.
func (connection *scopeConnection) Invoke(ctx context.Context, method string, input, output any, options ...grpc.CallOption) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("call Anki collection %s: %w", method, err)
	}
	owner := connection.owner
	owner.mutex.Lock()
	if owner.closing || !owner.opened || owner.uncertain {
		owner.mutex.Unlock()
		return fmt.Errorf("call Anki collection %s: collection is not open", method)
	}
	owner.active++
	if method == ankicol.BackendCollectionService_CloseCollection_FullMethodName || method == import_export.BackendImportExportService_ExportCollectionPackage_FullMethodName {
		owner.opened = false
	}
	owner.mutex.Unlock()
	var trailers metadata.MD
	err := connection.ClientConnInterface.Invoke(ctx, method, input, output, append(options, grpc.Trailer(&trailers))...)
	owner.mutex.Lock()
	owner.uncertain = owner.uncertain || (err != nil && len(trailers.Get(ankiclient.ErrorTrailer)) != 1)
	owner.active--
	if owner.closing && owner.active == 0 {
		close(owner.drained)
	}
	owner.mutex.Unlock()
	if err != nil {
		return fmt.Errorf("call Anki collection %s: %w", method, err)
	}
	return nil
}

package native

import (
	"fmt"
	"net/url"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
)

// init registers direct native constructors without initializing backend state.
func init() {
	factory := ankigrpc.ToClientFactory(newInprocessClient)
	ankigrpc.RegisterClientFactory("inprocess", factory)
	ankigrpc.RegisterClientFactory("inprocess+cgo", factory)
}

// newInprocessClient creates fresh native state on Start; labels never share clients.
func newInprocessClient(address string) (*Client, error) {
	if _, err := url.Parse(address); err != nil {
		return nil, fmt.Errorf("parse in-process address: %w", err)
	}
	return New(), nil
}

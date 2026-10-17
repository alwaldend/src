package app

import (
	"context"
	"errors"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/grpc"
)

// App coordinates complete workflows through scoped collection services.
// Its zero value is ready to use and must not be copied after collection access.
type App struct {
	// collections reserves sources and publishes successful changed archive scopes.
	collections collection.Store
}

// Export writes an editable text representation of a collection.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Export(ctx context.Context, input, output string) error {
	return errors.New("export is not implemented")
}

// Plan describes changes needed to reconcile a collection with its text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Plan(ctx context.Context, input, text string) (any, error) {
	return nil, errors.New("plan is not implemented")
}

// Apply reconciles an offline collection with its text representation.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Apply(ctx context.Context, input, text string) error {
	return errors.New("apply is not implemented")
}

// Build creates an archive from a base collection and desired text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Build(ctx context.Context, base, text, output string) error {
	return errors.New("build is not implemented")
}

// Dump loads and encodes config files without opening Anki.
func (application *App) Dump(ctx context.Context, input, format string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	data, err := config.New(input).Dump(format)
	if err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("finish configuration dump: %w", err)
	}
	return data, nil
}

// ColStateDump reads stored bindings and closes the collection without publishing changes.
func (application *App) ColStateDump(ctx context.Context, input, key string) (*collectionpb.State, error) {
	var value *collectionpb.State
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		var err error
		value, err = state.New(connection, key).Dump(ctx)
		if err != nil {
			return false, fmt.Errorf("read collection state: %w", err)
		}
		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("dump collection state: %w", err)
	}
	return value, nil
}

// ColStateImport loads configs and imports a complete binding batch within one native scope.
func (application *App) ColStateImport(ctx context.Context, input, text, key string, bindings map[string]string) error {
	configs, err := config.New(text).Load()
	if err != nil {
		return fmt.Errorf("load import configs: %w", err)
	}
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, key).Import(ctx, configs, bindings)
		if err != nil {
			return false, fmt.Errorf("import collection bindings: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("import collection state: %w", err)
	}
	return nil
}

// ColStateUpdate merges external bindings and publishes a changed package after scope success.
func (application *App) ColStateUpdate(ctx context.Context, input, key string, value *collectionpb.State) error {
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, key).Update(ctx, value)
		if err != nil {
			return false, fmt.Errorf("merge collection bindings: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("update collection state: %w", err)
	}
	return nil
}

// ColStateReplace replaces all stored bindings and publishes changed packages after scope success.
func (application *App) ColStateReplace(ctx context.Context, input, key string, value *collectionpb.State) error {
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, key).Replace(ctx, value)
		if err != nil {
			return false, fmt.Errorf("replace collection bindings: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("replace collection state: %w", err)
	}
	return nil
}

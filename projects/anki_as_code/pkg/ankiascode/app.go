package ankiascode

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/grpc"
)

// App coordinates configuration and native collection workflows through their services.
// Its zero value is ready to use and must not be copied after collection access.
type App struct {
	// collections reserves sources and publishes successful changed archive scopes.
	collections collection.Store
}

// ConfigDump loads and encodes config files without opening Anki.
func (application *App) ConfigDump(ctx context.Context, input, format string) ([]byte, error) {
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

// ConfigExport writes supplied configs as TOML without opening Anki or validating relationships.
// Declared files are overwritten; other entries are untouched and errors may leave partial output.
func (application *App) ConfigExport(ctx context.Context, output string, configs *collectionpb.CollectionConfigSet) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("export configuration: %w", err)
	}
	output, err := expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve configuration output: %w", err)
	}
	if err := config.New(output).Export(configs); err != nil {
		return fmt.Errorf("export configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish configuration export: %w", err)
	}
	return nil
}

// ColExport writes editable TOML without exporting state or changing the source collection.
func (application *App) ColExport(ctx context.Context, input, output, stateKey string) error {
	output, err := expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	source, err := expandPath(input)
	if err != nil {
		return fmt.Errorf("resolve source collection: %w", err)
	}
	source, err = filepath.EvalSymlinks(source)
	if err != nil {
		return fmt.Errorf("resolve export source aliases: %w", err)
	}
	output, err = resolveExportOutputPath(output)
	if err != nil {
		return fmt.Errorf("resolve export destination: %w", err)
	}
	relative, err := filepath.Rel(output, source)
	if err != nil {
		return fmt.Errorf("compare export and collection paths: %w", err)
	}
	if filepath.IsLocal(relative) {
		return fmt.Errorf("export destination must not contain the source collection")
	}
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		exporter := export.New(resource.New(connection), state.New(connection, stateKey))
		if err := exporter.Export(ctx, output); err != nil {
			return false, fmt.Errorf("export prepared collection: %w", err)
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("export collection: %w", err)
	}
	return nil
}

// ColPlan describes changes needed to reconcile a collection with its text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) ColPlan(ctx context.Context, input, text, stateKey string) (*collectionpb.CollectionChangePlan, error) {
	return nil, errors.New("plan is not implemented")
}

// ColApply reconciles an offline collection with its text representation.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) ColApply(ctx context.Context, input, text, stateKey string) error {
	return errors.New("apply is not implemented")
}

// ColBuild creates an archive from a base collection and desired text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) ColBuild(ctx context.Context, base, text, output, stateKey string) error {
	return errors.New("build is not implemented")
}

// ColStateDump reads stored bindings and closes the selected collection without publication.
func (application *App) ColStateDump(ctx context.Context, input, stateKey string) (*collectionpb.State, error) {
	var value *collectionpb.State
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		var err error
		value, err = state.New(connection, stateKey).Dump(ctx)
		if err != nil {
			return false, fmt.Errorf("read collection state: %w", err)
		}
		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("dump state: %w", err)
	}
	return value, nil
}

// ColStateImport loads configs and binds a complete batch using the selected collection.
func (application *App) ColStateImport(ctx context.Context, input, text, stateKey string, bindings map[string]string) error {
	configs, err := config.New(text).Load()
	if err != nil {
		return fmt.Errorf("load import configs: %w", err)
	}
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, stateKey).Import(ctx, configs, bindings)
		if err != nil {
			return false, fmt.Errorf("bind collection objects: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("import state: %w", err)
	}
	return nil
}

// ColStateUpdate merges bindings and publishes changed archives after scope success.
func (application *App) ColStateUpdate(ctx context.Context, input, stateKey string, value *collectionpb.State) error {
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, stateKey).Update(ctx, value)
		if err != nil {
			return false, fmt.Errorf("merge collection state: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	return nil
}

// ColStateReplace installs bindings and publishes changed archives after scope success.
func (application *App) ColStateReplace(ctx context.Context, input, stateKey string, value *collectionpb.State) error {
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		changed, err := state.New(connection, stateKey).Replace(ctx, value)
		if err != nil {
			return false, fmt.Errorf("replace collection state: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("replace state: %w", err)
	}
	return nil
}

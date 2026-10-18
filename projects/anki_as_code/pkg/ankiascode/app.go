package ankiascode

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	planning "git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rebuild"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/reconcile"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// App coordinates configuration and native collection workflows through their services.
// Its zero value supports configuration operations. Use New to configure a native
// worker before collection access; it must not be copied after collection access.
type App struct {
	// collections reserves sources and publishes successful changed archive scopes.
	collections collection.Store
}

// Options configures an application's native worker without starting it.
type Options struct {
	// BackendBinary supplies the executable required for collection scopes; no default is discovered.
	BackendBinary string
	// BackendArguments supplies its server-mode arguments.
	BackendArguments []string
}

// New constructs an application with explicit native worker configuration.
func New(options Options) *App {
	return &App{collections: collection.New(ankiclient.BackendSTDOptions{Binary: options.BackendBinary, Arguments: options.BackendArguments})}
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

// ColExport optionally writes snapshot bindings as TOML after exporting configs.
// The source collection and its stored state remain unchanged; stateOutput may be empty.
func (application *App) ColExport(ctx context.Context, input, output, stateOutput, stateKey string) error {
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
	output, err = export.ResolveOutputPath(output)
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
	if stateOutput != "" {
		stateOutput, err = expandPath(stateOutput)
		if err != nil {
			return fmt.Errorf("resolve state output: %w", err)
		}
		stateOutput, err = export.ResolveOutputPath(stateOutput)
		if err != nil {
			return fmt.Errorf("resolve state output aliases: %w", err)
		}
		sourceInfo, err := os.Stat(source)
		if err != nil {
			return fmt.Errorf("inspect source collection: %w", err)
		}
		if info, err := os.Stat(stateOutput); err == nil {
			if !info.Mode().IsRegular() || os.SameFile(sourceInfo, info) {
				return fmt.Errorf("state output must be a regular file distinct from the source collection")
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect state output: %w", err)
		}
	}
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		exporter := export.New(resource.New(connection), state.New(connection, stateKey))
		if err := exporter.Export(ctx, output, stateOutput); err != nil {
			return false, fmt.Errorf("export prepared collection: %w", err)
		}
		return false, nil
	}); err != nil {
		return fmt.Errorf("export collection: %w", err)
	}
	return nil
}

// ColPlan compares declarations with an API snapshot without native writes.
func (application *App) ColPlan(ctx context.Context, input, dir, stateKey string) (*collectionpb.CollectionChangePlan, error) {
	var report *collectionpb.CollectionChangePlan
	if err := application.withPlan(ctx, input, dir, stateKey, func(_ *reconcile.Reconciler, planned *collectionpb.CollectionPlan) (bool, error) {
		report = planned.Report
		return false, nil
	}); err != nil {
		return nil, fmt.Errorf("plan collection: %w", err)
	}
	return report, nil
}

// ColBuild delegates staged copy-and-apply and atomic publication to rebuild.Builder.
func (application *App) ColBuild(ctx context.Context, base, dir, output, stateKey string) error {
	source, err := expandPath(base)
	if err != nil {
		return fmt.Errorf("resolve build source: %w", err)
	}
	output, err = expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve build output: %w", err)
	}
	builder := rebuild.New(func(ctx context.Context, archive string) error {
		if err := application.ColApply(ctx, archive, dir, stateKey); err != nil {
			return fmt.Errorf("apply build configs: %w", err)
		}
		return nil
	})
	if err := builder.Build(ctx, source, output); err != nil {
		return fmt.Errorf("build collection: %w", err)
	}
	return nil
}

// ColApply reconciles an offline database or archive in place.
func (application *App) ColApply(ctx context.Context, input, dir, stateKey string) error {
	if err := application.withPlan(ctx, input, dir, stateKey, func(reconciler *reconcile.Reconciler, planned *collectionpb.CollectionPlan) (bool, error) {
		changed, err := reconciler.Apply(ctx, planned)
		if err != nil {
			return false, fmt.Errorf("apply desired state: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("apply collection: %w", err)
	}
	return nil
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

// withPlan loads configs before borrowing one collection and prepares a plan for its operation.
// The callback shares the captured resource/state clients and controls archive publication.
func (application *App) withPlan(ctx context.Context, input, dir, stateKey string, operation func(*reconcile.Reconciler, *collectionpb.CollectionPlan) (bool, error)) error {
	configs, err := config.New(dir).Load()
	if err != nil {
		return fmt.Errorf("load desired configs: %w", err)
	}
	if err := application.collections.WithCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
		resources, states := resource.New(connection), state.New(connection, stateKey)
		snapshot, err := resources.Snapshot(ctx)
		if err != nil {
			return false, fmt.Errorf("capture resources: %w", err)
		}
		stored, err := states.Dump(ctx)
		if err != nil {
			return false, fmt.Errorf("capture bindings: %w", err)
		}
		planned, err := planning.Build(configs, snapshot, stored)
		if err != nil {
			return false, fmt.Errorf("build execution plan: %w", err)
		}
		changed, err := operation(reconcile.New(resources, states), planned)
		if err != nil {
			return false, fmt.Errorf("run planned operation: %w", err)
		}
		return changed, nil
	}); err != nil {
		return fmt.Errorf("prepare collection plan: %w", err)
	}
	return nil
}

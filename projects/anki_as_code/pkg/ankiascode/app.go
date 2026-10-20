package ankiascode

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	planning "git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rebuild"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/reconcile"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/rpc"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"google.golang.org/grpc"
)

// App coordinates configuration, collection and server operations through their owners.
// Its zero value supports configuration operations and borrowed connections. Use
// New to configure a native backend for path sources; do not copy App after use.
type App struct {
	// clientFactory supplies the same address-selected constructors to API and collection scopes.
	clientFactory ankigrpc.ClientFactory
	// collections reserves sources and owns managed archive publication.
	collections collection.Store
	// serverRunner hosts the server CLI through the adapter supplied by the executable.
	serverRunner func([]string) error
}

// Options supplies client construction, cleanup policy and an optional server hosting adapter.
type Options struct {
	// NewClient constructs clients and selects its default when an operation address is empty.
	NewClient ankigrpc.ClientFactory
	// ServerRunner hosts a full server argument vector, including its program name.
	// It owns serving until return or process termination; collection operations do not use it.
	ServerRunner func([]string) error
	// CollectionCloseTimeout bounds publication and close independently of request cancellation; nonpositive values select ten minutes.
	CollectionCloseTimeout time.Duration
}

// New constructs an application from explicit backend and server hosting dependencies.
func New(options Options) *App {
	return &App{
		clientFactory: options.NewClient,
		collections: collection.New(collection.Options{
			NewClient: options.NewClient, CloseTimeout: options.CollectionCloseTimeout,
		}),
		serverRunner: options.ServerRunner,
	}
}

// ServerRun invokes the configured server CLI with a display name and unchanged arguments.
// ctx controls entry only; the injected runner owns hosting and its process until it returns.
func (application *App) ServerRun(ctx context.Context, program string, arguments []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start server CLI: %w", errors.Join(err, context.Cause(ctx)))
	}
	if application.serverRunner == nil {
		return fmt.Errorf("server CLI runner is not configured")
	}
	if err := application.serverRunner(append([]string{program}, arguments...)); err != nil {
		return fmt.Errorf("run server CLI: %w", err)
	}
	return nil
}

// ServerServe delegates to ServerRun with the server serve subcommand prepended.
func (application *App) ServerServe(ctx context.Context, program string, arguments []string) error {
	if err := application.ServerRun(ctx, program, append([]string{"serve"}, arguments...)); err != nil {
		return fmt.Errorf("serve Anki: %w", err)
	}
	return nil
}

// APICall validates and invokes one upstream method before returning its encoded response.
// Native cleanup and enabled managed-archive publication finish before success;
// borrowed connections and remote collection lifetime remain with their owners.
func (application *App) APICall(ctx context.Context, request ankiascodeapi.APIRequest) (output []byte, err error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("start Anki API call: %w", errors.Join(err, context.Cause(ctx)))
	}
	options := request.Collection
	var path string
	var connection grpc.ClientConnInterface
	if options.Address != "" {
		if !isNilConnection(options.Connection) {
			return nil, fmt.Errorf("API access requires exactly one path, server address or open connection")
		}
		path = options.Path
	} else {
		path, connection, err = resolveCollectionOptions(options)
		if err != nil {
			return nil, fmt.Errorf("select API collection: %w", err)
		}
	}
	method, err := rpc.Resolve(request.Method)
	if err != nil {
		return nil, fmt.Errorf("resolve API method: %w", err)
	}
	if path != "" {
		switch method.Descriptor.Name() {
		case "OpenCollection", "CloseCollection", "ImportCollectionPackage":
			return nil, fmt.Errorf("%s requires connected access because managed paths own collection lifetime", method.Path)
		}
	}
	if request.InputFormat == "" {
		request.InputFormat = "json"
	}
	if request.OutputFormat == "" {
		request.OutputFormat = "json"
	}
	operation, err := rpc.Prepare(method, request.Input, request.InputFormat, request.OutputFormat)
	if err != nil {
		return nil, fmt.Errorf("prepare Anki API request: %w", err)
	}
	run := func(connection grpc.ClientConnInterface) error {
		var err error
		output, err = operation(ctx, connection)
		if err != nil {
			return fmt.Errorf("call prepared Anki API: %w", errors.Join(err, ctx.Err(), context.Cause(ctx)))
		}
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("finish Anki API call: %w", errors.Join(err, context.Cause(ctx)))
		}
		return nil
	}
	if options.Address != "" && path == "" {
		if err := ankigrpc.WithClient(ctx, application.clientFactory, options.Address, run); err != nil {
			return nil, fmt.Errorf("call addressed Anki API: %w", err)
		}
	} else if connection != nil {
		if err := run(connection); err != nil {
			return nil, fmt.Errorf("call connected Anki API: %w", err)
		}
	} else if err := application.collections.WithCollection(ctx, collection.CollectionOptions{Path: path, Address: options.Address, WriteBack: request.WriteBack}, run); err != nil {
		return nil, fmt.Errorf("call offline Anki API: %w", err)
	}
	return output, nil
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
func (application *App) ConfigExport(ctx context.Context, output string, configs *ankiascodepb.CollectionConfigSet) error {
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
// Collection resources and stored state remain unchanged; managed archives may be repackaged.
// stateOutput may be empty.
// Borrowed connection owners must keep outputs distinct from native collection files.
func (application *App) ColExport(ctx context.Context, options ankiascodeapi.CollectionOptions, output, stateOutput string) error {
	path, _, err := resolveCollectionOptions(options)
	if err != nil {
		return fmt.Errorf("resolve export source: %w", err)
	}
	output, err = expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}
	output, err = export.ResolveOutputPath(output)
	if err != nil {
		return fmt.Errorf("resolve export destination: %w", err)
	}
	var source string
	if path != "" {
		source, err = expandPath(path)
		if err != nil {
			return fmt.Errorf("resolve source collection: %w", err)
		}
		source, err = filepath.EvalSymlinks(source)
		if err != nil {
			return fmt.Errorf("resolve export source aliases: %w", err)
		}
		relative, err := filepath.Rel(output, source)
		if err != nil {
			return fmt.Errorf("compare export and collection paths: %w", err)
		}
		if filepath.IsLocal(relative) {
			return fmt.Errorf("export destination must not contain the source collection")
		}
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
		var sourceInfo os.FileInfo
		if source != "" {
			sourceInfo, err = os.Stat(source)
			if err != nil {
				return fmt.Errorf("inspect source collection: %w", err)
			}
		}
		if info, err := os.Stat(stateOutput); err == nil {
			if !info.Mode().IsRegular() || sourceInfo != nil && os.SameFile(sourceInfo, info) {
				return fmt.Errorf("state output must be a regular file distinct from the source collection")
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect state output: %w", err)
		}
	}
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		exporter := export.New(resource.New(connection), state.New(connection, options.StateKey))
		if err := exporter.Export(ctx, output, stateOutput); err != nil {
			return fmt.Errorf("export prepared collection: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("export collection: %w", err)
	}
	return nil
}

// ColPlan compares declarations with an API snapshot without native writes.
func (application *App) ColPlan(ctx context.Context, options ankiascodeapi.CollectionOptions, dir string) (*ankiascodepb.CollectionChangePlan, error) {
	var report *ankiascodepb.CollectionChangePlan
	if err := application.withPlan(ctx, options, dir, func(_ *reconcile.Reconciler, planned *ankiascodepb.CollectionPlan) error {
		report = planned.Report
		return nil
	}); err != nil {
		return nil, fmt.Errorf("plan collection: %w", err)
	}
	return report, nil
}

// ColBuild delegates staged copy-and-apply and atomic publication to rebuild.Builder.
func (application *App) ColBuild(ctx context.Context, options ankiascodeapi.CollectionOptions, dir, output string) error {
	path, connection, err := resolveCollectionOptions(options)
	if err != nil {
		return fmt.Errorf("resolve build collection: %w", err)
	}
	if connection != nil {
		return fmt.Errorf("build requires a collection path")
	}
	source, err := expandPath(path)
	if err != nil {
		return fmt.Errorf("resolve build source: %w", err)
	}
	output, err = expandPath(output)
	if err != nil {
		return fmt.Errorf("resolve build output: %w", err)
	}
	builder := rebuild.New(func(ctx context.Context, archive string) error {
		staged := options
		staged.Path = archive
		if err := application.ColApply(ctx, staged, dir); err != nil {
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
func (application *App) ColApply(ctx context.Context, options ankiascodeapi.CollectionOptions, dir string) error {
	if err := application.withPlan(ctx, options, dir, func(reconciler *reconcile.Reconciler, planned *ankiascodepb.CollectionPlan) error {
		if _, err := reconciler.Apply(ctx, planned); err != nil {
			return fmt.Errorf("apply desired state: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("apply collection: %w", err)
	}
	return nil
}

// ColStateDump reads stored bindings; path sources can repackage unchanged content.
func (application *App) ColStateDump(ctx context.Context, options ankiascodeapi.CollectionOptions) (*ankiascodepb.State, error) {
	var value *ankiascodepb.State
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		var err error
		value, err = state.New(connection, options.StateKey).Dump(ctx)
		if err != nil {
			return fmt.Errorf("read collection state: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("dump state: %w", err)
	}
	return value, nil
}

// ColStateImport loads configs and binds a complete batch using the selected collection.
func (application *App) ColStateImport(ctx context.Context, options ankiascodeapi.CollectionOptions, text string, bindings map[string]string) error {
	configs, err := config.New(text).Load()
	if err != nil {
		return fmt.Errorf("load import configs: %w", err)
	}
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		if _, err := state.New(connection, options.StateKey).Import(ctx, configs, bindings); err != nil {
			return fmt.Errorf("bind collection objects: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("import state: %w", err)
	}
	return nil
}

// ColStateUpdate merges bindings; enabled archive publication is independent of operation results.
func (application *App) ColStateUpdate(ctx context.Context, options ankiascodeapi.CollectionOptions, value *ankiascodepb.State) error {
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		if _, err := state.New(connection, options.StateKey).Update(ctx, value); err != nil {
			return fmt.Errorf("merge collection state: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("update state: %w", err)
	}
	return nil
}

// ColStateReplace installs bindings; enabled archive publication is independent of operation results.
func (application *App) ColStateReplace(ctx context.Context, options ankiascodeapi.CollectionOptions, value *ankiascodepb.State) error {
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		if _, err := state.New(connection, options.StateKey).Replace(ctx, value); err != nil {
			return fmt.Errorf("replace collection state: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("replace state: %w", err)
	}
	return nil
}

// withPlan loads configs before borrowing one collection and prepares a plan for its operation.
// The callback shares the captured resource/state clients while the collection backend owns publication.
func (application *App) withPlan(ctx context.Context, options ankiascodeapi.CollectionOptions, dir string, operation func(*reconcile.Reconciler, *ankiascodepb.CollectionPlan) error) error {
	configs, err := config.New(dir).Load()
	if err != nil {
		return fmt.Errorf("load desired configs: %w", err)
	}
	if err := application.withCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
		resources, states := resource.New(connection), state.New(connection, options.StateKey)
		snapshot, err := resources.Snapshot(ctx)
		if err != nil {
			return fmt.Errorf("capture resources: %w", err)
		}
		stored, err := states.Dump(ctx)
		if err != nil {
			return fmt.Errorf("capture bindings: %w", err)
		}
		planned, err := planning.Build(configs, snapshot, stored)
		if err != nil {
			return fmt.Errorf("build execution plan: %w", err)
		}
		if err := operation(reconcile.New(resources, states), planned); err != nil {
			return fmt.Errorf("run planned operation: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("prepare collection plan: %w", err)
	}
	return nil
}

// withCollection manages path sources and borrows open connections without changing their lifetime.
func (application *App) withCollection(ctx context.Context, options ankiascodeapi.CollectionOptions, operation func(grpc.ClientConnInterface) error) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start collection operation: %w", errors.Join(err, context.Cause(ctx)))
	}
	path, connection, err := resolveCollectionOptions(options)
	if err != nil {
		return fmt.Errorf("select collection: %w", err)
	}
	if connection == nil {
		if err := application.collections.WithCollection(ctx, collection.CollectionOptions{Path: path, Address: options.Address, WriteBack: true}, operation); err != nil {
			return fmt.Errorf("run managed collection: %w", err)
		}
		return nil
	}
	if err := operation(connection); err != nil {
		return fmt.Errorf("run borrowed collection: %w", errors.Join(err, ctx.Err(), context.Cause(ctx)))
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish borrowed collection: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}

// resolveCollectionOptions validates managed and borrowed access without opening native resources.
func resolveCollectionOptions(options ankiascodeapi.CollectionOptions) (string, grpc.ClientConnInterface, error) {
	connection := options.Connection
	if isNilConnection(connection) {
		connection = nil
	}
	if (options.Path == "") == (connection == nil) || connection != nil && options.Address != "" {
		return "", nil, fmt.Errorf("collection access requires exactly one path or open connection; borrowed connections cannot include an address")
	}
	return options.Path, connection, nil
}

// isNilConnection recognizes absent connections boxed in interfaces.
func isNilConnection(connection grpc.ClientConnInterface) bool {
	if connection == nil {
		return true
	}
	reflected := reflect.ValueOf(connection)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

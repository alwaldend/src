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
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"google.golang.org/grpc"
)

// App coordinates configuration, collection and server operations through their owners.
// Its zero value supports configuration operations and borrowed connections. Use
// New to configure a native worker for path sources; do not copy App after use.
type App struct {
	// collections reserves sources and owns managed archive publication.
	collections collection.Store
	// serverRunner hosts the native CLI through the adapter supplied by the executable.
	serverRunner func([]string) error
}

// Options supplies explicit worker settings and an optional native hosting adapter.
type Options struct {
	// ServerRunner hosts a full native argument vector, including its program name.
	// It owns serving until return or process termination; collection operations do not use it.
	ServerRunner func([]string) error
	// BackendBinary supplies the executable required for collection scopes; no default is discovered.
	BackendBinary string
	// BackendArguments supplies its server-mode arguments.
	BackendArguments []string
	// CollectionCloseTimeout bounds publication and close independently of request cancellation; nonpositive values select ten minutes.
	CollectionCloseTimeout time.Duration
}

// New constructs an application from explicit worker and native hosting dependencies.
func New(options Options) *App {
	return &App{
		collections:  collection.New(collection.Options{Backend: ankiclient.BackendSTDOptions{Binary: options.BackendBinary, Arguments: options.BackendArguments}, CloseTimeout: options.CollectionCloseTimeout}),
		serverRunner: options.ServerRunner,
	}
}

// ServerRun invokes the configured native CLI with a display name and unchanged arguments.
// ctx controls entry only; the injected runner owns hosting and its process until it returns.
func (application *App) ServerRun(ctx context.Context, program string, arguments []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start native CLI: %w", errors.Join(err, context.Cause(ctx)))
	}
	if application.serverRunner == nil {
		return fmt.Errorf("native CLI runner is not configured")
	}
	if err := application.serverRunner(append([]string{program}, arguments...)); err != nil {
		return fmt.Errorf("run native CLI: %w", err)
	}
	return nil
}

// ServerServe delegates to ServerRun with the native serve subcommand prepended.
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
	hasCollection := !isNilCollectionValue(request.Collection)
	if hasCollection == (request.Backend != "") {
		return nil, fmt.Errorf("API access requires exactly one collection source or backend address")
	}
	var path string
	var connection grpc.ClientConnInterface
	if hasCollection {
		path, connection, err = resolveCollectionSource(request.Collection)
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
	if request.Backend != "" {
		remote := ankiclient.NewGRPC(ankiclient.BackendGRPCOptions{Address: request.Backend})
		defer func() {
			if stopErr := remote.Stop(context.WithoutCancel(ctx)); stopErr != nil {
				output = nil
				err = errors.Join(err, fmt.Errorf("close remote Anki connection: %w", stopErr))
			}
		}()
		if err := remote.Start(ctx); err != nil {
			return nil, fmt.Errorf("start remote Anki connection: %w", err)
		}
		connection, err = remote.Connection()
		if err != nil {
			return nil, fmt.Errorf("get remote Anki connection: %w", err)
		}
	}
	if connection != nil {
		if err := run(connection); err != nil {
			return nil, fmt.Errorf("call connected Anki API: %w", err)
		}
	} else if err := application.collections.WithCollection(ctx, collection.CollectionOptions{Path: path, WriteBack: request.WriteBack}, func(connection grpc.ClientConnInterface) (bool, error) {
		if err := run(connection); err != nil {
			return false, fmt.Errorf("call offline Anki API: %w", err)
		}
		return false, nil
	}); err != nil {
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
// The source collection and its stored state remain unchanged; stateOutput may be empty.
// Borrowed connection owners must keep outputs distinct from native collection files.
func (application *App) ColExport(ctx context.Context, input ankiascodeapi.CollectionSource, output, stateOutput, stateKey string) error {
	path, _, err := resolveCollectionSource(input)
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
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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
func (application *App) ColPlan(ctx context.Context, input ankiascodeapi.CollectionSource, dir, stateKey string) (*ankiascodepb.CollectionChangePlan, error) {
	var report *ankiascodepb.CollectionChangePlan
	if err := application.withPlan(ctx, input, dir, stateKey, func(_ *reconcile.Reconciler, planned *ankiascodepb.CollectionPlan) (bool, error) {
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
		if err := application.ColApply(ctx, ankiascodeapi.CollectionPath(archive), dir, stateKey); err != nil {
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
func (application *App) ColApply(ctx context.Context, input ankiascodeapi.CollectionSource, dir, stateKey string) error {
	if err := application.withPlan(ctx, input, dir, stateKey, func(reconciler *reconcile.Reconciler, planned *ankiascodepb.CollectionPlan) (bool, error) {
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

// ColStateDump reads stored bindings; path sources can repackage unchanged content.
func (application *App) ColStateDump(ctx context.Context, input ankiascodeapi.CollectionSource, stateKey string) (*ankiascodepb.State, error) {
	var value *ankiascodepb.State
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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
func (application *App) ColStateImport(ctx context.Context, input ankiascodeapi.CollectionSource, text, stateKey string, bindings map[string]string) error {
	configs, err := config.New(text).Load()
	if err != nil {
		return fmt.Errorf("load import configs: %w", err)
	}
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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

// ColStateUpdate merges bindings; enabled archive publication is independent of operation results.
func (application *App) ColStateUpdate(ctx context.Context, input ankiascodeapi.CollectionSource, stateKey string, value *ankiascodepb.State) error {
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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

// ColStateReplace installs bindings; enabled archive publication is independent of operation results.
func (application *App) ColStateReplace(ctx context.Context, input ankiascodeapi.CollectionSource, stateKey string, value *ankiascodepb.State) error {
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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
// The callback shares the captured resource/state clients while the collection backend owns publication.
func (application *App) withPlan(ctx context.Context, input ankiascodeapi.CollectionSource, dir, stateKey string, operation func(*reconcile.Reconciler, *ankiascodepb.CollectionPlan) (bool, error)) error {
	configs, err := config.New(dir).Load()
	if err != nil {
		return fmt.Errorf("load desired configs: %w", err)
	}
	if err := application.withCollection(ctx, input, func(connection grpc.ClientConnInterface) (bool, error) {
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

// withCollection manages path sources and borrows open connections without changing their lifetime.
func (application *App) withCollection(ctx context.Context, source ankiascodeapi.CollectionSource, operation func(grpc.ClientConnInterface) (bool, error)) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("start collection operation: %w", errors.Join(err, context.Cause(ctx)))
	}
	path, connection, err := resolveCollectionSource(source)
	if err != nil {
		return fmt.Errorf("select collection: %w", err)
	}
	if connection == nil {
		if err := application.collections.WithCollection(ctx, collection.CollectionOptions{Path: path, WriteBack: true}, operation); err != nil {
			return fmt.Errorf("run managed collection: %w", err)
		}
		return nil
	}
	if _, err := operation(connection); err != nil {
		return fmt.Errorf("run borrowed collection: %w", errors.Join(err, ctx.Err(), context.Cause(ctx)))
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish borrowed collection: %w", errors.Join(err, context.Cause(ctx)))
	}
	return nil
}

// resolveCollectionSource validates source alternatives without opening native resources.
func resolveCollectionSource(source ankiascodeapi.CollectionSource) (string, grpc.ClientConnInterface, error) {
	if isNilCollectionValue(source) {
		return "", nil, fmt.Errorf("collection source requires exactly one path or open connection")
	}
	path, connection := source.Path(), source.Connection()
	if isNilCollectionValue(connection) {
		connection = nil
	}
	if (path == "") == (connection == nil) {
		return "", nil, fmt.Errorf("collection source requires exactly one path or open connection")
	}
	return path, connection, nil
}

// isNilCollectionValue recognizes absent sources and connections boxed in interfaces.
func isNilCollectionValue(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

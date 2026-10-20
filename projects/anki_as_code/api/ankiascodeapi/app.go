package ankiascodeapi

import (
	"context"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"google.golang.org/grpc"
)

// App defines configuration, collection, state, API and server workflows for callers
// and command implementations without depending on a concrete application or
// native worker. Collection operations select their source per call; borrowed
// connections retain caller-owned lifetime and package publication.
type App interface {
	// APICall validates and invokes one upstream RPC, returning its encoded response after owned cleanup.
	APICall(ctx context.Context, request APIRequest) ([]byte, error)

	// ServerRun hosts opaque native CLI arguments through the configured adapter.
	// The context controls startup only; the adapter owns hosting until return.
	ServerRun(ctx context.Context, program string, arguments []string) error
	// ServerServe delegates to ServerRun with the native serve subcommand prepended.
	ServerServe(ctx context.Context, program string, arguments []string) error
	// ConfigDump returns encoded collection configurations.
	ConfigDump(ctx context.Context, text, format string) ([]byte, error)
	// ConfigExport writes a parsed configuration set as a TOML tree.
	ConfigExport(ctx context.Context, output string, configs *ankiascodepb.CollectionConfigSet) error
	// ColExport writes editable configs and optional snapshot bindings.
	ColExport(ctx context.Context, options CollectionOptions, output, stateOutput string) error
	// ColPlan describes proposed collection changes.
	ColPlan(ctx context.Context, options CollectionOptions, text string) (*ankiascodepb.CollectionChangePlan, error)
	// ColApply reconciles the supplied collection in place.
	ColApply(ctx context.Context, options CollectionOptions, text string) error
	// ColBuild stages a base archive path and writes an independent reconciled archive.
	ColBuild(ctx context.Context, options CollectionOptions, text, output string) error
	// ColStateDump reads stored collection bindings.
	ColStateDump(ctx context.Context, options CollectionOptions) (*ankiascodepb.State, error)
	// ColStateImport binds existing objects to configuration paths.
	ColStateImport(ctx context.Context, options CollectionOptions, text string, bindings map[string]string) error
	// ColStateUpdate merges external bindings into stored state.
	ColStateUpdate(ctx context.Context, options CollectionOptions, value *ankiascodepb.State) error
	// ColStateReplace installs the complete supplied state object.
	ColStateReplace(ctx context.Context, options CollectionOptions, value *ankiascodepb.State) error
}

// CollectionOptions selects collection access and the namespace for one application operation.
// Col workflows require Path or Connection. An empty Address uses the configured
// factory default for a path, and an empty StateKey selects the default namespace.
// Raw API calls also accept Address alone to use the server's current collection.
type CollectionOptions struct {
	// Path identifies a managed database or archive; App owns preparation and cleanup.
	Path string
	// Address selects a client using subprocess:// or a plaintext TCP/Unix endpoint.
	// Managed paths require shared source/temporary paths and exclusive server access.
	// External server processes remain caller-owned; subprocess clients own their worker.
	// Address cannot accompany Connection.
	Address string
	// Connection borrows an already-open collection with caller-owned lifetime and publication.
	// It cannot accompany Path or Address; raw RPCs retain their requested native effects.
	Connection grpc.ClientConnInterface
	// StateKey selects the namespace for Col workflows; empty selects the default namespace.
	// Raw API request payloads control native semantics and do not use this field.
	StateKey string
}

// APIRequest selects an upstream method, protobuf payload and collection access mode.
// Collection requires exactly one path, client address or open connection.
type APIRequest struct {
	// Collection selects managed path access, a client address, or a borrowed open connection.
	Collection CollectionOptions
	// WriteBack enables publication of managed source archives, including after request errors.
	// It defaults to false and has no effect on borrowed connections or remote servers.
	WriteBack bool
	// Method is the fully qualified gRPC route, including its leading slash.
	Method string
	// Input contains one protobuf JSON or TOML request.
	Input []byte
	// InputFormat selects json or toml; empty selects json.
	InputFormat string
	// OutputFormat selects json or toml; empty selects json.
	OutputFormat string
}

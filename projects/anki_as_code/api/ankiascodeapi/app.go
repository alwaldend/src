package ankiascodeapi

import (
	"context"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
)

// App defines configuration, collection, state and server workflows for callers
// and command implementations without depending on a concrete application or
// native worker. Collection operations select their source per call; borrowed
// connections retain caller-owned lifetime and package publication.
type App interface {
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
	ColExport(ctx context.Context, collection CollectionSource, output, stateOutput, stateKey string) error
	// ColPlan describes proposed collection changes.
	ColPlan(ctx context.Context, collection CollectionSource, text, stateKey string) (*ankiascodepb.CollectionChangePlan, error)
	// ColApply reconciles the supplied collection in place.
	ColApply(ctx context.Context, collection CollectionSource, text, stateKey string) error
	// ColBuild stages a base archive path and writes an independent reconciled archive.
	ColBuild(ctx context.Context, base, text, output, stateKey string) error
	// ColStateDump reads stored collection bindings.
	ColStateDump(ctx context.Context, collection CollectionSource, stateKey string) (*ankiascodepb.State, error)
	// ColStateImport binds existing objects to configuration paths.
	ColStateImport(ctx context.Context, collection CollectionSource, text, stateKey string, bindings map[string]string) error
	// ColStateUpdate merges external bindings into stored state.
	ColStateUpdate(ctx context.Context, collection CollectionSource, stateKey string, value *ankiascodepb.State) error
	// ColStateReplace installs the complete supplied state object.
	ColStateReplace(ctx context.Context, collection CollectionSource, stateKey string, value *ankiascodepb.State) error
}

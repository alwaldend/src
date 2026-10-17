package cli

import (
	"context"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
)

// application describes the workflows consumed by the complete command tree.
// Tests substitute these operations without constructing the concrete App.
type application interface {
	// ConfigDump returns encoded collection configurations.
	ConfigDump(ctx context.Context, text, format string) ([]byte, error)
	// ConfigExport writes a parsed configuration set as a TOML tree.
	ConfigExport(ctx context.Context, output string, configs *collectionpb.CollectionConfigSet) error
	// ColExport writes editable configs and optional snapshot bindings.
	ColExport(ctx context.Context, collection, output, stateOutput, stateKey string) error
	// ColPlan describes proposed collection changes.
	ColPlan(ctx context.Context, collection, text, stateKey string) (*collectionpb.CollectionChangePlan, error)
	// ColApply reconciles the supplied collection in place.
	ColApply(ctx context.Context, collection, text, stateKey string) error
	// ColBuild writes an independent reconciled archive.
	ColBuild(ctx context.Context, base, text, output, stateKey string) error
	// ColStateDump reads stored collection bindings.
	ColStateDump(ctx context.Context, collection, stateKey string) (*collectionpb.State, error)
	// ColStateImport binds existing objects to configuration paths.
	ColStateImport(ctx context.Context, collection, text, stateKey string, bindings map[string]string) error
	// ColStateUpdate merges external bindings into stored state.
	ColStateUpdate(ctx context.Context, collection, stateKey string, value *collectionpb.State) error
	// ColStateReplace installs the complete supplied state object.
	ColStateReplace(ctx context.Context, collection, stateKey string, value *collectionpb.State) error
}

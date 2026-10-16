package export

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
)

// Exporter projects a prepared native collection into editable configs.
// The caller owns the connection, collection and backend lifetime.
type Exporter struct {
	// resources reads native values through the shared resource store.
	resources *resource.Store
	// states reads existing bindings for filenames and preset aliases without writes.
	states *state.Store
}

// New borrows resource and state stores for the same already-open collection.
// Both remain caller-owned. No native or filesystem work occurs.
func New(resources *resource.Store, states *state.Store) *Exporter {
	return &Exporter{resources: resources, states: states}
}

// Export writes editable TOML without changing native data or exporting state.
// The caller supplies an absolute output path that does not contain the source
// collection. Other files are untouched; filesystem errors may leave partial output.
// Use a fresh directory for a clean snapshot. State operations remain separate.
func (exporter *Exporter) Export(ctx context.Context, output string) error {
	current, err := exporter.resources.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("read export snapshot: %w", err)
	}
	stored, err := exporter.states.Dump(ctx)
	if err != nil {
		return fmt.Errorf("read export bindings: %w", err)
	}
	configs, err := exportConfigs(ctx, current, stored, exporter.resources.Rendering)
	if err != nil {
		return fmt.Errorf("project export configs: %w", err)
	}
	if err := config.Export(configs, output); err != nil {
		return fmt.Errorf("write exported configs: %w", err)
	}
	return nil
}

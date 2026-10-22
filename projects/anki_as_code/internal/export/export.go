package export

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
)

// Exporter projects a prepared native collection into editable configs.
// The caller owns the connection, collection and backend lifetime.
type Exporter struct {
	// filesystem supplies all config and state-output access.
	filesystem filesystem.FS
	// resources reads native values through the shared resource store.
	resources *resource.Store
	// states reads existing bindings for filenames and preset aliases without writes.
	states *state.Store
}

// New borrows resource and state stores for the same already-open collection.
// Both remain caller-owned. storage supplies output access, defaulting to OS storage.
// No native or filesystem work occurs.
func New(resources *resource.Store, states *state.Store, storage filesystem.FS) *Exporter {
	return &Exporter{resources: resources, states: states, filesystem: storage}
}

// Export writes configs and optionally publishes their version-1 TOML bindings.
// Paths use the injected filesystem's namespace and must not alias the source collection. State publication
// follows successful config output; failures may leave partial configs, but no partial state file.
// Existing stored bindings supply names only; the generated state describes exported objects.
// An empty stateOutput selects config-only export; state installation remains separate.
func (exporter *Exporter) Export(ctx context.Context, output, stateOutput string) error {
	output, err := exporter.filesystem.ResolvePath(output, true)
	if err != nil {
		return fmt.Errorf("resolve configuration output: %w", err)
	}
	if stateOutput != "" {
		resolved, err := exporter.filesystem.ResolvePath(stateOutput, true)
		if err != nil {
			return fmt.Errorf("resolve state output: %w", err)
		}
		stateOutput = resolved
	}
	current, err := exporter.resources.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("read export snapshot: %w", err)
	}
	stored, err := exporter.states.Dump(ctx)
	if err != nil {
		return fmt.Errorf("read export bindings: %w", err)
	}
	configs, identities, err := exportConfigs(ctx, current, stored, exporter.resources.Rendering)
	if err != nil {
		return fmt.Errorf("project export configs: %w", err)
	}
	var encoded []byte
	if stateOutput != "" {
		if err := validateStateOutput(exporter.filesystem, output, stateOutput, identities); err != nil {
			return fmt.Errorf("validate state output: %w", err)
		}
		bindings, err := state.Export(configs, current, identities)
		if err != nil {
			return fmt.Errorf("generate exported state: %w", err)
		}
		encoded, err = codec.Encode(bindings, codec.TOML)
		if err != nil {
			return fmt.Errorf("encode exported state: %w", err)
		}
	}
	if err := config.New(output, exporter.filesystem).Export(configs); err != nil {
		return fmt.Errorf("write exported configs: %w", err)
	}
	if stateOutput != "" {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("publish exported state: %w", err)
		}
		if err := writeStateOutput(exporter.filesystem, stateOutput, encoded); err != nil {
			return fmt.Errorf("write exported state: %w", err)
		}
	}
	return nil
}

// validateStateOutput rejects declaration filenames and ancestors needed by config output.
func validateStateOutput(storage filesystem.FS, output, destination string, identities map[string]int64) error {
	ancestor, err := filepath.Rel(destination, output)
	if err != nil {
		return fmt.Errorf("compare state and config directories: %w", err)
	}
	if filepath.IsLocal(ancestor) {
		return fmt.Errorf("state output must not replace the config directory or its ancestors")
	}
	relative, err := filepath.Rel(output, destination)
	if err != nil {
		return fmt.Errorf("compare export paths: %w", err)
	}
	paths := []string{"collection.anki.toml"}
	for path := range identities {
		paths = append(paths, path)
	}
	stateInfo, err := storage.Stat(destination)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect state destination: %w", err)
	}
	for _, path := range paths {
		configPath, err := storage.ResolvePath(filepath.Join(output, filepath.FromSlash(path)), true)
		if err != nil {
			return fmt.Errorf("resolve config destination %q: %w", path, err)
		}
		contained, err := filepath.Rel(destination, configPath)
		if err != nil {
			return fmt.Errorf("compare state and config destination %q: %w", path, err)
		}
		if filepath.IsLocal(contained) {
			return fmt.Errorf("state output conflicts with resolved config destination %q", path)
		}
		nested, err := filepath.Rel(configPath, destination)
		if err != nil {
			return fmt.Errorf("compare config and state destination %q: %w", path, err)
		}
		if filepath.IsLocal(nested) {
			return fmt.Errorf("state output must not be nested beneath config %q", path)
		}
		info, err := storage.Stat(configPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect config destination %q: %w", path, err)
		}
		if stateInfo != nil && info != nil && os.SameFile(stateInfo, info) {
			return fmt.Errorf("state output aliases config %q", path)
		}
	}
	if !filepath.IsLocal(relative) {
		return nil
	}
	address := filepath.ToSlash(relative)
	// Declaration suffixes would make the config loader interpret state as a config.
	if filepath.Base(destination) == "collection.anki.toml" || strings.HasSuffix(address, ".note.anki.toml") || strings.HasSuffix(address, ".notetype.anki.toml") || filepath.Base(destination) == "deck.anki.toml" {
		return fmt.Errorf("state output %q uses a configuration filename", destination)
	}
	for _, path := range paths {
		if strings.HasPrefix(path, address+"/") {
			return fmt.Errorf("state output %q conflicts with a config directory", destination)
		}
	}
	return nil
}

// writeStateOutput stages complete state beside its destination before atomic replacement.
func writeStateOutput(storage filesystem.FS, destination string, data []byte) (resultErr error) {
	if err := storage.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create state output directory: %w", err)
	}
	file, staged, err := storage.TempFile(filepath.Dir(destination), ".export-state-*")
	if err != nil {
		return fmt.Errorf("stage state output: %w", err)
	}
	defer func() {
		if err := storage.Remove(staged); err != nil && !errors.Is(err, os.ErrNotExist) {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove staged state: %w", err))
		}
	}()
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write staged state: %w", err)
	}
	if err := storage.Rename(staged, destination); err != nil {
		return fmt.Errorf("publish state output %q: %w", destination, err)
	}
	return nil
}

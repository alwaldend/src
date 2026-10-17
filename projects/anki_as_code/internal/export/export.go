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

// Export writes configs and optionally publishes their version-1 TOML bindings.
// Paths must be absolute and must not alias the source collection. State publication
// follows successful config output; failures may leave partial configs, but no partial state file.
// Existing stored bindings supply names only; the generated state describes exported objects.
// An empty stateOutput selects config-only export; state installation remains separate.
func (exporter *Exporter) Export(ctx context.Context, output, stateOutput string) error {
	if stateOutput != "" {
		resolved, err := ResolveOutputPath(stateOutput)
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
		if err := validateStateOutput(output, stateOutput, identities); err != nil {
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
	if err := config.New(output).Export(configs); err != nil {
		return fmt.Errorf("write exported configs: %w", err)
	}
	if stateOutput != "" {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("publish exported state: %w", err)
		}
		if err := writeStateOutput(stateOutput, encoded); err != nil {
			return fmt.Errorf("write exported state: %w", err)
		}
	}
	return nil
}

// validateStateOutput rejects declaration filenames and ancestors needed by config output.
func validateStateOutput(output, destination string, identities map[string]int64) error {
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
	stateInfo, err := os.Stat(destination)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inspect state destination: %w", err)
	}
	for _, path := range paths {
		configPath, err := ResolveOutputPath(filepath.Join(output, filepath.FromSlash(path)))
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
		info, err := os.Stat(configPath)
		if err != nil && !os.IsNotExist(err) {
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

// ResolveOutputPath resolves existing symlinks, including dangling links, for export and App.
// The input must be absolute; nonexistent suffixes remain available for output creation.
func ResolveOutputPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("export path %q must be absolute", path)
	}
	pending := filepath.Clean(path)
	for links := 0; links < 255; links++ {
		ancestor := pending
		for {
			info, err := os.Lstat(ancestor)
			if err == nil {
				suffix, err := filepath.Rel(ancestor, pending)
				if err != nil {
					return "", fmt.Errorf("resolve export suffix: %w", err)
				}
				if info.Mode()&os.ModeSymlink != 0 {
					target, err := os.Readlink(ancestor)
					if err != nil {
						return "", fmt.Errorf("read export symlink %q: %w", ancestor, err)
					}
					if !filepath.IsAbs(target) {
						target = filepath.Join(filepath.Dir(ancestor), target)
					}
					pending = filepath.Join(target, suffix)
					break
				}
				resolved, err := filepath.EvalSymlinks(ancestor)
				if err != nil {
					return "", fmt.Errorf("resolve export ancestor %q: %w", ancestor, err)
				}
				return filepath.Join(resolved, suffix), nil
			}
			parent := filepath.Dir(ancestor)
			if !os.IsNotExist(err) || parent == ancestor {
				return "", fmt.Errorf("inspect export path %q: %w", ancestor, err)
			}
			ancestor = parent
		}
	}
	return "", fmt.Errorf("too many symlinks resolving export path %q", path)
}

// writeStateOutput stages complete state beside its destination before atomic replacement.
func writeStateOutput(destination string, data []byte) (resultErr error) {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create state output directory: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(destination), ".export-state-*")
	if err != nil {
		return fmt.Errorf("stage state output: %w", err)
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil && !os.IsNotExist(err) {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove staged state: %w", err))
		}
	}()
	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		return fmt.Errorf("write staged state: %w", err)
	}
	if err := os.Rename(file.Name(), destination); err != nil {
		return fmt.Errorf("publish state output %q: %w", destination, err)
	}
	return nil
}

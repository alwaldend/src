package config

import (
	"fmt"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config/layout"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
)

// Version identifies the supported configuration schema for collection exporters and loaders.
const Version = layout.Version

// Store reads and writes schema-backed configuration files beneath a selected root.
// It retains no parsed messages or open files. Callers own collection semantics
// and coordinate filesystem access.
type Store struct {
	// filesystem owns every configuration read, write and path inspection.
	filesystem filesystem.FS
	// root selects the configuration directory, or a root file for read operations.
	root string
}

// New selects a configuration root without accessing the filesystem.
// Load and Dump accept a directory or root config file; Export requires a directory.
// storage supplies all file access; its zero value selects the OS filesystem.
func New(root string, storage filesystem.FS) *Store {
	return &Store{root: root, filesystem: storage}
}

// Load parses the selected root file or directory into grouped protobuf configs.
// Syntax, schema types, versions, and filesystem layout are checked; callers
// validate semantics. Each call reads fresh values without changing source files.
func (store *Store) Load() (*ankiascodepb.CollectionConfigSet, error) {
	configs, err := layout.Load(store.filesystem, store.root)
	if err != nil {
		return nil, fmt.Errorf("load configuration layout: %w", err)
	}
	return configs, nil
}

// Dump loads and encodes the selected configuration set without opening Anki.
// JSON is the default; TOML is also supported.
func (store *Store) Dump(format string) ([]byte, error) {
	if format == "" {
		format = string(codec.JSON)
	}
	if codec.Format(format) != codec.JSON && codec.Format(format) != codec.TOML {
		return nil, fmt.Errorf("unsupported configuration format %q; use json or toml", format)
	}
	configs, err := store.Load()
	if err != nil {
		return nil, fmt.Errorf("load configuration set: %w", err)
	}
	data, err := codec.Encode(configs, codec.Format(format))
	if err != nil {
		return nil, fmt.Errorf("encode configuration set: %w", err)
	}
	return data, nil
}

// Export writes supplied configurations as TOML files beneath the selected directory.
// Callers own validation and unique, representable paths. Declared files are
// overwritten; other files are untouched. Failures may leave partial output.
func (store *Store) Export(configs *ankiascodepb.CollectionConfigSet) error {
	if err := layout.Export(store.filesystem, configs, store.root); err != nil {
		return fmt.Errorf("export configuration layout: %w", err)
	}
	return nil
}

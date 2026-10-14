package config

import (
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config/layout"
)

// Version identifies the supported configuration schema for collection exporters and loaders.
const Version = layout.Version

// Load parses a root config file or its directory into the grouped protobuf configuration set.
// Syntax, schema types, versions, and filesystem layout are checked; callers validate semantics.
// Source files remain unchanged, and no Anki collection is opened.
func Load(input string) (*collectionpb.CollectionConfigSet, error) {
	configs, err := layout.Load(input)
	if err != nil {
		return nil, fmt.Errorf("load configuration layout: %w", err)
	}
	return configs, nil
}

// Export writes supplied configurations as TOML files in the fixed hierarchy.
// Callers own validation and unique, representable paths. Declared files are
// overwritten; other files are untouched. Failures may leave partial output.
func Export(configs *collectionpb.CollectionConfigSet, output string) error {
	if err := layout.Export(configs, output); err != nil {
		return fmt.Errorf("export configuration layout: %w", err)
	}
	return nil
}

// IsConfigFile identifies declaration filenames for callers managing export directories.
// It performs no filesystem inspection or configuration validation.
func IsConfigFile(name string) bool {
	return layout.IsConfigFile(name)
}

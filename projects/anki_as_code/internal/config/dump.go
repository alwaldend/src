package config

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config/codec"
)

// Dump loads and encodes the configuration set for the CLI adapter.
// JSON is the default; TOML is also supported. No Anki operations or ID allocation occur.
func Dump(input, format string) ([]byte, error) {
	if format == "" {
		format = "json"
	}
	if format != "json" && format != "toml" {
		return nil, fmt.Errorf("unsupported configuration format %q; use json or toml", format)
	}
	configs, err := Load(input)
	if err != nil {
		return nil, fmt.Errorf("load configuration set: %w", err)
	}
	if format == "toml" {
		data, err := codec.EncodeTOML(configs)
		if err != nil {
			return nil, fmt.Errorf("encode configuration set as TOML: %w", err)
		}
		return data, nil
	}
	data, err := codec.EncodeJSON(configs)
	if err != nil {
		return nil, fmt.Errorf("encode configuration set as JSON: %w", err)
	}
	return data, nil
}

package config

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config/codec"
)

// Dump loads and encodes the configuration set for the CLI adapter.
// JSON is the default; TOML is also supported. No Anki operations or ID allocation occur.
func Dump(input, format string) ([]byte, error) {
	if format == "" {
		format = string(codec.JSON)
	}
	if codec.Format(format) != codec.JSON && codec.Format(format) != codec.TOML {
		return nil, fmt.Errorf("unsupported configuration format %q; use json or toml", format)
	}
	configs, err := Load(input)
	if err != nil {
		return nil, fmt.Errorf("load configuration set: %w", err)
	}
	data, err := codec.Encode(configs, codec.Format(format))
	if err != nil {
		return nil, fmt.Errorf("encode configuration set: %w", err)
	}
	return data, nil
}

package layout

import (
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/filesystem"
	"google.golang.org/protobuf/proto"
)

// readTOML reads a TOML declaration and validates it through its protobuf schema.
func readTOML(storage filesystem.FS, path string, message proto.Message) error {
	b, err := storage.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read TOML input: %w", err)
	}
	if err := codec.Decode(b, codec.TOML, message); err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}
	return nil
}

// writeTOML serializes a protobuf message to a TOML declaration file.
func writeTOML(storage filesystem.FS, path string, message proto.Message) error {
	b, err := codec.Encode(message, codec.TOML)
	if err != nil {
		return fmt.Errorf("serialize config %q: %w", path, err)
	}
	if err := storage.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	return nil
}

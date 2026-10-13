package layout

import (
	"errors"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config/codec"
	"google.golang.org/protobuf/proto"
)

// readFile reads declaration bytes and reports read and close failures.
func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open text %q: %w", path, err)
	}
	b, readErr := io.ReadAll(f)
	closeErr := f.Close()
	if err := errors.Join(readErr, closeErr); err != nil {
		return nil, fmt.Errorf("read text %q: %w", path, err)
	}
	return b, nil
}

// readTOML reads a TOML declaration and validates it through its protobuf schema.
func readTOML(path string, message proto.Message) error {
	b, err := readFile(path)
	if err != nil {
		return fmt.Errorf("read TOML input: %w", err)
	}
	if err := codec.DecodeTOML(b, message); err != nil {
		return fmt.Errorf("parse %q: %w", path, err)
	}
	return nil
}

// writeTOML serializes a protobuf message to a TOML declaration file.
func writeTOML(path string, message proto.Message) error {
	b, err := codec.EncodeTOML(message)
	if err != nil {
		return fmt.Errorf("serialize config %q: %w", path, err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write config %q: %w", path, err)
	}
	return nil
}

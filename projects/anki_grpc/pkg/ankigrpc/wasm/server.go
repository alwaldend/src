package wasm

import (
	_ "embed"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/wasmhost"
)

// guest embeds the complete offline Anki Rust CLI owned by this adapter.
//
//go:embed guest.wasm
var guest []byte

// Run executes the embedded Rust CLI using a complete argument vector.
// Callers retain process isolation by invoking serving only in a worker process.
func Run(arguments []string) error {
	if err := wasmhost.Run(arguments, guest); err != nil {
		return fmt.Errorf("run embedded Anki WASM server: %w", err)
	}
	return nil
}

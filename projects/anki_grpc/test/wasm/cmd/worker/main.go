package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/internal/wasmhost"
)

// guest embeds the actual cross-compiled Rust module in the cgo-disabled worker.
//
//go:embed fixture.wasm
var guest []byte

// main runs the server using the same argument-vector contract as the native adapter.
func main() {
	if err := wasmhost.Run(os.Args, guest); err != nil {
		fmt.Fprintln(os.Stderr, err)
		var status interface{ ExitCode() uint32 }
		if errors.As(err, &status) {
			os.Exit(int(status.ExitCode()))
		}
		os.Exit(1)
	}
}

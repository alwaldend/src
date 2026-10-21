package main

import (
	"errors"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/wasm"
)

// main adapts the process argument vector to the WASM server host.
func main() {
	if err := wasm.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		var status interface{ ExitCode() uint32 }
		if errors.As(err, &status) {
			os.Exit(int(status.ExitCode()))
		}
		os.Exit(1)
	}
}

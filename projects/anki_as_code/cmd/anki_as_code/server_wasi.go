//go:build anki_wasm

package main

import (
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/wasm"
)

// wasmServerRunner hosts the full Rust CLI inside the supervised WASM worker.
var wasmServerRunner = wasm.Run

// wasmClientFactory constructs independent in-process WASM instances without startup I/O.
var wasmClientFactory ankigrpc.ClientFactory = ankigrpc.ClientFactoryFunc[*wasm.Client](func(string) (*wasm.Client, error) {
	return wasm.New(), nil
})

//go:build !anki_wasm

package main

import "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"

// wasmServerRunner records that this executable has no embedded WASM engine.
var wasmServerRunner func([]string) error

// wasmClientFactory records that this executable has no in-process WASM engine.
var wasmClientFactory ankigrpc.ClientFactory

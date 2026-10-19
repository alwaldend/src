//go:build !anki_wasm

package main

import "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiservercgo"

// runServer selects the server adapter invoked by the worker command.
var runServer = ankiservercgo.Run

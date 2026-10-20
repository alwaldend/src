package main

import "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"

// runServer hosts the native CLI owned by this executable.
var runServer = native.RunServer

// subprocessSchemes names the worker strategies packaged by this executable.
var subprocessSchemes = []string{"subprocess", "subprocess+cgo"}

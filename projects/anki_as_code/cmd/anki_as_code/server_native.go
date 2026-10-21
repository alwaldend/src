//go:build anki_cgo

package main

import "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"

// nativeServerRunner hosts the native CLI and imports native in-process clients.
var nativeServerRunner = native.RunServer

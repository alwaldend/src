package main

import (
	"os"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/test/wasmcache"
)

// TestMain configures shared test resources before any executable workflows run.
func TestMain(tests *testing.M) { os.Exit(wasmcache.Run(tests)) }

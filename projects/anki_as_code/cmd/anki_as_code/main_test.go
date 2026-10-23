package main

import (
	"fmt"
	"os"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/test/wasmcache"
)

// TestMain configures shared test resources before any executable workflows run.
func TestMain(tests *testing.M) {
	if len(os.Args) == 4 && os.Args[1] == "mcp-worker-fixture" {
		if err := runMCPWorkerFixture(os.Args[2], os.Args[3]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	os.Exit(wasmcache.Run(tests))
}

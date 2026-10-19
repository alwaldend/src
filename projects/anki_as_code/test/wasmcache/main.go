package wasmcache

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// Run executes an E2E suite with a private cache shared by that suite's workers.
func Run(tests *testing.M) int {
	if directory := os.Getenv("TEST_TMPDIR"); directory != "" {
		if err := os.Setenv("ANKI_WASM_CACHE_DIR", filepath.Join(directory, "anki-wasm-compiled")); err != nil {
			fmt.Fprintf(os.Stderr, "configure test compilation cache: %v\n", err)
			return 1
		}
	}
	return tests.Run()
}

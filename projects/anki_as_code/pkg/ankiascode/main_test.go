package ankiascode

import (
	"os"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/test/wasmcache"
)

// TestMain configures shared test resources before application E2E workflows run.
func TestMain(tests *testing.M) { os.Exit(wasmcache.Run(tests)) }

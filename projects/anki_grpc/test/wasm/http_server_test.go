package wasm

import (
	"bytes"
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// httpServerProbe selects the guest that calls Anki's HTTP sync-server API.
var httpServerProbe = flag.String("http-server-probe", "", "HTTP sync-server WASI probe runfile")

// TestHTTPServerUnsupported verifies that server startup fails before filesystem or network work.
func TestHTTPServerUnsupported(t *testing.T) {
	path, err := runfiles.Rlocation(*httpServerProbe)
	if err != nil {
		t.Fatalf("resolve HTTP sync-server probe: %v", err)
	}
	code, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read HTTP sync-server probe: %v", err)
	}
	runtime := wazero.NewRuntime(t.Context())
	t.Cleanup(func() {
		if err := runtime.Close(context.Background()); err != nil {
			t.Errorf("close HTTP sync-server probe runtime: %v", err)
		}
	})
	if _, err := wasi_snapshot_preview1.Instantiate(t.Context(), runtime); err != nil {
		t.Fatalf("initialize HTTP sync-server WASI probe: %v", err)
	}
	var output bytes.Buffer
	if _, err := runtime.InstantiateWithConfig(t.Context(), code, wazero.NewModuleConfig().WithStdout(&output).WithStderr(&output)); err != nil {
		t.Fatalf("execute HTTP sync-server WASI probe: %v; %s", err, output.String())
	}
	if !strings.Contains(output.String(), "unsupported in offline WASM backend: HTTP sync server") {
		t.Fatalf("HTTP sync server did not reject startup explicitly: %s", output.String())
	}
}

package wasm

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/generic"
)

// TestGuestTimezones verifies historical offsets, local-time gaps and overlaps in Rust.
func TestGuestTimezones(t *testing.T) {
	path := workerExecutable(t)
	for _, fixture := range []struct {
		// zone selects the host timezone inherited by a fresh worker.
		zone string
		// method selects the Rust Chrono assertions for that timezone.
		method string
	}{
		{"America/New_York", "TimezoneNewYork"},
		{"Australia/Lord_Howe", "TimezoneLordHowe"},
		{"Pacific/Apia", "TimezoneApia"},
	} {
		t.Run(fixture.zone, func(t *testing.T) {
			t.Setenv("TZ", fixture.zone)
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
			defer cancel()
			worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: path, Arguments: []string{"serve"}})
			if err := worker.Start(ctx); err != nil {
				t.Fatalf("start timezone worker: %v", err)
			}
			t.Cleanup(func() {
				stop, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := worker.Stop(stop); err != nil {
					t.Errorf("reap timezone worker: %v", err)
				}
			})
			connection, err := worker.Connection()
			if err != nil {
				t.Fatalf("get timezone connection: %v", err)
			}
			if err := connection.Invoke(ctx, "/fixture/"+fixture.method, &generic.Empty{}, &generic.Empty{}); err != nil {
				t.Fatalf("Rust timezone contract in %s: %v", fixture.zone, err)
			}
		})
	}
	if t.Failed() {
		return
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "timezone-verification.txt"), []byte("New York seasonal offsets/gap/overlap; Lord Howe half-hour transitions; Apia skipped date passed through Rust Chrono and Go host\n"), 0o600); err != nil {
			t.Fatalf("retain timezone verification: %v", err)
		}
	}
}

package injector

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// TestCleanerConcurrentRegistrationAndStop verifies path ownership during shutdown.
func TestCleanerConcurrentRegistrationAndStop(t *testing.T) {
	cleaner := NewCleaner()
	dir := t.TempDir()
	var wg sync.WaitGroup
	for i := range 64 {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			path := filepath.Join(dir, fmt.Sprint(i))
			if err := os.WriteFile(path, []byte("synthetic fixture"), 0o600); err != nil {
				t.Error(err)
				return
			}
			if err := cleaner.Add(path); err != nil {
				if err := os.RemoveAll(path); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	if err := cleaner.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	wg.Wait()
	if err := cleaner.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 0 {
		t.Fatalf("%d files remained after cleanup", len(files))
	}
	if err := cleaner.Add("unused"); err == nil {
		t.Fatal("registration after stop succeeded")
	}
}

// TestFileFetcherCleanup verifies private-file permissions and rollback.
func TestFileFetcherCleanup(t *testing.T) {
	for _, stopped := range []bool{false, true} {
		t.Run(fmt.Sprint(stopped), func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("TMPDIR", dir)
			cleaner := NewCleaner()
			if stopped {
				if err := cleaner.Stop(context.Background()); err != nil {
					t.Fatal(err)
				}
			}
			fetcher := newFileFetcher(&templater{}, cleaner)
			res, err := fetcher.Get(context.Background(), &Resource{Name: "fixture", Type: "file", Args: map[string]any{"value": "synthetic fixture"}}, nil)
			if stopped {
				if err == nil {
					t.Fatal("closed cleaner accepted resource")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				info, err := os.Stat(res.Files[0])
				if err != nil {
					t.Fatal(err)
				}
				if info.Mode().Perm() != 0o600 {
					t.Fatalf("mode = %v", info.Mode())
				}
				if err := cleaner.Stop(context.Background()); err != nil {
					t.Fatal(err)
				}
			}
			files, err := os.ReadDir(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != 0 {
				t.Fatalf("%d secret files remained", len(files))
			}
		})
	}
}

// TestTemplateErrorsRedactSourceAndValues verifies value-free diagnostics.
func TestTemplateErrorsRedactSourceAndValues(t *testing.T) {
	for _, tpl := range []string{`{{ synthetic_fixture_secret }}`, `{{index "synthetic_fixture_secret" "bad"}}`, `{{.Extra.synthetic_fixture_secret}}`} {
		_, err := (&templater{}).Template(context.Background(), tpl, nil, nil).Get()
		if err == nil {
			t.Fatal("expected template error")
		}
		if strings.Contains(err.Error(), "synthetic_fixture_secret") {
			t.Fatal("template error disclosed input")
		}
	}
}

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/tools/vault/injector/injector_proto"
)

func TestCleanerConcurrentRegistrationAndStop(t *testing.T) {
	cleaner := NewCleaner(nil)
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

func TestFileFetcherCleanup(t *testing.T) {
	for _, stopped := range []bool{false, true} {
		t.Run(fmt.Sprint(stopped), func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("TMPDIR", dir)
			cleaner := NewCleaner(nil)
			if stopped {
				if err := cleaner.Stop(context.Background()); err != nil {
					t.Fatal(err)
				}
			}
			fetcher := NewFileFetcher(nil, &Templater{}, cleaner)
			res, err := fetcher.Get(context.Background(), &injector_proto.Resource{Name: "fixture", Res: &injector_proto.Resource_File{File: &injector_proto.File{Value: "synthetic fixture"}}}, nil)
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

func TestTemplateErrorsRedactSourceAndValues(t *testing.T) {
	for _, tpl := range []string{`{{ synthetic_fixture_secret }}`, `{{index "synthetic_fixture_secret" "bad"}}`, `{{.Extra.synthetic_fixture_secret}}`} {
		_, err := (&Templater{}).Template(context.Background(), tpl, nil, nil).Get()
		if err == nil {
			t.Fatal("expected template error")
		}
		if strings.Contains(err.Error(), "synthetic_fixture_secret") {
			t.Fatal("template error disclosed input")
		}
	}
}

func fixtureVaultConfig(address string) *al_proto.Config {
	return &al_proto.Config{
		VaultConn: []*al_proto.VaultConn{{Name: "default", Config: &al_proto.VaultConfig{Address: address}, Tls: &al_proto.VaultTLS{}}},
		VaultAuth: []*al_proto.VaultAuth{{Name: "default", NoAuth: true}},
	}
}

func TestNoAuthOverridesInheritedVaultToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "synthetic-inherited-token")
	config := fixtureVaultConfig("http://127.0.0.1:1")
	fetcher := NewVaultEnvFetcher(&Templater{}, config, al.NewVault(config))
	result, err := fetcher.Get(context.Background(), &injector_proto.Resource{Res: &injector_proto.Resource_VaultEnv{VaultEnv: &injector_proto.VaultEnv{}}}, nil)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result.Env["VAULT_TOKEN"]
	if !ok || value != "" {
		t.Fatal("no_auth must explicitly override the inherited token")
	}
}

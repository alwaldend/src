package terraform_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

var (
	withProvider    = flag.String("with-provider", "", "runfile path of the executable declaring the local provider")
	withoutProvider = flag.String("without-provider", "", "runfile path of the executable without providers")
)

const providerVersion = "2.9.0"

func TestDeclaredProviderInstallationStaysOffline(t *testing.T) {
	rf, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	var requests atomic.Int64
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests.Add(1)
		http.Error(w, "provider network access is forbidden by this fixture", http.StatusBadGateway)
	}))
	t.Cleanup(proxy.Close)

	scratch := t.TempDir()
	invocationDirectory := t.TempDir()
	cache := t.TempDir()
	poisonConfig := filepath.Join(scratch, "poisoned.terraformrc")
	if err := os.WriteFile(poisonConfig, []byte("provider_installation {\n  direct {}\n}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	baseEnvironment := append(rf.Env(),
		"TEST_WORKSPACE="+os.Getenv("TEST_WORKSPACE"),
		"TEST_TMPDIR="+scratch,
		// Terraform changes to the explicit fixture directory before starting
		// plugins. A relative temporary path keeps Unix socket names short and
		// inside the test-owned workspace, even with long Bazel output paths.
		"TMPDIR=.",
		"HTTP_PROXY="+proxy.URL,
		"HTTPS_PROXY="+proxy.URL,
		"http_proxy="+proxy.URL,
		"https_proxy="+proxy.URL,
		"NO_PROXY=",
		"no_proxy=",
		"TF_CLI_CONFIG_FILE="+poisonConfig,
		"TF_PLUGIN_CACHE_MAY_BREAK_DEPENDENCY_LOCK_FILE=true",
	)
	run := func(target, directory, pluginCache string, args ...string) (string, error) {
		t.Helper()
		binary, err := rf.Rlocation(target)
		if err != nil {
			t.Fatal(err)
		}
		binary, err = filepath.Abs(binary)
		if err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, binary, append([]string{"--chdir", directory}, args...)...)
		command.Dir = invocationDirectory
		command.Env = append(append([]string{}, baseEnvironment...), "TF_PLUGIN_CACHE_DIR="+pluginCache)
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("Terraform command exceeded its deadline: %v\n%s", ctx.Err(), output)
		}
		return string(output), err
	}

	declared := terraformWorkspace(t)
	if output, err := run(*withProvider, declared, cache, "init", "-backend=false", "-input=false", "-no-color"); err != nil {
		t.Fatalf("initialize declared provider: %v\n%s", err, output)
	}
	if output, err := run(*withProvider, declared, cache, "validate", "-no-color"); err != nil {
		t.Fatalf("validate using installed provider: %v\n%s", err, output)
	}
	lock, err := os.ReadFile(filepath.Join(declared, ".terraform.lock.hcl"))
	if err != nil {
		t.Fatalf("read generated lock in caller-owned workspace: %v", err)
	}
	normalizedLock := strings.Join(strings.Fields(string(lock)), " ")
	for _, want := range []string{`provider "registry.terraform.io/hashicorp/local"`, `version = "` + providerVersion + `"`} {
		if !strings.Contains(normalizedLock, want) {
			t.Errorf("generated lock does not record %s", want)
		}
	}
	installedProviders := filepath.Join(declared, ".terraform", "providers")
	installed := filepath.Join(installedProviders, "registry.terraform.io", "hashicorp", "local", providerVersion, runtime.GOOS+"_"+runtime.GOARCH)
	entries, err := os.ReadDir(installed)
	if err != nil {
		t.Fatalf("inspect installed pinned provider: %v", err)
	}
	foundBinary := false
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "terraform-provider-local_v"+providerVersion) {
			foundBinary = true
			info, err := entry.Info()
			if err != nil {
				t.Fatal(err)
			}
			if !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
				t.Errorf("installed provider %s is not an executable regular file", entry.Name())
			}
		}
	}
	if !foundBinary {
		t.Error("pinned provider executable was not installed")
	}
	if entries, err := os.ReadDir(cache); err != nil || len(entries) != 0 {
		t.Errorf("inherited plugin cache was used: entries=%v, error=%v", entries, err)
	}

	undeclared := terraformWorkspace(t)
	// A complete installation is available through the poisoned cache. The
	// target without a provider declaration must still reject the dependency.
	output, err := run(*withoutProvider, undeclared, installedProviders, "init", "-backend=false", "-input=false", "-no-color")
	if err == nil {
		t.Fatal("initialization succeeded without a declared provider")
	}
	normalizedOutput := strings.Join(strings.Fields(output), " ")
	for _, want := range []string{"hashicorp/local", "was not found in any of the search locations", "providers"} {
		if !strings.Contains(normalizedOutput, want) {
			t.Errorf("missing-provider diagnostic does not identify %q:\n%s", want, output)
		}
	}
	if count := requests.Load(); count != 0 {
		t.Errorf("Terraform made %d network requests during provider initialization or validation", count)
	}
	assertAbsent(t, filepath.Join(undeclared, ".terraform.lock.hcl"))
	for _, directory := range []string{invocationDirectory, scratch} {
		assertAbsent(t, filepath.Join(directory, ".terraform.lock.hcl"))
		assertAbsent(t, filepath.Join(directory, ".terraform"))
	}
	assertAbsent(t, filepath.Join(declared, "created.txt"))
	assertAbsent(t, filepath.Join(declared, "terraform.tfstate"))
}

func terraformWorkspace(t *testing.T) string {
	t.Helper()
	directory := t.TempDir()
	configuration := fmt.Sprintf(`terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "%s"
    }
  }
}

resource "local_file" "fixture" {
  filename = "${path.module}/created.txt"
  content  = "Offline Terraform validation fixture."
}
`, providerVersion)
	if err := os.WriteFile(filepath.Join(directory, "main.tf"), []byte(configuration), 0o600); err != nil {
		t.Fatal(err)
	}
	return directory
}

func assertAbsent(t *testing.T, name string) {
	t.Helper()
	if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("unexpected generated artifact %s: %v", name, err)
	}
}

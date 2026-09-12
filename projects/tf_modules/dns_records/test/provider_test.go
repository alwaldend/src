package dns_records_test

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var (
	terraform  = flag.String("terraform", "", "Pinned Terraform executable")
	cloudflare = flag.String("cloudflare", "", "Pinned Cloudflare provider")
	routeros   = flag.String("routeros", "", "Pinned RouterOS provider")
	proxmox    = flag.String("proxmox", "", "Pinned Proxmox provider")
)

func TestProviderDeclarations(t *testing.T) {
	runfiles := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"))
	resolve := func(path string) string { return filepath.Join(runfiles, path) }
	workspace := t.TempDir()
	source := resolve("projects/tf_modules/dns_records")
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tf") && !strings.HasSuffix(path, ".tftest.hcl") {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		destination := filepath.Join(workspace, relative)
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return err
		}
		return os.WriteFile(destination, content, 0o600)
	})
	if err != nil {
		t.Fatal(err)
	}
	cliConfig := installProviderMirror(t, workspace, []providerPackage{
		{"cloudflare", "cloudflare", "5.22.0", *cloudflare},
		{"terraform-routeros", "routeros", "1.99.1", *routeros},
	})
	for _, check := range []struct {
		name, directory string
		arguments       []string
	}{
		{"combined_mock_lifecycle", workspace, []string{"test", "-no-color"}},
		{"global_mock_lifecycle", filepath.Join(workspace, "global"), []string{"test", "-no-color"}},
	} {
		t.Run(check.name, func(t *testing.T) {
			for _, args := range [][]string{
				{"init", "-backend=false", "-input=false", "-no-color"},
				check.arguments,
			} {
				if output, err := runTerraform(check.directory, cliConfig, nil, args...); err != nil {
					t.Fatalf("terraform %v: %v\n%s", args, err, output)
				}
			}
		})
	}
}

type providerPackage struct {
	namespace, name, version, path string
}

func installProviderMirror(t *testing.T, workspace string, providers []providerPackage) string {
	t.Helper()
	runfiles := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"))
	mirror := filepath.Join(workspace, "mirror")
	for _, provider := range providers {
		directory := filepath.Join(mirror, "registry.terraform.io", provider.namespace, provider.name, provider.version, "linux_amd64")
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(runfiles, provider.path), filepath.Join(directory, "terraform-provider-"+provider.name+"_v"+provider.version)); err != nil {
			t.Fatal(err)
		}
	}
	cliConfig := filepath.Join(workspace, "terraform.rc")
	if err := os.WriteFile(cliConfig, []byte(fmt.Sprintf("provider_installation { filesystem_mirror { path = %q } }\n", mirror)), 0o600); err != nil {
		t.Fatal(err)
	}
	return cliConfig
}

func runTerraform(directory, cliConfig string, environment []string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	runfiles := filepath.Join(os.Getenv("TEST_SRCDIR"), os.Getenv("TEST_WORKSPACE"))
	command := exec.CommandContext(ctx, filepath.Join(runfiles, *terraform), args...)
	command.Dir = directory
	command.WaitDelay = 5 * time.Second
	// A relative temporary directory keeps provider Unix socket names within
	// the platform limit even when Bazel's isolated test root is long. The
	// explicit environment excludes inherited provider credentials.
	command.Env = []string{"PATH=" + os.Getenv("PATH"), "RUNFILES_DIR=" + os.Getenv("TEST_SRCDIR"), "TMPDIR=.", "TF_IN_AUTOMATION=1", "CHECKPOINT_DISABLE=1", "TF_CLI_CONFIG_FILE=" + cliConfig}
	command.Env = append(command.Env, environment...)
	return command.CombinedOutput()
}

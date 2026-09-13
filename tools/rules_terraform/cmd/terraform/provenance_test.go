package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func writeProviderPackage(t *testing.T, archive, directory string, contents map[string]string) {
	t.Helper()
	file, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	names := make([]string, 0, len(contents))
	for name := range contents {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		header := &zip.FileHeader{Name: name, Method: zip.Deflate}
		header.SetMode(0o755)
		entry, err := writer.CreateHeader(header)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(entry, contents[name]); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(directory, filepath.FromSlash(name)), contents[name])
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestInstalledProviderPackageMatchesDeclaredArchive(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "provider.zip")
	directory := filepath.Join(root, "installed")
	writeProviderPackage(t, archive, directory, map[string]string{
		"terraform-provider-fixture_v1.2.3": "provider executable",
		"LICENSE":                           "license",
		"docs/README":                       "additional metadata",
	})
	if err := verifyInstalledPackage(archive, directory); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(root, "cache-link")
	if err := os.Symlink(directory, linked); err != nil {
		t.Fatal(err)
	}
	if err := verifyInstalledPackage(archive, linked); err != nil {
		t.Fatalf("matching unpacked mirror or cache root rejected: %v", err)
	}
}

func TestInstalledProviderPackageRejectsChangedMissingAndExtraFiles(t *testing.T) {
	for _, change := range []string{"modified", "missing", "extra", "symlink"} {
		t.Run(change, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "provider.zip")
			directory := filepath.Join(root, "installed")
			const binary = "terraform-provider-fixture_v1.2.3"
			writeProviderPackage(t, archive, directory, map[string]string{binary: "declared executable", "LICENSE": "license"})
			switch change {
			case "modified":
				writeFile(t, filepath.Join(directory, binary), "foreign executable")
			case "missing":
				if err := os.Remove(filepath.Join(directory, binary)); err != nil {
					t.Fatal(err)
				}
			case "extra":
				writeFile(t, filepath.Join(directory, "terraform-provider-fixture_earlier"), "additional executable")
			case "symlink":
				if err := os.Remove(filepath.Join(directory, binary)); err != nil {
					t.Fatal(err)
				}
				target := filepath.Join(root, "outside-executable")
				writeFile(t, target, "declared executable")
				if err := os.Symlink(target, filepath.Join(directory, binary)); err != nil {
					t.Fatal(err)
				}
			}
			if err := verifyInstalledPackage(archive, directory); err == nil {
				t.Fatal("installed package accepted after " + change)
			}
		})
	}
}

func TestProviderArchiveRejectsSymlinkAndUnsafeEntries(t *testing.T) {
	for _, fixture := range []struct {
		name string
		mode os.FileMode
	}{
		{"terraform-provider-fixture", os.ModeSymlink | 0o777},
		{"../outside", 0o644},
		{"/absolute", 0o644},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			root := t.TempDir()
			archive := filepath.Join(root, "provider.zip")
			file, err := os.Create(archive)
			if err != nil {
				t.Fatal(err)
			}
			writer := zip.NewWriter(file)
			header := &zip.FileHeader{Name: fixture.name}
			header.SetMode(fixture.mode)
			entry, err := writer.CreateHeader(header)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.WriteString(entry, "fixture"); err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
			if err := verifyInstalledPackage(archive, root); err == nil {
				t.Fatal("unsafe archive entry accepted")
			}
		})
	}
}

func TestProviderSelectionUsesEffectiveTerraformDataDirectory(t *testing.T) {
	for _, kind := range []string{"default", "absolute", "relative"} {
		t.Run(kind, func(t *testing.T) {
			root := t.TempDir()
			workingDirectory := filepath.Join(root, "work")
			provider := fixtureProvider()
			var environ []string
			dataDirectory := filepath.Join(workingDirectory, ".terraform")
			switch kind {
			case "absolute":
				dataDirectory = filepath.Join(root, "custom-data")
				environ = []string{"TF_DATA_DIR=" + dataDirectory}
			case "relative":
				dataDirectory = filepath.Join(workingDirectory, "relative-data")
				environ = []string{"TF_DATA_DIR=relative-data"}
			}
			archive := filepath.Join(root, "provider.zip")
			installed := filepath.Join(dataDirectory, "providers", filepath.FromSlash(provider.Source), provider.Version, provider.Platform)
			writeProviderPackage(t, archive, installed, map[string]string{"terraform-provider-fixture": "declared executable"})
			call := invocation{runfiles: mapResolver{provider.Archive: archive}}
			config := configuration{Providers: []providerPackage{provider}}
			execute := func(args []string, output io.Writer) error {
				if !reflect.DeepEqual(args, []string{"version", "-json"}) {
					t.Fatalf("verification invoked provider/backend operation: %q", args)
				}
				fmt.Fprintf(output, `{"provider_selections":{%q:%q}}`, provider.Source, provider.Version)
				return nil
			}
			if err := call.verifyProviderSelections(config, workingDirectory, environ, execute); err != nil {
				t.Fatal(err)
			}
			writeFile(t, filepath.Join(installed, "terraform-provider-fixture"), "modified executable")
			if err := call.verifyProviderSelections(config, workingDirectory, environ, execute); err == nil {
				t.Fatal("modified package in effective TF_DATA_DIR was accepted")
			}
		})
	}
}

func TestProviderSelectionsRejectStaleAndUndeclaredProviders(t *testing.T) {
	provider := fixtureProvider()
	for _, selection := range []map[string]string{
		{provider.Source: "0.1.0"},
		{"registry.terraform.io/foreign/provider": "1.0.0"},
		{provider.Source: provider.Version, "registry.terraform.io/foreign/provider": "1.0.0"},
	} {
		root := t.TempDir()
		archive := filepath.Join(root, "provider.zip")
		installed := filepath.Join(root, ".terraform", "providers", filepath.FromSlash(provider.Source), provider.Version, provider.Platform)
		writeProviderPackage(t, archive, installed, map[string]string{"terraform-provider-fixture": "declared executable"})
		call := invocation{runfiles: mapResolver{provider.Archive: archive}}
		config := configuration{Providers: []providerPackage{provider}}
		execute := func(_ []string, output io.Writer) error {
			return json.NewEncoder(output).Encode(map[string]any{"provider_selections": selection})
		}
		if err := call.verifyProviderSelections(config, root, nil, execute); err == nil {
			t.Fatalf("foreign lock selections accepted: %v", selection)
		}
	}
}

func TestDirectAndInitializedCommandsVerifyBeforeLaunchingProviders(t *testing.T) {
	for _, direct := range []bool{true, false} {
		t.Run(fmt.Sprintf("direct_%t", direct), func(t *testing.T) {
			root := t.TempDir()
			provider := fixtureProvider()
			archive := filepath.Join(root, "provider.zip")
			writeProviderPackage(t, archive, filepath.Join(root, "fixture-package"), map[string]string{"terraform-provider-fixture": "declared executable"})
			args := []string{"--chdir", root, "plan"}
			if direct {
				args = append([]string{"--direct"}, args...)
			}
			config := configuration{Terraform: "terraform/terraform", Args: args, Providers: []providerPackage{provider}}
			var commands []string
			call := invocation{
				environ:  []string{"TMPDIR=" + root},
				runfiles: mapResolver{"terraform/terraform": filepath.Join(root, "terraform"), provider.Archive: archive},
				stdout:   io.Discard, stderr: io.Discard,
				execute: func(cmd *exec.Cmd) error {
					commands = append(commands, cmd.Args[2])
					if cmd.Args[2] == "version" {
						fmt.Fprintf(cmd.Stdout, `{"provider_selections":{%q:"0.1.0"}}`, provider.Source)
					}
					return nil
				},
			}
			if code, err := call.run(config, nil); code != 2 || err == nil || !strings.Contains(err.Error(), "selects 0.1.0") {
				t.Fatalf("stale provider lock reached plan: %d, %v", code, err)
			}
			want := []string{"version"}
			if !direct {
				want = append([]string{"init"}, want...)
			}
			if !reflect.DeepEqual(commands, want) {
				t.Fatalf("commands = %q, want %q", commands, want)
			}
		})
	}
}

func TestInitializationFreeCommandsDoNotReadProviderCache(t *testing.T) {
	for _, command := range []string{"init", "fmt", "version", "--version", "help", "--help"} {
		if commandUsesProviders([]string{command}) {
			t.Errorf("%s unnecessarily inspects provider cache", command)
		}
	}
	for _, command := range []string{"plan", "apply", "validate", "test", "show", "providers", "future-command"} {
		if !commandUsesProviders([]string{command}) {
			t.Errorf("%s bypasses provider verification", command)
		}
	}
}

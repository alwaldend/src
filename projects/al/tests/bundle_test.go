package tests

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

// TestDevelopmentBundle checks that an ordinary unstamped Bazel invocation
// produces inspectable, cache-stable development bytes and retains its manifest.
func TestDevelopmentBundle(t *testing.T) {
	directory, err := runfiles.Rlocation(os.Getenv("DEVELOPMENT_BUNDLE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(directory, "release.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		// Name identifies the development version.
		Name string `json:"name"`
		// Channel distinguishes development from publishable releases.
		Channel string `json:"channel"`
		// SourceCommit must not invalidate ordinary development outputs.
		SourceCommit string `json:"sourceCommit"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Name != "0.0.0-dev" || manifest.Channel != "development" || manifest.SourceCommit != "" {
		t.Fatalf("incorrect development manifest: %s", data)
	}
	payload, err := os.ReadFile(filepath.Join(directory, "files", "payload.txt"))
	if err != nil || string(payload) != "Standalone release fixture.\n" {
		t.Fatalf("incorrect development payload %q: %v", payload, err)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "development-release.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestStandaloneBundle inspects a bundle produced through the public Bzlmod API.
func TestStandaloneBundle(t *testing.T) {
	directory, err := runfiles.Rlocation(os.Getenv("BUNDLE_PATH"))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(directory, "release.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		// Name records the explicit version fixture.
		Name string `json:"name"`
		// SourceCommit binds the rule output to supplied build metadata.
		SourceCommit string `json:"sourceCommit"`
		// Items lists the output payloads.
		Items []struct {
			// File describes one payload.
			File struct {
				// Name identifies the payload basename.
				Name string `json:"name"`
				// Hashes records the emitted digest algorithms and values.
				Hashes []struct{ Algo, Content string } `json:"hashes"`
			} `json:"file"`
		} `json:"items"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Name != "2026.44.0" || manifest.SourceCommit != "0123456789012345678901234567890123456789" || len(manifest.Items) != 1 {
		t.Fatalf("incorrect rule manifest: %s", data)
	}
	payload, err := os.ReadFile(filepath.Join(directory, "files", manifest.Items[0].File.Name))
	if err != nil {
		t.Fatal(err)
	}
	if string(payload) != "Standalone release fixture.\n" {
		t.Fatalf("incorrect bundled bytes: %q", payload)
	}
	verified := false
	for _, hash := range manifest.Items[0].File.Hashes {
		if hash.Algo == "SHA-256" && hash.Content == fmt.Sprintf("%x", sha256.Sum256(payload)) {
			verified = true
		}
	}
	if !verified {
		t.Fatalf("rule output does not match its digest: %s", data)
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "standalone-release.json"), data, 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

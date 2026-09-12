package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFixture(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRunDiscoversNestedFilesAndPrintsTable(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "infra/empty/dnsconfig.json", `{"records":{}}`)
	writeFixture(t, root, "projects/nested/MODULE.bazel", `module(name = "nested")`)
	writeFixture(t, root, "projects/nested/dnsconfig.json", `{"records":{"site":{"CNAME":{"name":"site","target":"example.com."},"dsp":["global"]}}}`)
	var output bytes.Buffer
	if err := run(nil, root, &output); err != nil {
		t.Fatal(err)
	}
	for _, row := range []string{
		"| infra/empty/dnsconfig.json | - | - | - |",
		"| projects/nested/dnsconfig.json | site.alwaldend.com | CNAME | global |",
	} {
		if !strings.Contains(output.String(), row) {
			t.Fatalf("runtime table missing %s:\n%s", row, output.String())
		}
	}
}

func TestRunRejectsCrossFileDomainAcrossTypesAndViews(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "infra/first/dnsconfig.json", `{"records":{"first":{"A":{"name":"@","address":"192.0.2.1"},"dsp":["global"]}}}`)
	writeFixture(t, root, "projects/second/dnsconfig.json", `{"records":{"second":{"AAAA":{"name":"EXAMPLE.COM.","address":"2001:db8::1"},"dsp":["dc1"]}}}`)
	var output bytes.Buffer
	err := run([]string{"--workspace", root, "--zone", "example.com"}, "", &output)
	if err == nil || !strings.Contains(err.Error(), "domain example.com is managed by multiple dnsconfig.json files") {
		t.Fatalf("split domain ownership accepted: %v", err)
	}
	if output.Len() != 0 {
		t.Fatalf("invalid ownership printed a success table: %s", output.String())
	}
}

func TestRunRereadsFilesWithoutAnInventory(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "infra/first/dnsconfig.json", `{"records":{"first":{"A":{"name":"service","address":"192.0.2.1"},"dsp":["all"]}}}`)
	if err := run(nil, root, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "projects/new/dnsconfig.json", `{"records":{"second":{"A":{"name":"service.alwaldend.com","address":"192.0.2.2"},"dsp":["dc1"]}}}`)
	if err := run(nil, root, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "projects/new/dnsconfig.json") {
		t.Fatalf("newly added duplicate source was not discovered: %v", err)
	}
}

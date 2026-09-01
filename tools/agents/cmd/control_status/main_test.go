package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"git.alwaldend.com/alwaldend/src/tools/agents/control"
)

func TestControlStatusEmitsPackagesAndAsset(t *testing.T) {
	root := t.TempDir()
	kernel, err := control.New(control.KernelOptions{
		Root:      filepath.Join(root, "control"),
		Namespace: "task-a",
		RuntimeID: "runtime-a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterPackage("pkg.a", "project", "rev-1", "hash-a", 5_000); err != nil {
		t.Fatal(err)
	}
	if err := kernel.Mark("pkg.a", control.PackageReady); err != nil {
		t.Fatal(err)
	}
	if err := kernel.Snapshot(); err != nil {
		t.Fatal(err)
	}
	if _, err := kernel.PublishLease("worker-a", "rev-1", "", 3_600_000); err != nil {
		t.Fatal(err)
	}
	if _, err := kernel.PublishAsset("worker-a", "rev-1", "hash-a", "manifests/a.json"); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", filepath.Join(root, "workspace"),
		"--control-root", "../control",
		"--namespace", "task-a",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result status
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("invalid status JSON: %v", err)
	}
	if !result.Healthy {
		t.Fatalf("healthy control reported unhealthy")
	}
	if len(result.Packages) != 1 || result.Packages[0].State != control.PackageReady {
		t.Fatalf("unexpected packages: %+v", result.Packages)
	}
	if result.Asset == nil || result.Asset.Revision != "rev-1" {
		t.Fatalf("unexpected asset: %+v", result.Asset)
	}
}

func TestControlStatusMarkdownUsesSameData(t *testing.T) {
	root := t.TempDir()
	kernel, err := control.New(control.KernelOptions{
		Root:      filepath.Join(root, "control"),
		Namespace: "task-a",
		RuntimeID: "runtime-a",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.RegisterPackage("pkg.a", "project", "rev-1", "hash-a", 5_000); err != nil {
		t.Fatal(err)
	}
	var markdownOut bytes.Buffer
	if err := run([]string{
		"--workspace-root", filepath.Join(root, "workspace"),
		"--control-root", "../control",
		"--namespace", "task-a",
		"--markdown",
	}, &markdownOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(markdownOut.String(), "task-a") {
		t.Fatalf("markdown does not state the same data: %s", markdownOut.String())
	}
}

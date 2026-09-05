package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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
	if _, err := kernel.PublishLease("worker-a", "rev-1", "", time.Hour); err != nil {
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
	if result.Healthy == nil || !*result.Healthy {
		t.Fatalf("healthy control reported unhealthy")
	}
	if result.Freshness != "unknown" || result.HealthBasis != "persisted-package-snapshot" {
		t.Fatalf("persisted health lacks its limitations: %+v", result)
	}
	if len(result.Packages) != 1 || result.Packages[0].State != control.PackageReady {
		t.Fatalf("unexpected packages: %+v", result.Packages)
	}
	if result.Asset == nil || result.Asset.Revision != "rev-1" {
		t.Fatalf("unexpected asset: %+v", result.Asset)
	}
	if result.RuntimeID != "" || result.Asset.RuntimeID != "runtime-a" {
		t.Fatalf("publication identity must not identify the snapshot writer: %+v", result)
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
	if !strings.Contains(markdownOut.String(), "Snapshot healthy: unavailable") {
		t.Fatalf("missing snapshot must have unavailable health: %s", markdownOut.String())
	}
}

func TestControlStatusMissingRootDoesNotCreateDirectories(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	if err := run([]string{"--workspace-root", root, "--control-root", "missing"}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result status
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Healthy != nil || result.RuntimeID != "" || result.Freshness != "unavailable" {
		t.Fatalf("missing observations invented runtime state: %+v", result)
	}
	if !strings.Contains(strings.Join(result.Unavailable, "\n"), "package snapshot unavailable") {
		t.Fatalf("missing snapshot is not explained: %+v", result)
	}
	if _, err := os.Stat(filepath.Join(root, "missing")); !os.IsNotExist(err) {
		t.Fatalf("status must not create a control root: %v", err)
	}
}

func TestControlStatusPersistedFailedPackageIsUnhealthy(t *testing.T) {
	root := t.TempDir()
	kernel, err := control.New(control.KernelOptions{
		Root: filepath.Join(root, "control"), Namespace: "task-a", RuntimeID: "writer",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.Mark("pkg.a", control.PackageFailed); err != nil {
		t.Fatal(err)
	}
	if err := kernel.Snapshot(); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root, "--control-root", "control", "--namespace", "task-a",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result status
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Healthy == nil || *result.Healthy {
		t.Fatalf("failed persisted package must be unhealthy: %+v", result)
	}
}

func TestControlStatusUnavailableSnapshots(t *testing.T) {
	for _, content := range []string{
		`[]`, `null`, `{`, `[{}]`,
		`[{"id":"pkg.a","state":"invented","observationTime":"2020-01-01T00:00:00Z"}]`,
		`[{"id":"pkg.a","state":"ready","observationTime":"yesterday"}]`,
		`[{"id":"pkg.a","state":"ready","observationTime":"2020-01-01T00:00:00Z","unknown":true}]`,
	} {
		t.Run(content, func(t *testing.T) {
			root := t.TempDir()
			assetDir := filepath.Join(root, "control", "assets", "default")
			if err := os.MkdirAll(assetDir, 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(assetDir, control.SnapshotFile), []byte(content), 0o644); err != nil {
				t.Fatal(err)
			}
			var stdout bytes.Buffer
			if err := run([]string{"--workspace-root", root, "--control-root", "control"}, &stdout); err != nil {
				t.Fatal(err)
			}
			var result status
			if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.Healthy != nil || result.Freshness != "unavailable" {
				t.Fatalf("invalid snapshot invented health: %+v", result)
			}
			if _, err := os.Stat(filepath.Join(root, "control", "locks")); !os.IsNotExist(err) {
				t.Fatalf("reading existing observations must not create locks: %v", err)
			}
		})
	}
}

func TestSnapshotHealthStates(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		state control.PackageState
		want  bool
	}{
		{control.PackageReady, true},
		{control.PackageDegraded, true},
		{control.PackageDisabled, true},
		{control.PackageLoading, false},
		{control.PackageFailed, false},
		{control.PackageTimeout, false},
		{control.PackageDraining, false},
	} {
		t.Run(string(test.state), func(t *testing.T) {
			packages := []control.PackageStatus{
				{ID: "pkg.a", State: control.PackageReady, ObservationTime: now.Format(time.RFC3339Nano)},
				{ID: "pkg.b", State: test.state, ObservationTime: now.Format(time.RFC3339Nano)},
			}
			healthy, freshness, err := snapshotHealth(packages, now, time.Minute)
			if err != nil || healthy == nil || *healthy != test.want || freshness != "within-age-bound" {
				t.Fatalf("snapshot health = %v, %s, %v; want %t within age bound", healthy, freshness, err, test.want)
			}
		})
	}
}

func TestSnapshotHealthObservationAge(t *testing.T) {
	now := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	for _, test := range []struct {
		name      string
		age       time.Duration
		maxAge    time.Duration
		freshness string
		available bool
	}{
		{"no freshness policy", time.Hour, 0, "unknown", true},
		{"age equals bound", time.Minute, time.Minute, "within-age-bound", true},
		{"stale", time.Hour, time.Minute, "stale", false},
		{"future timestamp", -time.Minute, time.Minute, "unavailable", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			packages := []control.PackageStatus{{
				ID: "pkg.a", State: control.PackageReady,
				ObservationTime: now.Add(-test.age).Format(time.RFC3339Nano),
			}}
			healthy, freshness, err := snapshotHealth(packages, now, test.maxAge)
			if (healthy != nil) != test.available || (err == nil) != test.available || freshness != test.freshness {
				t.Fatalf("snapshot health = %v, %s, %v; want available=%t freshness=%s", healthy, freshness, err, test.available, test.freshness)
			}
		})
	}
}

func TestControlStatusRejectsNegativeAgeBound(t *testing.T) {
	if _, err := parseFlags([]string{"--workspace-root", ".", "--max-snapshot-age", "-1s"}); err == nil {
		t.Fatal("negative age bound must be rejected")
	}
}

func TestControlStatusRejectsStaleSnapshot(t *testing.T) {
	root := t.TempDir()
	kernel, err := control.New(control.KernelOptions{
		Root: filepath.Join(root, "control"), Namespace: "default",
		ObserveNow: func() time.Time { return time.Now().Add(-time.Hour) },
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := kernel.Mark("pkg.a", control.PackageReady); err != nil {
		t.Fatal(err)
	}
	if err := kernel.Snapshot(); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	if err := run([]string{
		"--workspace-root", root, "--control-root", "control", "--max-snapshot-age", "5m",
	}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result status
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Healthy != nil || result.Freshness != "stale" || result.MaxSnapshotAge != "5m0s" {
		t.Fatalf("stale observations must have unavailable health: %+v", result)
	}
}
